package chordstore

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	chord "github.com/euforia/go-chord"
)

var (
	testReplicas = 4
	testKey      = []byte("mytestkey")

	testKeyValue = map[string][]byte{
		"foo":    []byte("foo"),
		"bar":    []byte("bar"),
		"baz":    []byte("baz"),
		"xyz":    []byte("xyz"),
		"bizzle": []byte("bizzle"),
	}
)

func initConfig(p int, j ...string) (*Config, error) {
	laddr := fmt.Sprintf("127.0.0.1:%d", p)

	c1 := DefaultConfig(laddr)

	c1.Chord.StabilizeMin = time.Second
	c1.Chord.StabilizeMax = 5 * time.Second
	c1.Chord.Peers = j

	var err error
	c1.Listener, err = net.Listen("tcp", laddr)
	return c1, err
}

func Test_ChordStore(t *testing.T) {
	c1, err := initConfig(44332)
	if err != nil {
		t.Fatal(err)
	}

	cs1, e1 := NewChordStore(c1, &MemKeyValueStore{})
	if e1 != nil {
		t.Fatal(e1)
	}

	cd := c1.ChordDelegate()
	if cd == nil {
		t.Log("delegate not set")
		t.Fail()
	}

	<-time.After(200 * time.Millisecond)
	c2, err := initConfig(55443, "127.0.0.1:44332")
	if err != nil {
		t.Fatal(err)
	}

	cs2, e2 := NewChordStore(c2, &MemKeyValueStore{})
	if e2 != nil {
		t.Fatal(e2)
	}
	// Put
	rsp, err := cs2.PutKey(testReplicas, testKey, []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rsp) != testReplicas {
		t.Fatal("mismatch")
	}
	t.Logf("PUT %+v", *rsp[0])

	// Get
	rsp2, e2 := cs1.GetKey(testReplicas, testKey)
	if e2 != nil {
		t.Fatal(e2)
	}
	if len(rsp2) != testReplicas {
		t.Fatal("mismatch")
	}
	if string(rsp2[0].Data) != "value" {
		t.Fatal("wrong value stored")
	}
	t.Logf("GET %+v", rsp2)

	// Update
	prevHash := sha256.Sum256([]byte("value"))
	_, e21 := cs2.UpdateKey(testReplicas, prevHash[:], testKey, []byte("newValue"))
	if e21 != nil {
		t.Fatal(e21)
	}

	r5, e5 := cs1.GetKey(testReplicas, testKey)
	if e5 != nil {
		t.Fatal(e5)
	}

	tr := r5[0]
	if tr.Err != nil {
		t.Fatal(tr.Err)
	}
	if string(tr.Data) != "newValue" {
		t.Fatal("update value mismatch")
	}

	rf, e := cs1.UpdateKey(1, prevHash[:], testKey, []byte("newValue"))
	if e != nil {
		t.Log("request failed")
		t.Fail()
	}
	if rf[0].Err == nil {
		t.Fatal("should fail with invalid hash")
	}

	//<-time.After(100 * time.Millisecond)
	dr, err := cs2.RemoveKey(testReplicas, testKey)
	if err != nil {
		t.Fatal(err)
	}

	rsp3, err := cs2.GetKey(testReplicas, testKey)
	if err != nil {
		t.Fatal(err)
	}
	if rsp3[0].Err == nil {
		t.Fatal("should have an error")
	}

	if len(dr) != testReplicas {
		t.Fatal("mismatch on remove")
	}

	c1.Ring.Shutdown()
	c2.Ring.Shutdown()

	cs1.Shutdown()
	cs2.Shutdown()
}

func Test_VnodeData_Marshal(t *testing.T) {
	tvn := VnodeData{
		Vnode: &chord.Vnode{
			Id:   []byte("foo"),
			Host: "0.0.0.0:0",
		},
		Data: []byte("testdata"),
	}
	b, err := json.Marshal(tvn)
	if err != nil {
		t.Fatal(err)
	}

	var v VnodeData
	if err = json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}

	if string(v.Vnode.Id) != "foo" {
		t.Fatal("wrong value")
	}
	if string(v.Vnode.Host) != "0.0.0.0:0" {
		t.Fatal("wrong value")
	}
	if string(v.Data) != "testdata" {
		t.Fatal("wrong value")
	}
}
