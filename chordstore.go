package chordstore

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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

	GetObject(vn *chord.Vnode, key []byte) (io.Reader, error)
	PutObject(vn *chord.Vnode, key []byte, rd io.Reader) error
	RemoveObject(vn *chord.Vnode, key []byte) error
}

// VnodeStore are operations for local vnodes. It also instantiates new stores
type VnodeStore interface {
	New(*chord.Vnode) (VnodeStore, error) // Instantiate a new store

	GetKey(key []byte) ([]byte, error)
	PutKey(key, value []byte) error
	UpdateKey(prevHash, key, value []byte) error
	RemoveKey(key []byte) error

	Snapshot(io.Writer) error
	Restore(io.Reader) error

	GetObject(key []byte) (io.Reader, error)
	PutObject(key []byte, rd io.Reader) error
	RemoveObject(key []byte) error
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

/*func (cs *ChordStore) GetObject(n int, key []byte) (io.Reader, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return nil, err
	}

	for _, vn := range vns {
		//log.Printf("GET try=%d vnode=%s key=%x", i+1, shortID(vn), key)
		rd, e := cs.store.GetObject(vn, key)
		if e == nil {
			return rd, nil
		}
		//log.Printf("%d GET ERR %s %x %v", i, shortID(vn), key, e)
		err = e
	}
	return nil, err
}*/

func (cs *ChordStore) GetObject(n int, key []byte) ([]*VnodeDataIO, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return nil, err
	}

	vds := make([]*VnodeDataIO, len(vns))
	for i, vn := range vns {
		vds[i] = &VnodeDataIO{Vnode: vn}

		buf := bytes.NewBuffer(nil)
		//log.Printf("GET try=%d vnode=%s key=%x", i+1, shortID(vn), key)
		rd, err := cs.store.GetObject(vn, key)
		if err == nil {
			if _, err = io.Copy(buf, rd); err == nil {
				vds[i].r = buf
				continue
			}
		}
		vds[i].Err = err
	}
	return vds, nil
}

// PutObject from reader returning the sha256 hash as the key
func (cs *ChordStore) PutObject(n int, key []byte, rd io.Reader) ([]*VnodeDataIO, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return nil, err
	}

	// Copy initial data
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, rd); err != nil {
		return nil, err
	}

	vds := make([]*VnodeDataIO, len(vns))
	for i, vn := range vns {
		vds[i] = &VnodeDataIO{Vnode: vn}
		vds[i].Err = cs.store.PutObject(vn, key, bytes.NewBuffer(buf.Bytes()))
	}

	return vds, nil
}

/*// PutObject from reader returning the sha256 hash as the key
func (cs *ChordStore) PutObject(n int, key []byte, rd io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, rd)
	if err != nil {
		return err
	}

	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return err
	}

	for _, vn := range vns {
		e := cs.store.PutObject(vn, key, bytes.NewBuffer(buf.Bytes()))
		err = mergeErrors(err, e)
	}

	return err
}*/

// RemoveObject with n copies
func (cs *ChordStore) RemoveObject(n int, key []byte) ([]*VnodeDataIO, error) {
	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return nil, err
	}

	vds := make([]*VnodeDataIO, len(vns))
	for i, vn := range vns {
		vds[i] = &VnodeDataIO{Vnode: vn}
		vds[i].Err = cs.store.RemoveObject(vn, key)
	}

	return vds, nil
}

/*func (cs *ChordStore) RemoveObject(n int, key []byte) error {
	vns, err := cs.ring.Lookup(n, key)
	if err != nil {
		return err
	}

	for _, vn := range vns {
		err = mergeErrors(err, cs.store.RemoveObject(vn, key))
	}

	return err
}*/

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
func (cs *ChordStore) UpdateKey(n int, key, value []byte) ([]*VnodeData, error) {
	// Get current
	rsp, err := cs.GetKey(n, key)
	if err != nil {
		return nil, err
	}
	// Check all copies for same hash to use as previousHash
	hash := rsp[0].Hash()
	for i := 1; i < len(rsp); i++ {
		h := rsp[i].Hash()
		if !bytes.Equal(hash[:], h[:]) {
			// TODO:
			//cs.healQ <- HealRequest{Vnode: rsp[i].Vnode, Key: key}
			return nil, fmt.Errorf("inconsistent key %x %x", hash, h)
		}
	}

	// Update each vnode from GetKey
	out := make([]*VnodeData, len(rsp))
	for i, r := range rsp {
		o := &VnodeData{Vnode: r.Vnode}
		o.Err = cs.store.UpdateKey(r.Vnode, hash[:], key, value)
		out[i] = o
	}
	return out, nil
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
	// Receive header with vnode
	var vn chord.Vnode
	if err := stream.RecvMsg(&vn); err != nil {
		return err
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
	if err := cs.store.Restore(&vn, buf); err != nil {
		ersp.Err = err.Error()
	}

	return stream.SendAndClose(ersp)
}

// PutObjectRPC server-side call
func (cs *ChordStore) PutObjectRPC(stream DHT_PutObjectRPCServer) error {
	// Receive header with vnode, and object id
	var args DHTBytes
	if err := stream.RecvMsg(&args); err != nil {
		return stream.SendAndClose(&ErrResponse{Err: err.Error()})
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

	rsp := &ErrResponse{}
	if err := cs.store.PutObject(args.Vn, args.B, buf); err != nil {
		rsp.Err = err.Error()
	}

	return stream.SendAndClose(rsp)
}

// GetObjectRPC is the server side call
func (cs *ChordStore) GetObjectRPC(key *DHTBytes, stream DHT_GetObjectRPCServer) error {
	var (
		vn     = key.Vn
		objkey = key.B
	)

	rd, err := cs.store.GetObject(vn, objkey)
	if err != nil {
		return err
	}

	out := make([]byte, 65519)
	for {
		n, err := rd.Read(out)
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

func (cs *ChordStore) RemoveObjectRPC(ctx context.Context, key *DHTBytes) (*ErrResponse, error) {
	var (
		vn     = key.Vn
		objkey = key.B
	)

	rsp := &ErrResponse{}
	if err := cs.store.RemoveObject(vn, objkey); err != nil {
		rsp.Err = err.Error()
	}

	return rsp, nil
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

// Hash returns the sha256 hash of the data
func (vd *VnodeData) Hash() []byte {
	s := sha256.Sum256(vd.Data)
	return s[:]
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

type VnodeDataIO struct {
	Vnode *chord.Vnode
	Err   error
	r     io.Reader
}

func (vd *VnodeDataIO) Reader() io.Reader {
	return vd.r
}
