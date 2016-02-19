package main

import (
  "encoding/json"
  "github.com/AndreasBriese/bbloom"
  "github.com/gorilla/mux"
  "net/http"
)

var bbfilters = make(map[string]bbloom.Bloom)

func BBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  bf := bbfilters[name]
  if r.Method == "PUT" {
    bf := bbloom.New(float64(1<<16), float64(0.01))
    bbfilters[name] = bf
    resp.Result = true
  } else if r.Method == "POST" {
    bf.AddTS(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = bf.HasTS(item)
  }
  json.NewEncoder(w).Encode(resp)
}
