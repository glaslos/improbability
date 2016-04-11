package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Response default struct
type Response struct {
	Item     string `json:"item"`
	Result   bool   `json:"result"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

// Middleware for auth check
func Middleware(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.Method, r.URL)
		authKey := r.URL.Query().Get("auth")
		if authKey != "changeme" {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}
		//log.Println(auth_key)
		router.ServeHTTP(w, r)
	})
}

func main() {
	router := GetRouter()
	log.Fatal(http.ListenAndServe(":8008", Middleware(router)))
}
