package main

import (
  "github.com/seiflotfy/cuckoofilter"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)

var cfilters = make(map[string]*cuckoofilter.CuckooFilter)

func CuckooFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  cf := cfilters[name]
  if r.Method == "PUT" {
    cf := cuckoofilter.NewDefaultCuckooFilter()
    cfilters[name] = cf
    resp.Result = true
  } else if r.Method == "POST" {
    resp.Result = cf.InsertUnique(item)
  } else if r.Method == "GET" {
    resp.Result = cf.Lookup(item)
  } else if r.Method == "DELETE" {
    resp.Result = cf.Delete(item)
  }
  json.NewEncoder(w).Encode(resp)
}
