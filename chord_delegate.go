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

// NewPredecessor is called when a new predecessor is found
func (cd *ChordDelegate) NewPredecessor(local, remoteNew, remotePrev *chord.Vnode) {
	log.Printf("[chord] NewPredecessor local=%s remote=%s old=%s", shortID(local), shortID(remoteNew), shortID(remotePrev))

	buf := new(bytes.Buffer)
	err := cd.Store.Snapshot(local, buf)

	if err != nil {
		if err != io.EOF {
			log.Println("ERR", err, shortID(local))
			return
		}
	}

	// Skip if no data
	if buf.Len() < 1 {
		log.Println("Nothing to transfer.  No data!")
		return
	}

	log.Printf("[chord] NewPredecessor: Copy %s ----> %s", shortID(local), shortID(remoteNew))

	if err = cd.Store.Restore(remoteNew, buf); err != nil {
		log.Println("ERR", err)
	}

}

// Leaving is called when local node is leaving the ring
func (cd *ChordDelegate) Leaving(local, pred, succ *chord.Vnode) {
	log.Printf("[chord] Leaving local=%s succ=%s", shortID(local), shortID(succ))
}

// PredecessorLeaving is called when a predecessor leaves
func (cd *ChordDelegate) PredecessorLeaving(local, remote *chord.Vnode) {
	log.Printf("[chord] PredecessorLeaving local=%s remote=%s", shortID(local), shortID(remote))
}

// SuccessorLeaving is called when a successor leaves
func (cd *ChordDelegate) SuccessorLeaving(local, remote *chord.Vnode) {
	log.Printf("[chord] SuccessorLeaving local=%s remote=%s", shortID(local), shortID(remote))
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
