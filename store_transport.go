package chordstore

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	chord "github.com/euforia/go-chord"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type rpcOutClient struct {
	host string
	c    DHTClient
	conn *grpc.ClientConn
}

// ChordStoreTransport is a remote store
type ChordStoreTransport struct {
	plock    sync.Mutex                 // connection pool lock
	pool     map[string][]*rpcOutClient //conneciton pool
	shutdown int32
}

// NewChordStoreTransport initialzed with and empty pool
func NewChordStoreTransport() *ChordStoreTransport {
	return &ChordStoreTransport{pool: map[string][]*rpcOutClient{}}
}

func (st *ChordStoreTransport) Snapshot(vn *chord.Vnode, wr io.Writer) error {
	out, err := st.getClient(vn.Host)
	if err != nil {
		return err
	}
	defer st.returnClient(out)

	cli, err := out.c.SnapshotRPC(context.Background(), vn)
	if err != nil {
		return err
	}

	for {

		ds, e := cli.Recv()
		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}

		wr.Write(ds.Data)
	}

	err = mergeErrors(err, cli.CloseSend())
	return err
}

func (st *ChordStoreTransport) Restore(vn *chord.Vnode, rd io.Reader) error {
	out, err := st.getClient(vn.Host)
	if err != nil {
		return err
	}
	defer st.returnClient(out)

	ctx := context.Background()
	cli, err := out.c.RestoreRPC(ctx)
	if err != nil {
		return err
	}
	if err = cli.SendMsg(vn); err != nil {
		return err
	}

	buf := make([]byte, 65519)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		ds := &DataStream{Data: buf[:n]}
		if err = cli.Send(ds); err != nil {
			return err
		}
	}

	ersp, err := cli.CloseAndRecv()
	if err == nil && len(ersp.Err) > 0 {
		err = fmt.Errorf(ersp.Err)
	}

	return err
}

// GetKey via an rpc call
func (st *ChordStoreTransport) GetKey(vn *chord.Vnode, key []byte) ([]byte, error) {
	out, err := st.getClient(vn.Host)
	if err == nil {
		defer st.returnClient(out)

		var resp *DHTBytesErr
		resp, err = out.c.GetKeyRPC(context.Background(), &DHTBytes{B: key, Vn: vn})
		if err == nil {
			if resp.Err == "" {
				return resp.B, nil
			}
			err = fmt.Errorf(resp.Err)
		}
	}
	return nil, err
}

func (st *ChordStoreTransport) GetObject(vn *chord.Vnode, key []byte) (io.Reader, error) {
	out, err := st.getClient(vn.Host)
	if err != nil {
		return nil, err
	}
	defer st.returnClient(out)

	cli, err := out.c.GetObjectRPC(context.Background(), &DHTBytes{B: key, Vn: vn})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	for {
		ds, e := cli.Recv()
		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}
		buf.Write(ds.Data)
	}

	return buf, err
}

func (st *ChordStoreTransport) PutObject(vn *chord.Vnode, key []byte, rd io.Reader) error {
	out, err := st.getClient(vn.Host)
	if err != nil {
		return err
	}
	defer st.returnClient(out)

	cli, err := out.c.PutObjectRPC(context.Background())
	if err != nil {
		return err
	}

	// Send header with vnode and objec id
	mt := &DHTBytes{Vn: vn, B: key}
	if err = cli.SendMsg(mt); err != nil {
		return err
	}

	// Send object data
	buf := make([]byte, 65519)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		/*if err == io.EOF {
			break
		} else if err != nil {
			return err
		}*/

		ds := &DataStream{Data: buf[:n]}
		if err = cli.Send(ds); err != nil {
			return err
		}
	}

	ersp, err := cli.CloseAndRecv()
	if err != nil {
		return err
	}

	if len(ersp.Err) > 0 {
		return fmt.Errorf(ersp.Err)
	}
	return nil
}

// PutKey writes a key value to the vnode
func (st *ChordStoreTransport) PutKey(vn *chord.Vnode, key, value []byte) error {
	out, err := st.getClient(vn.Host)
	if err == nil {
		defer st.returnClient(out)

		var resp *chord.ErrResponse
		if resp, err = out.c.PutKeyRPC(context.Background(), &DHTKeyValue{Vn: vn, Key: key, Value: value}); err == nil {
			if resp.Err == "" {
				return nil
			}
			err = fmt.Errorf(resp.Err)
		}
	}
	return err
}

// UpdateKey updates a key value to the vnode.  The previousHash is that of the previous
// value of the key.
func (st *ChordStoreTransport) UpdateKey(vn *chord.Vnode, prevHash, key, value []byte) error {
	out, err := st.getClient(vn.Host)
	if err == nil {
		defer st.returnClient(out)

		var resp *chord.ErrResponse
		if resp, err = out.c.UpdateKeyRPC(context.Background(),
			&DHTHashKeyValue{Vn: vn, PrevHash: prevHash, Key: key, Value: value}); err == nil {

			if resp.Err == "" {
				return nil
			}
			err = fmt.Errorf(resp.Err)
		}
	}
	return err
}

// RemoveKey from a specific vnode
func (st *ChordStoreTransport) RemoveKey(vn *chord.Vnode, key []byte) error {
	out, err := st.getClient(vn.Host)
	if err == nil {
		defer st.returnClient(out)
		var resp *chord.ErrResponse
		if resp, err = out.c.RemoveKeyRPC(context.Background(), &DHTBytes{B: key, Vn: vn}); err == nil {
			if resp.Err == "" {
				return nil
			}
			err = fmt.Errorf(resp.Err)
		}
	}
	return err
}

// RemoveObject from a vnode
func (st *ChordStoreTransport) RemoveObject(vn *chord.Vnode, key []byte) error {
	out, err := st.getClient(vn.Host)
	if err == nil {
		defer st.returnClient(out)
		var resp *chord.ErrResponse
		if resp, err = out.c.RemoveObjectRPC(context.Background(), &DHTBytes{B: key, Vn: vn}); err == nil {
			if resp.Err == "" {
				return nil
			}
			err = fmt.Errorf(resp.Err)
		}
	}
	return err
}

// Shutdown the store transport
func (st *ChordStoreTransport) Shutdown() error {
	atomic.StoreInt32(&st.shutdown, 1)
	// Close all the outbound
	st.plock.Lock()
	for _, conns := range st.pool {
		for _, out := range conns {
			out.conn.Close()
		}
	}
	st.pool = nil
	st.plock.Unlock()

	return nil
}

func (st *ChordStoreTransport) getClient(host string) (*rpcOutClient, error) {
	// Check if we have a conn cached
	var out *rpcOutClient
	st.plock.Lock()
	if atomic.LoadInt32(&st.shutdown) == 1 {
		st.plock.Unlock()
		return nil, fmt.Errorf("transport is shutdown")
	}

	list, ok := st.pool[host]
	if ok && len(list) > 0 {
		out = list[len(list)-1]
		list = list[:len(list)-1]
		st.pool[host] = list
	}
	st.plock.Unlock()
	// Make a new connection
	if out == nil {
		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err == nil {
			return &rpcOutClient{host: host, c: NewDHTClient(conn), conn: conn}, nil
		}
		return nil, err
	}
	// return an existing connection
	return out, nil
}
func (st *ChordStoreTransport) returnClient(c *rpcOutClient) {
	// Push back into the pool
	st.plock.Lock()
	defer st.plock.Unlock()
	if atomic.LoadInt32(&st.shutdown) == 1 {
		c.conn.Close()
		return
	}
	list, _ := st.pool[c.host]
	st.pool[c.host] = append(list, c)
}
