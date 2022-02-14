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

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

	fetch()
}

func handler(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write(cache)

}

func fetch() {

	client := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	mutex.Lock()
	defer mutex.Unlock()

	var readErr error
	cache, readErr = ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

}
