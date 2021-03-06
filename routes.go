package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Index handles requests to the api index endpoint
func Index(w http.ResponseWriter, r *http.Request) {
	resp := Response{Method: r.Method, Endpoint: r.URL.Path}
	json.NewEncoder(w).Encode(resp)
}

// GetRouter returns a router for the http server
func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/", Index)
	router.HandleFunc("/api/cuckoofilter/{name:[a-z]+}", CuckooFilter)
	router.HandleFunc("/api/bloomd/{name:[a-z]+}", BloomdFilter)
	router.HandleFunc("/api/cbfilter/{name:[a-z]+}", CBFilter)
	router.HandleFunc("/api/bbfilter/{name:[a-z]+}", BBFilter)
	router.HandleFunc("/api/pbfilter/{name:[a-z]+}", PBFilter)
	router.HandleFunc("/api/sbfilter/{name:[a-z]+}", SBFilter)
	router.HandleFunc("/api/hllplus/{name:[a-z]+}", HLLPlus)
	router.HandleFunc("/api/skiplist/{name:[a-z]+}", SkipList)
	router.HandleFunc("/api/treap/{name:[a-z]+}", Treap)
	router.HandleFunc("/api/cml/{name:[a-z]+}", CountMinLog)
	router.HandleFunc("/api/cms/{name:[a-z]+}", CountMinSketch)
	router.HandleFunc("/api/topk/{name:[a-z]+}", topK)
	return router
}
