package main

import "github.com/gorilla/mux"

func GetRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/api/", Index)
  router.HandleFunc("/api/cuckoofilter/{name:[a-z]+}", CuckooFilter)
  router.HandleFunc("/api/bloomd/{name:[a-z]+}", BloomdFilter)
  router.HandleFunc("/api/cbfilter/{name:[a-z]+}", CBFilter)
  router.HandleFunc("/api/bbfilter/{name:[a-z]+}", BBFilter)
  router.HandleFunc("/api/pbfilter/{name:[a-z]+}", PBFilter)
  router.HandleFunc("/api/sbfilter/{name:[a-z]+}", SBFilter)
  return router
}
