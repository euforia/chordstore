package chordstore

import (
	"bytes"
	"fmt"
	"log"

	chord "github.com/euforia/go-chord"
)

type HealRequest struct {
	Vnode *chord.Vnode
	Key   []byte
}

type HealingEngine struct {
	q    chan HealRequest
	stop chan bool

	cs *ChordStore
}

func NewHealingEngine(cs *ChordStore) *HealingEngine {
	return &HealingEngine{
		q:    make(chan HealRequest),
		stop: make(chan bool, 1),
		cs:   cs,
	}
}

func (he *HealingEngine) Start() {
	log.Printf("[heal] Engine started. Waiting for requests...")
	for {
		select {
		case hr := <-he.q:
			if err := he.healKey(hr); err != nil {
				log.Printf("ERR failed to heal key: %s %v", hr.Key, err)
			}

		case <-he.stop:
			goto cleanup
		}
	}

cleanup:
	close(he.q)
	close(he.stop)
}

func (he *HealingEngine) healKey(hr HealRequest) error {
	vnds, err := he.cs.GetKey(3, hr.Key)
	if err != nil {
		return err
	}

	pok := []*VnodeData{}
	heal := []*VnodeData{}
	for i := 0; i < len(vnds); i++ {
		if vnds[i].Err != nil {
			heal = append(heal, vnds[i])
		} else {
			pok = append(pok, vnds[i])
		}
	}

	if len(pok) == 0 {
		return fmt.Errorf("fatal: all keys exausted: '%s'", hr.Key)
	}

	h := pok[0].Hash()
	if len(pok) > 1 {
		for i := 1; i < len(pok); i++ {
			hc := pok[i].Hash()
			if !bytes.Equal(h, hc) {
				return fmt.Errorf("inconsistent key data: '%x' != '%x'", h, hc)
			}
		}
	}

	for _, v := range heal {
		if err := he.cs.store.PutKey(v.Vnode, hr.Key, pok[0].Data); err != nil {
			log.Printf("ERR Failed heal key %s: %v", hr.Key, err)
		} else {
			log.Printf("Healed %s/%s", v.Vnode.StringID(), hr.Key)
		}
	}

	return nil
}
