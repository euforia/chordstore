package chordstore

import (
	"bytes"
	"encoding/json"
	"io"

	chord "github.com/euforia/go-chord"
	context "golang.org/x/net/context"
)

// Store is the overall store abstracting local and remote vnodes
type Store interface {
	GetKey(vn *chord.Vnode, key []byte) ([]byte, error)
	PutKey(vn *chord.Vnode, key, value []byte) error
	UpdateKey(vn *chord.Vnode, prevHash, key, value []byte) error
	RemoveKey(vn *chord.Vnode, key []byte) error
	Shutdown() error

	Snapshot(vn *chord.Vnode, wr io.Writer) error
	Restore(vn *chord.Vnode, rd io.Reader) error
}

// VnodeStore are operations for local vnodes
type VnodeStore interface {
	New() (VnodeStore, error)
	GetKey(key []byte) ([]byte, error)
	PutKey(key, value []byte) error
	UpdateKey(prevHash, key, value []byte) error
	RemoveKey(key []byte) error

	Snapshot(io.Writer) error
	Restore(io.Reader) error
}

// ChordStore implements chord ring base storage
type ChordStore struct {
	ring  *chord.Ring
	store *TransparentStore
}

// NewChordStore instantiaties a new chord store using the given VnodeStore for
// local storage.
func NewChordStore(cfg *Config, vnstore VnodeStore) (*ChordStore, error) {
	if err := initChordRing(cfg); err != nil {
		return nil, err
	}

	vnodes, err := cfg.Ring.ListVnodes(cfg.Ring.Hostname())
	if err != nil {
		return nil, err
	}

	cs := &ChordStore{ring: cfg.Ring}
	if cs.store, err = NewTransparentStore(vnstore, vnodes...); err != nil {
		return nil, err
	}
	cfg.ChordDelegate().Store = cs.store

	// register server with grpc
	RegisterDHTServer(cfg.Server, cs)
	return cs, nil
}

// PutKey with value on the ring with a replica count of n
func (cs *ChordStore) PutKey(n int, key, value []byte) ([]*VnodeData, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err == nil {
		out := make([]*VnodeData, len(vns))
		for i, vn := range vns {
			o := &VnodeData{Vnode: vn}
			o.Err = cs.store.PutKey(vn, key, value)
			out[i] = o
		}
		return out, nil
	}
	return nil, err
}

// UpdateKey with value on the ring with a replica count of n
func (cs *ChordStore) UpdateKey(n int, prevHash, key, value []byte) ([]*VnodeData, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err == nil {
		out := make([]*VnodeData, len(vns))
		for i, vn := range vns {
			o := &VnodeData{Vnode: vn}
			o.Err = cs.store.UpdateKey(vn, prevHash, key, value)
			out[i] = o
		}
		return out, nil
	}
	return nil, err
}

// GetKey with n replicas
func (cs *ChordStore) GetKey(n int, key []byte) ([]*VnodeData, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err == nil {
		out := make([]*VnodeData, len(vns))
		for i, vn := range vns {
			o := &VnodeData{Vnode: vn}
			o.Data, o.Err = cs.store.GetKey(vn, key)
			out[i] = o
		}
		return out, nil
	}
	return nil, err
}

// RemoveKey with n replicas.
func (cs *ChordStore) RemoveKey(n int, key []byte) ([]*VnodeData, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err == nil {
		out := make([]*VnodeData, len(vns))
		for i, vn := range vns {
			o := &VnodeData{Vnode: vn}
			o.Err = cs.store.RemoveKey(vn, key)
			out[i] = o
		}
		return out, nil
	}
	return nil, err
}

// PutKeyRPC server-side
func (cs *ChordStore) PutKeyRPC(ctx context.Context, dkv *DHTKeyValue) (*ErrResponse, error) {
	resp := &ErrResponse{}
	if err := cs.store.PutKey(dkv.Vn, dkv.Key, dkv.Value); err != nil {
		resp.Err = err.Error()
	}
	return resp, nil
}

// UpdateKeyRPC server-side
func (cs *ChordStore) UpdateKeyRPC(ctx context.Context, dkv *DHTHashKeyValue) (*ErrResponse, error) {
	resp := &ErrResponse{}
	if err := cs.store.UpdateKey(dkv.Vn, dkv.PrevHash, dkv.Key, dkv.Value); err != nil {
		resp.Err = err.Error()
	}
	return resp, nil
}

// GetKeyRPC  server-side
func (cs *ChordStore) GetKeyRPC(ctx context.Context, key *DHTBytes) (*DHTBytesErr, error) {
	resp := &DHTBytesErr{}
	b, err := cs.store.GetKey(key.Vn, key.B)
	if err == nil {
		resp.B = b
	} else {
		resp.Err = err.Error()
	}
	return resp, nil
}

// RemoveKeyRPC server-side
func (cs *ChordStore) RemoveKeyRPC(ctx context.Context, key *DHTBytes) (*ErrResponse, error) {
	resp := &ErrResponse{}
	if err := cs.store.RemoveKey(key.Vn, key.B); err != nil {
		resp.Err = err.Error()
	}
	return resp, nil
}

// SnapshotRPC server-side
func (cs *ChordStore) SnapshotRPC(vn *chord.Vnode, stream DHT_SnapshotRPCServer) error {
	//log.Println("SERVER SIDE SNAPSHOT", shortID(vn))
	buf := new(bytes.Buffer)

	err := cs.store.Snapshot(vn, buf)
	if err != nil {
		//log.Println("ERR", err, shortID(vn))
		if err == io.EOF {
			return nil
		}
		return err
	}

	out := make([]byte, 65519)
	for {
		n, err := buf.Read(out)
		if err != nil {
			if err != io.EOF {
				//log.Println("ERR", err, shortID(vn))
				return err
			}
			break
		}

		ds := &DataStream{Data: out[:n]}
		if err = stream.Send(ds); err != nil {
			return err
		}
	}
	return nil
}

// RestoreRPC server-side call
func (cs *ChordStore) RestoreRPC(stream DHT_RestoreRPCServer) error {

	ctx := stream.Context()
	vn, ok := ctx.Value("vnode").(*chord.Vnode)
	if !ok {
		return stream.SendAndClose(&ErrResponse{Err: "vnode not provided"})
	}

	buf := new(bytes.Buffer)

	for {
		dc, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return stream.SendAndClose(&ErrResponse{Err: err.Error()})
		}

		buf.Write(dc.Data)
	}

	ersp := &ErrResponse{}
	if err := cs.store.Restore(vn, buf); err != nil {
		ersp.Err = err.Error()
	}

	return stream.SendAndClose(ersp)
}

// Shutdown underlying stores.  This does not shutdown any underlying services.
func (cs *ChordStore) Shutdown() error {
	return cs.store.Shutdown()
}

// VnodeData is data from a vnode on the ring
type VnodeData struct {
	Vnode *chord.Vnode
	Data  []byte
	Err   error
}

// MarshalJSON custom
func (vd *VnodeData) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"vnode": map[string]string{
			"id":   vd.Vnode.StringID(),
			"host": vd.Vnode.Host,
		},
	}
	if vd.Data != nil {
		m["data"] = vd.Data
	}
	if vd.Err != nil {
		m["error"] = vd.Err.Error()
	}
	return json.Marshal(m)
}
