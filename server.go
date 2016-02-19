package main

import (
  "log"
  "net/http"
  "encoding/json"
)

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
