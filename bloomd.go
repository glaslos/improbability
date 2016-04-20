package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/geetarista/go-bloomd/bloomd"
	"github.com/gorilla/mux"
)

var bfilters = make(map[string]*bloomd.Filter)

func BloomdFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	item := r.URL.Query().Get("item")
	bf := bfilters[name]
	var err error
	resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
	if r.Method == "PUT" {
		client := bloomd.NewClient("localhost:8673")
		filter := bloomd.Filter{
			Name: name,
			Conn: client.Conn,
		}
		client.CreateFilter(&filter)
		bfilters[name] = client.GetFilter(name)
		resp.Result = true
	} else if r.Method == "POST" {
		resp.Result, err = bf.Set(item)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "GET" {
		mutex.Lock()
		ret, err := bf.Multi([]string{item})
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		if len(ret) == 0 {
			resp.Result = false
		} else {
			resp.Result = ret[0]
		}
	}
	json.NewEncoder(w).Encode(resp)
}
