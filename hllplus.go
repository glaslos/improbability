package main

import (
  "hash"
	"hash/fnv"
  "github.com/clarkduvall/hyperloglog"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)

type HLLResponse struct {
  Response
  Value     uint  `json:"value"`
}

var hlls = make(map[string]*hyperloglog.HyperLogLogPlus)

func hash64(s string) hash.Hash64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h
}

func HLLPlus(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := HLLResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
  hll := hlls[name]
  if r.Method == "PUT" {
    hll, _ := hyperloglog.NewPlus(16)
    hlls[name] = hll
    resp.Result = true
  } else if r.Method == "POST" {
    hll.Add(hash64(item))
    resp.Result = true
  } else if r.Method == "GET" {
    hll.Count()
    resp.Result = true
  }
  json.NewEncoder(w).Encode(resp)
}
