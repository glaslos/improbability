package main

import (
  cml "github.com/seiflotfy/count-min-log"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "strconv"
)

var cmlogs = make(map[string]*cml.Sketch)

func CountMinLog(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
  cmlog := cmlogs[name]
  if r.Method == "PUT" {
    cmlog, _ := cml.NewDefaultSketch()
    cmlogs[name] = cmlog
    resp.Result = true
  } else if r.Method == "POST" {
    resp.Result = cmlog.IncreaseCount([]byte(item))
  } else if r.Method == "GET" {
    resp.Value = strconv.FormatFloat(cmlog.Frequency([]byte(item)), 'f', 2, 64)
    resp.Result = true
  }
  json.NewEncoder(w).Encode(resp)
}
