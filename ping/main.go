package main

import (
	"net/http"
	"sync"
)

const (
	maxWorker = 10
	url       = "https://www.google.com"
)

var wg sync.WaitGroup

func Worker() {
	http.Get(url)
	wg.Done()

}

func main() {
	for i := 0; i <= maxWorker; i++ {
		wg.Add(1)
		go Worker()
	}

	wg.Wait()

}
