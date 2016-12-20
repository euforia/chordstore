package chordstore

import (
	"bytes"
	"testing"

	chord "github.com/euforia/go-chord"
)

var (
	testVn1 = &chord.Vnode{Id: []byte("foobarbaz1"), Host: "127.0.0.1:9876"}
	testVn2 = &chord.Vnode{Id: []byte("foobarbaz2"), Host: "127.0.0.1:9877"}
)

func Test_MemKeyValueStore(t *testing.T) {
	st := &MemKeyValueStore{}
	kvs, _ := st.New(testVn1)
	kvs1 := kvs.(*MemKeyValueStore)
	kvs1.m = testKeyValue

	buf := new(bytes.Buffer)

	err := kvs1.Snapshot(buf)
	if err != nil {
		t.Fatal(err)
	}

	kvs, _ = st.New(testVn2)
	kvs2 := kvs.(*MemKeyValueStore)

	if err := kvs2.Restore(buf); err != nil {
		t.Fatal(err)
	}

	if len(kvs2.m) != 5 {
		t.Fatal("key count mismatch", len(kvs2.m), len(kvs1.m))
	}

	val, _ := kvs2.GetKey([]byte("bizzle"))
	if string(val) != "bizzle" {
		t.Fatal("value mismatch", string(val))
	}
}
