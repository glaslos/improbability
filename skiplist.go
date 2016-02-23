package main

import (
  "github.com/ryszard/goskiplist/skiplist"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "strconv"
)

var slists = make(map[string]*skiplist.SkipList)

func SkipList(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
  slist := slists[name]
  if r.Method == "PUT" {
    slist := skiplist.NewIntMap()
    slists[name] = slist
    resp.Result = true
  } else if r.Method == "POST" {
    key := slist.Len()
    key ++
    slist.Set(key, item)
    resp.Result = true
  } else if r.Method == "GET" {
    i, _ := strconv.Atoi(item)
    str, ok := slist.Get(i)
    if ok {
      if test_string, ok := str.(string); ok {
        resp.Value = test_string
        resp.Result = true
      } else {
        resp.Result = false
      }
      resp.Result = true
    } else {
      resp.Result = false
    }
  }
  json.NewEncoder(w).Encode(resp)
}
