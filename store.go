package chordstore

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"sync"

	chord "github.com/euforia/go-chord"
	"gopkg.in/vmihailenco/msgpack.v2"
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

// GetObject from the given vnode
func (ts *TransparentStore) GetObject(vn *chord.Vnode, key []byte) (io.Reader, error) {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.GetObject(key)
	}
	return ts.remote.GetObject(vn, key)
}

// PutObject data from the reader to the vnode
func (ts *TransparentStore) PutObject(vn *chord.Vnode, key []byte, rd io.Reader) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.PutObject(key, rd)
	}
	return ts.remote.PutObject(vn, key, rd)
}

// RemoveObject from vnode with the given key
func (ts *TransparentStore) RemoveObject(vn *chord.Vnode, key []byte) error {
	if st, ok := ts.local[vn.StringID()]; ok {
		return st.RemoveObject(key)
	}
	return ts.remote.RemoveObject(vn, key)
}

// Shutdown remote and local stores
func (ts *TransparentStore) Shutdown() error {
	return ts.remote.Shutdown()
}

// MemKeyValueStore is an in memory key-value store
type MemKeyValueStore struct {
	mu sync.Mutex
	// key-value
	m map[string][]byte
	// objects
	o map[string][]byte
}

func (s *MemKeyValueStore) New() (VnodeStore, error) {
	return &MemKeyValueStore{m: map[string][]byte{}, o: map[string][]byte{}}, nil
}

func (s *MemKeyValueStore) GetObject(key []byte) (io.Reader, error) {
	v, ok := s.o[fmt.Sprintf("%x", key)]
	if !ok {
		return nil, fmt.Errorf("key not found: %x", key)
	}

	buf := new(bytes.Buffer)
	_, err := buf.Write(v)
	return buf, err
}

func (s *MemKeyValueStore) PutObject(key []byte, rd io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, rd)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.o[fmt.Sprintf("%x", key)] = buf.Bytes()
	return nil
}

func (s *MemKeyValueStore) RemoveObject(key []byte) error {
	k := fmt.Sprintf("%x", key)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.o[k]; ok {
		delete(s.o, k)
		return nil
	}
	return fmt.Errorf("key not found: %x", key)
}

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

	return fmt.Errorf("invalid previous hash: %x != %x", pv, prevHash)
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
	if len(s.m) == 0 && len(s.o) == 0 {
		return io.EOF
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	zw := zlib.NewWriter(wr)
	defer zw.Close()

	log.Printf("[Snapshot] keys=%d objects=%d", len(s.m), len(s.o))

	for k, v := range s.o {
		log.Printf("%s %d", k, len(v))
	}

	menc := msgpack.NewEncoder(zw)
	return menc.Encode(s.m, s.o)
}

// Restore dataset from reader de-compressing and de-serializing the data to the
// datastructure
func (s *MemKeyValueStore) Restore(r io.Reader) error {

	rd, err := zlib.NewReader(r)
	if err != nil {
		return err
	}
	defer rd.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	var tk, to map[string][]byte
	dec := msgpack.NewDecoder(rd)

	if err = dec.Decode(&tk, &to); err != nil {
		return err
	}

	log.Printf("[Restore] Received keys=%d objects=%d", len(tk), len(to))

	for k, v := range tk {
		s.m[k] = v
	}
	for k, v := range to {
		s.o[k] = v
	}

	return nil
}
