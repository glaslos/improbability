package main

import (
  "log"
  "net/http"
  "encoding/json"
  "sync"

  "github.com/zhenjl/bloom"
  "github.com/AndreasBriese/bbloom"
  "github.com/spencerkimball/cbfilter"
  "github.com/seiflotfy/cuckoofilter"
  "github.com/geetarista/go-bloomd/bloomd"
)

var cfilters = make(map[string]*cuckoofilter.CuckooFilter)
var bfilters = make(map[string]*bloomd.Filter)
var cbfilters = make(map[string]*cbfilter.Filter)
var bbfilters = make(map[string]bbloom.Bloom)
var pbfilters = make(map[string]bloom.Bloom)
var sbfilters = make(map[string]bloom.Bloom)

var mutex = &sync.Mutex{}

type Response struct {
  Item      string  `json:"item"`
  Result    bool    `json:"result"`
  Method    string  `json:"method"`
  Endpoint  string  `json:"endpoint"`
}

func main() {
  router := GetRouter()
  log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
  resp := Response{Method: r.Method, Endpoint: r.URL.Path}
  json.NewEncoder(w).Encode(resp)
}
