package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	maxWorkers = 100000
	url        = "http://localhost:8082/"
)

var wg sync.WaitGroup

func Worker(id int) {
	defer wg.Done()

	client := &http.Client{
		//Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		//fmt.Printf("Worker %d: Error: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	//fmt.Printf("Worker %d: Status Code: %d\n", id, resp.Status)
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, Go!")
	})

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		//fmt.Printf("Server failed: %v\n", err)
	}
}

func main() {
	go startServer()

	time.Sleep(2 * time.Second)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go Worker(i)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}
