package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	gotopk "github.com/dgryski/go-topk"
	"github.com/gorilla/mux"
)

type tkResponse struct {
	Response
	Value []gotopk.Element `json:"value"`
}

var topks = make(map[string]*gotopk.Stream)

func topK(w http.ResponseWriter, r *http.Request) {
	topkMutex.Lock()
	defer topkMutex.Unlock()
	vars := mux.Vars(r)
	name := vars["name"]
	item := r.URL.Query().Get("item")
	topk := topks[name]
	resp := tkResponse{Response: Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}}
	if r.Method == "PUT" {
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		if size == 0 {
			size = 100
		}
		topk = gotopk.New(size)
		topks[name] = topk
		resp.Result = true
	} else if r.Method == "POST" {
		if len(item) != 0 {
			topk.Insert(item, 1)
		} else {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()
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
	return
}
