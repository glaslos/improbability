package main

import (
  gotopk "github.com/dgryski/go-topk"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)

type TKResponse struct {
  Response
  Value     []gotopk.Element  `json:"value"`
}

var topks = make(map[string]*gotopk.Stream)

func TopK(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := TKResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
  topk := topks[name]
  if r.Method == "PUT" {
    topk := gotopk.New(100)
    topks[name] = topk
    resp.Result = true
  } else if r.Method == "POST" {
    topk.Insert(item, 1)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Value = topk.Keys()
    resp.Result = true
  }
  json.NewEncoder(w).Encode(resp)
}
