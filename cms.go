package main

import (
  "github.com/shenwei356/countminsketch"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
)

type CMSResponse struct {
  Response
  Value     uint  `json:"value"`
}

var cmsketchs = make(map[string]*countminsketch.CountMinSketch)

func CountMinSketch(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := CMSResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
  cmsketch := cmsketchs[name]
  if r.Method == "PUT" {
    var varepsilon, delta float64
    varepsilon, delta = 0.1, 0.9
    cmsketch, _ := countminsketch.NewWithEstimates(varepsilon, delta)
    cmsketchs[name] = cmsketch
    resp.Result = true
  } else if r.Method == "POST" {
    cmsketch.UpdateString(item, 1)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Value = cmsketch.EstimateString(item)
    resp.Result = true
  }
  json.NewEncoder(w).Encode(resp)
}
