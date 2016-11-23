package chordstore

import (
	"bytes"
	"testing"
)

func Test_MemKeyValueStore(t *testing.T) {
	st := &MemKeyValueStore{}
	kvs, _ := st.New()
	kvs1 := kvs.(*MemKeyValueStore)
	kvs1.m = testKeyValue

	buf := new(bytes.Buffer)

	err := kvs1.Snapshot(buf)
	if err != nil {
		t.Fatal(err)
	}

	kvs, _ = st.New()
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
