package chordstore

import (
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

	c1.Chord.StabilizeMin = 15 * time.Millisecond
	c1.Chord.StabilizeMax = 45 * time.Millisecond
	c1.Chord.Peers = j

	var err error
	c1.Listener, err = net.Listen("tcp", laddr)
	return c1, err
}

func Test_ChordDelegate(t *testing.T) {
	c1, _ := initConfig(0)
	defer c1.Listener.Close()

	cd := c1.ChordDelegate()
	if cd == nil {
		t.Fatal("delegate not set")
	}
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
	for _, r := range rsp {
		if r.Err != nil {
			t.Fatal(r.Err)
		}

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
	//prevHash := sha256.Sum256([]byte("value"))
	_, e21 := cs2.UpdateKey(testReplicas, testKey, []byte("newValue"))
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

	rf, e := cs1.UpdateKey(testReplicas, testKey, []byte("newValue2"))
	if e != nil {
		t.Log("request failed")
		t.Fail()
	}
	for _, v := range rf {
		if v.Err != nil {
			t.Fatal(err)
		}
	}
	gf, ef := cs2.GetKey(testReplicas, testKey)
	if ef != nil {
		t.Log("request failed")
		t.Fail()
	}
	for _, v := range gf {
		if v.Err != nil {
			t.Fatal(err)
		} else if string(v.Data) != "newValue2" {
			t.Fatal("value mismatch")
		}

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

func Test_ChordStore_Snapshot_Restore(t *testing.T) {
	c1, err := initConfig(22345)
	if err != nil {
		t.Fatal(err)
	}

	cs1, e1 := NewChordStore(c1, &MemKeyValueStore{})
	if e1 != nil {
		t.Fatal(e1)
	}

	<-time.After(200 * time.Millisecond)
	prsp, err := cs1.PutKey(3, []byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range prsp {
		if r.Err != nil {
			t.Fatal(r.Err)
		}
	}

	gr, _ := cs1.GetKey(3, []byte("key"))
	for _, r := range gr {
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		if string(r.Data) != "value" {
			t.Fatal("value mismatch")
		}
	}

	// Add node
	<-time.After(300 * time.Millisecond)
	c2, err := initConfig(65432, "127.0.0.1:22345")
	if err != nil {
		t.Fatal(err)
	}

	cs2, e2 := NewChordStore(c2, &MemKeyValueStore{})
	if e2 != nil {
		t.Fatal(e2)
	}

	// Wait for stabilize
	<-time.After(1000 * time.Millisecond)
	rsp, err := cs2.GetKey(3, []byte("key"))
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range rsp {
		if r.Err != nil {
			t.Fatal(r.Err)
		}
	}

	if len(rsp) != 3 {
		t.Fatalf("length mismatch have=%d want=3", len(rsp))
	}

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
