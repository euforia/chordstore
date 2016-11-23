package chordstore

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"sync"

	chord "github.com/euforia/go-chord"
)

// TransparentStore abstracts remote and local datastores
type TransparentStore struct {
	remote Store
	local  map[string]VnodeStore
}

// NewTransparentStore Initialized with the given vnodes
func NewTransparentStore(vnstore VnodeStore, vnodes ...*chord.Vnode) (*TransparentStore, error) {
	ts := &TransparentStore{
		local:  map[string]VnodeStore{},
		remote: NewChordStoreTransport(),
	}
	err := ts.init(vnstore, vnodes...)
	return ts, err
}

func (ts *TransparentStore) init(vnstore VnodeStore, vnodes ...*chord.Vnode) (err error) {

	for _, vn := range vnodes {
		if ts.local[vn.StringID()], err = vnstore.New(); err != nil {
			return
		}
	}

	return
}

// GetKey from local or remote vnode
func (ts *TransparentStore) GetKey(vn *chord.Vnode, key []byte) ([]byte, error) {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.GetKey(key)
	}
	return ts.remote.GetKey(vn, key)
}

// PutKey to local or remote vnode
func (ts *TransparentStore) PutKey(vn *chord.Vnode, key, value []byte) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.PutKey(key, value)
	}
	return ts.remote.PutKey(vn, key, value)
}

// UpdateKey to local or remote vnode
func (ts *TransparentStore) UpdateKey(vn *chord.Vnode, prevHash, key, value []byte) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.UpdateKey(prevHash, key, value)
	}
	return ts.remote.UpdateKey(vn, prevHash, key, value)
}

// RemoveKey from local or remote vnode
func (ts *TransparentStore) RemoveKey(vn *chord.Vnode, key []byte) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.RemoveKey(key)
	}
	return ts.remote.RemoveKey(vn, key)
}

// Snapshot a local or remote vnode
func (ts *TransparentStore) Snapshot(vn *chord.Vnode, wr io.Writer) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.Snapshot(wr)
	}
	return ts.remote.Snapshot(vn, wr)
}

// Restore a local or remote vnode
func (ts *TransparentStore) Restore(vn *chord.Vnode, rd io.Reader) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.Restore(rd)
	}
	return ts.remote.Restore(vn, rd)
}

// Shutdown remote and local stores
func (ts *TransparentStore) Shutdown() error {
	return ts.remote.Shutdown()
}

// MemKeyValueStore is an in memory key-value store
type MemKeyValueStore struct {
	mu sync.Mutex
	m  map[string][]byte
}

func (s *MemKeyValueStore) New() (VnodeStore, error) {
	return &MemKeyValueStore{m: map[string][]byte{}}, nil
}

// NewMemKeyValueStore instantiates new in memory key value store
/*func NewMemKeyValueStore() *MemKeyValueStore {
	return &MemKeyValueStore{m: map[string][]byte{}}
}*/

// GetKey from datastore
func (s *MemKeyValueStore) GetKey(key []byte) ([]byte, error) {
	if val, ok := s.m[string(key)]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("not found: %s", key)
}

// PutKey key-value
func (s *MemKeyValueStore) PutKey(key []byte, v []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := string(key)
	s.m[k] = v

	return nil
}

// UpdateKey that exists.  previousHash is the hash of the previous value of the
// key.
func (s *MemKeyValueStore) UpdateKey(prevHash, key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := string(key)
	cv, ok := s.m[k]
	if !ok {
		return fmt.Errorf("key not found: %s", k)
	}

	pv := sha256.Sum256(cv)
	if bytes.Equal(pv[:], prevHash) {
		s.m[k] = value
		return nil
	}

	return fmt.Errorf("invalid previous hash")
}

// RemoveKey a key from the datastore
func (s *MemKeyValueStore) RemoveKey(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, string(key))
	return nil
}

// Snapshot the dataset serializing and compressing it to the writer.
func (s *MemKeyValueStore) Snapshot(wr io.Writer) error {
	if len(s.m) == 0 {
		return io.EOF
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	zw := zlib.NewWriter(wr)
	defer zw.Close()

	for k, v := range s.m {
		l := append(append([]byte(k), '\x00'), v...)
		l = append(l, '\n')
		zw.Write(l)
		zw.Flush()
	}

	return nil
}

// Restore dataset from reader de-compressing and de-serializing the data to the
// datastructure
func (s *MemKeyValueStore) Restore(r io.Reader) error {

	rd, err := zlib.NewReader(r)
	if err != nil {
		return err
	}
	defer rd.Close()

	buf := bufio.NewReader(rd)

	s.mu.Lock()
	defer s.mu.Unlock()
	cnt := 0
	for {
		k, err := buf.ReadBytes('\x00')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		k = k[:len(k)-1]

		v, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s.m[string(k)] = v[:len(v)-1]
		cnt++
	}

	log.Println("Accept keys:", cnt)

	return nil
}
