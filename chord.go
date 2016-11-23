package chordstore

import (
	"fmt"
	"log"
	"net"
	"time"

	chord "github.com/euforia/go-chord"

	"google.golang.org/grpc"
)

// ChordConfig is an extend config for chord.
type ChordConfig struct {
	*chord.Config
	// Port to listen on
	BindAddr string
	// Address of an existing cluster member to join
	Peers []string
	// Chord transport.  This needs to be explicitly shutdown.
	Transport *chord.GRPCTransport `json:"-"`
	// Timeout for chord rpc calls
	Timeout time.Duration
	// Time before a outgoing connection idle conn is reaped
	ConnMaxIdle time.Duration
}

// Config contains the ring config along with the app and transport config.
type Config struct {
	// Underlying chord config
	Chord *ChordConfig
	// GRPC server. This is so multiple services can be registered with grpc
	Server *grpc.Server `json:"-"`
	// This can be provided or a tcp listener is created using the bind address.
	// It needs to be closed on shutdown
	Listener net.Listener `json:"-"`
	// Key replication count.
	Replicas int
	// Actual chord ring
	Ring *chord.Ring
}

// ChordDelegate returns a type delegate
func (cfg *Config) ChordDelegate() *ChordDelegate {
	d, _ := cfg.Chord.Delegate.(*ChordDelegate)
	return d
}

// DefaultConfig returns a semi-sane default configuration
func DefaultConfig(bindAddr string) *Config {
	c := &Config{
		Chord: &ChordConfig{
			BindAddr:    bindAddr,
			Timeout:     time.Second * 5,
			ConnMaxIdle: time.Second * 300,
			Peers:       []string{},
		},
		Server:   grpc.NewServer(),
		Replicas: 3,
	}

	c.Chord.Config = chord.DefaultConfig(c.Chord.BindAddr)
	c.Chord.StabilizeMin = time.Duration(5 * time.Second)
	c.Chord.StabilizeMax = time.Duration(15 * time.Second)
	c.Chord.Delegate = &ChordDelegate{}
	//c.Chord.HashFunc = sha256.New
	return c
}

// InitChord initializes the ring
func initChordRing(cfg *Config) (err error) {
	cfg.Chord.Transport = chord.NewGRPCTransport(
		cfg.Listener, cfg.Server, cfg.Chord.Timeout, cfg.Chord.ConnMaxIdle)

	if len(cfg.Chord.Peers) == 0 {
		log.Println("[chord] Creating ring...")
		cfg.Ring, err = chord.Create(cfg.Chord.Config, cfg.Chord.Transport)
		return
	}

	log.Println("[chord] Joining ring...")
	for _, peer := range cfg.Chord.Peers {
		log.Printf("[chord] Trying peer: %s", peer)
		// NOTE: If the peer has not cleanly left the ring and issues a join,
		// it may fail as other nodes are trying to contact it while it is trying
		// to join.
		ring, e := chord.Join(cfg.Chord.Config, cfg.Chord.Transport, peer)
		if e == nil {
			cfg.Ring = ring
			return
		}
		log.Printf("[chord] Peer contact failed %s: %s", peer, err)
	}
	return fmt.Errorf("exhausted all peers")
}
