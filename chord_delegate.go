package chordstore

import (
	"bytes"
	"io"
	"log"

	chord "github.com/euforia/go-chord"
)

// ChordDelegate handles transferring and replicating keys to the appropriate
// vnode based on chord events.
type ChordDelegate struct {
	Store Store
}

func (cd *ChordDelegate) transferVnodeData(src, dst *chord.Vnode) error {
	buf := new(bytes.Buffer)
	err := cd.Store.Snapshot(src, buf)

	if err != nil {
		if err != io.EOF {
			return err
		}
	}

	// Skip if no data
	if buf.Len() < 1 {
		log.Println("DBG [transfer] Nothing to transfer.  No data!")
		return nil
	}

	log.Printf("DBG [transfer] Copying %s --> %s", shortID(src), shortID(dst))
	return cd.Store.Restore(dst, buf)
}

// NewPredecessor is called when a new predecessor is found
func (cd *ChordDelegate) NewPredecessor(local, remoteNew, remotePrev *chord.Vnode) {
	log.Printf("DBG [chord] NewPredecessor local=%s remote=%s old=%s", shortID(local), shortID(remoteNew), shortID(remotePrev))
	// Ship a copy of the local vnode to the remote
	if err := cd.transferVnodeData(local, remoteNew); err != nil {
		log.Printf("ERR [transfer] %s %s %v", local.StringID(), remoteNew.StringID(), err)
	}

}

// Leaving is called when local node is leaving the ring
func (cd *ChordDelegate) Leaving(local, pred, succ *chord.Vnode) {
	log.Printf("DBG [chord] Leaving local=%s succ=%s", shortID(local), shortID(succ))
	if err := cd.transferVnodeData(local, succ); err != nil {
		log.Printf("ERR [transfer] %s %s %v", local.StringID(), succ.StringID(), err)
	}
}

// PredecessorLeaving is called when a predecessor leaves
func (cd *ChordDelegate) PredecessorLeaving(local, remote *chord.Vnode) {
	log.Printf("DBG [chord] PredecessorLeaving local=%s remote=%s", shortID(local), shortID(remote))
}

// SuccessorLeaving is called when a successor leaves
func (cd *ChordDelegate) SuccessorLeaving(local, remote *chord.Vnode) {
	log.Printf("DBG [chord] SuccessorLeaving local=%s remote=%s", shortID(local), shortID(remote))
}

// Shutdown is called when the node is shutting down
func (cd *ChordDelegate) Shutdown() {
	log.Println("[chord] Shutdown")
}

func shortID(vn *chord.Vnode) string {
	if vn != nil {
		return vn.Host + "/" + vn.StringID()[:12]
	}
	return "<nil>"
}
