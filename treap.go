package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "math/rand"
  "github.com/steveyen/gtreap"
  "bytes"
)

var treaps = make(map[string]*gtreap.Treap)

func stringCompare(a, b interface{}) int {
    return bytes.Compare([]byte(a.(string)), []byte(b.(string)))
}

func Treap(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
  treap := treaps[name]
  if r.Method == "PUT" {
    treap := gtreap.NewTreap(stringCompare)
    treaps[name] = treap
    resp.Result = true
  } else if r.Method == "POST" {
    treap = treap.Upsert(item, rand.Int())
    treaps[name] = treap
    resp.Result = true
  } else if r.Method == "GET" {
    if str, ok := treap.Get(item).(string); ok {
      resp.Value = str
      resp.Result = true
      treaps[name] = treap
    } else {
      resp.Result = false
    }
  } else if r.Method == "DELETE" {
    treap = treap.Delete(item)
    resp.Result = true
    treaps[name] = treap
  }
  json.NewEncoder(w).Encode(resp)
}
