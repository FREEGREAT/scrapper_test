package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	maxWorker = 10
	url       = "https://www.google.com"
)

var wg sync.WaitGroup

func Worker(id int) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Worker %d: Error: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Worker %d: Status Code: %d\n", id, resp.StatusCode)
}

func main() {
	for i := 0; i <= maxWorker; i++ {
		wg.Add(1)
		time.Sleep(2 * time.Second)
		go Worker(i)

	}

	wg.Wait()
	fmt.Println("All workers completed")
}
