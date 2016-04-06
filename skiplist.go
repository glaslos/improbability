package main

import (
  "github.com/ryszard/goskiplist/skiplist"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "strconv"
)

type SLResponse struct {
  Response
  Value     string  `json:"value"`
}

var slists = make(map[string]*skiplist.SkipList)

func SkipList(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := SLResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
  if r.Method == "PUT" {
    slist := skiplist.NewIntMap()
    slists[name] = slist
    resp.Result = true
  } else {
    slist, ok := slists[name]
    if ok != true {
      http.Error(w, "store not found", http.StatusNotFound)
      return
    } else {
      if r.Method == "POST" {
        key := slist.Len()
        key ++
        slist.Set(key, item)
        resp.Result = true
        resp.Value = strconv.Itoa(key)
      } else if r.Method == "GET" {
        i, _ := strconv.Atoi(item)
        val, ok := slist.Get(i)
        resp.Result = ok
        resp.Value, _ = val.(string)
      }
    }
  }
  json.NewEncoder(w).Encode(resp)
}
