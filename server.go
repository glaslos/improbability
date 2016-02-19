package main

import (
  "log"
  "net/http"
  "encoding/json"
  "sync"

  "github.com/gorilla/mux"
  "github.com/seiflotfy/cuckoofilter"
  "github.com/geetarista/go-bloomd/bloomd"
  "github.com/spencerkimball/cbfilter"
  "github.com/AndreasBriese/bbloom"
  "github.com/zhenjl/bloom"
  "github.com/zhenjl/bloom/partitioned"
  "github.com/spaolacci/murmur3"
)

var cfilters = make(map[string]*cuckoofilter.CuckooFilter)
var bfilters = make(map[string]*bloomd.Filter)
var cbfilters = make(map[string]*cbfilter.Filter)
var bbfilters = make(map[string]bbloom.Bloom)
var pbfilters = make(map[string]bloom.Bloom)

var mutex = &sync.Mutex{}

type Response struct {
  Item      string  `json:"item"`
  Result    bool    `json:"result"`
  Method    string  `json:"method"`
  Endpoint  string  `json:"endpoint"`
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/api/", Index)
  router.HandleFunc("/api/cuckoofilter/{name:[a-z]+}", CuckooFilter)
  router.HandleFunc("/api/bloomd/{name:[a-z]+}", BloomdFilter)
  router.HandleFunc("/api/cbfilter/{name:[a-z]+}", CBFilter)
  router.HandleFunc("/api/bbfilter/{name:[a-z]+}", BBFilter)
  router.HandleFunc("/api/pbfilter/{name:[a-z]+}", PBFilter)
  log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
  resp := Response{Method: r.Method, Endpoint: r.URL.Path}
  json.NewEncoder(w).Encode(resp)
}

func CuckooFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  cf := cfilters[name]
  if r.Method == "PUT" {
    cf := cuckoofilter.NewDefaultCuckooFilter()
    cfilters[name] = cf
    resp.Result = true
  } else if r.Method == "POST" {
    resp.Result = cf.InsertUnique(item)
  } else if r.Method == "GET" {
    resp.Result = cf.Lookup(item)
  } else if r.Method == "DELETE" {
    resp.Result = cf.Delete(item)
  }
  json.NewEncoder(w).Encode(resp)
}

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

func CBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := r.URL.Query().Get("item")
  resp := Response{Item: item, Method: r.Method, Endpoint: r.URL.Path}
  cbf := cbfilters[name]
  if r.Method == "PUT" {
    const (
      N = 100000
      B = 7
      FP = 0.01
    )
    cbf, err := cbfilter.NewFilter(N, B, FP)
    if err != nil {
       log.Fatal("error creating filter:", err)
    }
    cbfilters[name] = cbf
    resp.Result = true
  } else if r.Method == "POST" {
    cbf.AddKey(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = cbf.HasKey(item)
  } else if r.Method == "DELETE" {
    resp.Result = cbf.RemoveKey(item)
  }
  json.NewEncoder(w).Encode(resp)
}

func BBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  bf := bbfilters[name]
  if r.Method == "PUT" {
    bf := bbloom.New(float64(1<<16), float64(0.01))
    bbfilters[name] = bf
    resp.Result = true
  } else if r.Method == "POST" {
    bf.AddTS(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = bf.HasTS(item)
  }
  json.NewEncoder(w).Encode(resp)
}

func PBFilter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  item := []byte(r.URL.Query().Get("item"))
  resp := Response{Item: r.URL.Query().Get("item"), Method: r.Method, Endpoint: r.URL.Path}
  bf := pbfilters[name]
  if r.Method == "PUT" {
    bf := partitioned.New(100000)
    bf.SetHasher(murmur3.New64())
    pbfilters[name] = bf
    resp.Result = true
  } else if r.Method == "POST" {
    bf.Add(item)
    resp.Result = true
  } else if r.Method == "GET" {
    resp.Result = bf.Check(item)
  }
  json.NewEncoder(w).Encode(resp)
}
