package main

import (
	"flag"
	"log"
	"net"

	"github.com/euforia/chordstore"
)

var (
	bindAddr  = flag.String("b", "127.0.0.1:3243", "Bind address")
	httpAddr  = flag.String("http", "127.0.0.1:9090", "HTTP Bind address")
	joinAddrs = flag.String("j", "", "Initial cluster membders to join")
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	cfg := chordstore.DefaultConfig(*bindAddr)
	cfg.Chord.Peers = chordstore.ParsePeersList(*joinAddrs)

	var (
		chordStore *chordstore.ChordStore
		err        error
	)

	// Init listener
	if cfg.Listener, err = net.Listen("tcp", *bindAddr); err != nil {
		log.Fatal(err)
	}

	if chordStore, err = chordstore.NewChordStore(cfg, &chordstore.MemKeyValueStore{}); err != nil {
		log.Fatal(err)
	}

	/*he := chordstore.NewHealingEngine(chordStore)
	go he.Start()*/

	admServer := chordstore.NewAdminServer(cfg, chordStore)
	admServer.Start(*httpAddr)
	select {}

}
