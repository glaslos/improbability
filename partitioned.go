package main

import (
  "github.com/zhenjl/bloom/partitioned"
  "github.com/spaolacci/murmur3"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)

func PBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  bf := pbfilters[name]
  if r.Method == "PUT" {
    bf := partitioned.New(100000)
    bf.SetHasher(murmur3.New64())
    pbfilters[name] = bf
    resp.Result = true
  } else if r.Method == "POST" {
    bf.Add(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = bf.Check(item)
  }
  json.NewEncoder(w).Encode(resp)
}
