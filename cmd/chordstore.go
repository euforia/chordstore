package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/euforia/chordstore"
)

var (
	bindAddr  = flag.String("b", "127.0.0.1:3243", "Bind address")
	httpAddr  = flag.String("http", "127.0.0.1:9090", "HTTP Bind address")
	joinAddrs = flag.String("j", "", "Initial cluster membders to join")
)

func parseJoinAddrs() []string {
	out := []string{}
	for _, v := range strings.Split(*joinAddrs, ",") {
		if c := strings.TrimSpace(v); c != "" {
			out = append(out, c)
		}
	}
	return out
}

func init() {
	flag.Parse()
}

func main() {
	cfg := chordstore.DefaultConfig(*bindAddr)
	cfg.Chord.Peers = parseJoinAddrs()

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

	he := chordstore.NewHealingEngine(chordStore)
	go he.Start()

	httpServer := &Server{store: chordStore, cfg: cfg}
	if err = http.ListenAndServe(*httpAddr, httpServer); err != nil {
		log.Fatal(err)
	}

}
