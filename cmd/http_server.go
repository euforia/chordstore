package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/euforia/chordstore"
)

type Server struct {
	store *chordstore.ChordStore
	cfg   *chordstore.Config
}

func (svr *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(svr.cfg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (svr *Server) handleObject(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		oid = ctx.Value("oid").([]byte)
		n   = ctx.Value("n").(int)
	)

	switch r.Method {
	case "GET":
		s := strings.TrimPrefix(r.URL.Path, "/object/")
		if len(s) == 0 {
			w.WriteHeader(404)
			return
		}

		oid, err := hex.DecodeString(s)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		rd, err := svr.store.GetObject(n, oid)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		_, err = io.Copy(w, rd)

	case "POST":
		err := svr.store.PutObject(n, oid, r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		//log.Printf("%x", oid)

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("%x", oid)))

	default:
		w.WriteHeader(405)
		return
	}

}

func (svr *Server) handleLookup(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		key = ctx.Value("key").([]byte)
		n   = ctx.Value("n").(int)
	)

	rsp, err := svr.cfg.Ring.Lookup(n, key)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	b, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (svr *Server) handleKV(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		key = ctx.Value("key").([]byte)
		n   = ctx.Value("n").(int)

		rsp interface{}
		err error
	)

	switch r.Method {
	case "GET":
		rsp, err = svr.store.GetKey(n, key)

	case "POST":
		var value []byte
		if value, err = ioutil.ReadAll(r.Body); err == nil {
			rsp, err = svr.store.PutKey(n, key, value)
		}

	case "PUT":
		var value []byte
		if value, err = ioutil.ReadAll(r.Body); err == nil {
			rsp, err = svr.store.UpdateKey(n, key, value)
		}

	default:
		w.WriteHeader(405)
		return
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	b, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

// ServeHTTP routes the user request and sets the context
func (svr *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := parseN(r)
	if err != nil {
		n = 1
	}

	ctx := context.WithValue(context.Background(), "n", n)

	switch {
	case strings.HasPrefix(r.URL.Path, "/config"):
		svr.handleConfig(w, r.WithContext(ctx))

	case strings.HasPrefix(r.URL.Path, "/object/"):
		s := strings.TrimPrefix(r.URL.Path, "/object/")
		if len(s) == 0 {
			w.WriteHeader(404)
			return
		}

		oid, err := hex.DecodeString(s)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		ctx = context.WithValue(ctx, "oid", oid)
		svr.handleObject(w, r.WithContext(ctx))

	case strings.HasPrefix(r.URL.Path, "/kv"):
		key := strings.TrimPrefix(r.URL.Path, "/kv/")
		if len(key) == 0 {
			w.WriteHeader(404)
			return
		}
		ctx = context.WithValue(ctx, "key", []byte(key))
		svr.handleKV(w, r.WithContext(ctx))

	case strings.HasPrefix(r.URL.Path, "/lookup"):
		key := strings.TrimPrefix(r.URL.Path, "/lookup/")
		if len(key) == 0 {
			w.WriteHeader(404)
			return
		}
		ctx := context.WithValue(ctx, "key", []byte(key))
		svr.handleLookup(w, r.WithContext(ctx))

	default:
		w.WriteHeader(404)
		return
	}

}

func parseN(r *http.Request) (int, error) {
	n := 1
	nstr, ok := r.URL.Query()["n"]
	if ok && len(nstr) > 0 {
		i, err := strconv.ParseInt(nstr[0], 10, 32)
		if err != nil {
			return 0, err
		}
		n = int(i)
	}
	return n, nil
}
