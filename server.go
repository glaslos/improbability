package main

import (
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

type Response struct {
  Item      string  `json:"item"`
  Result    bool    `json:"result"`
  Method    string  `json:"method"`
  Endpoint  string  `json:"endpoint"`
}

func Middleware(router *mux.Router) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("middleware", r.URL)
        router.ServeHTTP(w, r)
    })
}

func main() {
  router := GetRouter()
  log.Fatal(http.ListenAndServe(":8080", Middleware(router)))
}
