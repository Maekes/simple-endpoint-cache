package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var cache []byte
var mutex sync.RWMutex
var url, isURL = os.LookupEnv("URL")

func main() {

	if !isURL {
		log.Fatal("Set URL env Variable")
	}

	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			fetch()
		}
	}()
	fetch()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Write(cache)
}

func fetch() {

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	if res.Body != nil {
		defer res.Body.Close()

		mutex.Lock()
		defer mutex.Unlock()

		var readErr error
		cache, readErr = ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(readErr)
		}
	}
}
