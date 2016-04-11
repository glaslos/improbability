package main

import (
	"encoding/json"
	"net/http"

	gotopk "github.com/dgryski/go-topk"
	"github.com/gorilla/mux"
)

type tkResponse struct {
	Response
	Value []gotopk.Element `json:"value"`
}

var topks = make(map[string]*gotopk.Stream)

func topK(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	item := r.URL.Query().Get("item")
	resp := tkResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
	topk := topks[name]
	if r.Method == "PUT" {
		topk := gotopk.New(100)
		topks[name] = topk
		resp.Result = true
	} else if r.Method == "POST" {
		if len(item) != 0 {
			topk.Insert(item, 1)
		} else {
			decoder := json.NewDecoder(r.Body)
			var hashes map[string]interface{}
			err := decoder.Decode(&hashes)
			if err != nil {
				panic(err)
			}
			for hash := range hashes {
				topk.Insert(hash, 1)
			}
		}
		resp.Result = true
	} else if r.Method == "GET" {
		resp.Value = topk.Keys()
		resp.Result = true
	}
	json.NewEncoder(w).Encode(resp)
}
