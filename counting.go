package main

import (
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/spencerkimball/cbfilter"
  "net/http"
  )

var cbfilters = make(map[string]*cbfilter.Filter)

func CBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
  cbf := cbfilters[name]
  if r.Method == "PUT" {
    const (
      N = 100000
      B = 7
      FP = 0.01
    )
    cbf, err := cbfilter.NewFilter(N, B, FP)
    if err != nil {
       log.Fatal("error creating filter:", err)
    }
    cbfilters[name] = cbf
    resp.Result = true
  } else if r.Method == "POST" {
    cbf.AddKey(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = cbf.HasKey(item)
  } else if r.Method == "DELETE" {
    resp.Result = cbf.RemoveKey(item)
  }
  json.NewEncoder(w).Encode(resp)
}
