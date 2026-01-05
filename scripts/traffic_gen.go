package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var targets = []string{
	"http://example.com",
	"http://google.com",
	"http://github.com",
	"http://api.openai.com/v1/chat/completions",
	"http://stackoverflow.com",
	"http://reddit.com",
	"http://go.dev",
}

func main() {
	proxyURL, _ := url.Parse("http://localhost:7100")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	workers := 5

	fmt.Printf("Starting Traffic Generator with %d workers...\n", workers)
	fmt.Println("Press Ctrl+C to stop.")

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				target := targets[rand.Intn(len(targets))]

				// Introduce some 404s randomly
				if rand.Intn(10) < 2 { // 20% chance
					target += "/non-existent-path"
				}

				req, _ := http.NewRequest("GET", target, nil)
				resp, err := client.Do(req)

				status := "ERR"
				if err == nil {
					status = fmt.Sprintf("%d", resp.StatusCode)
					resp.Body.Close()
				}

				fmt.Printf("[Worker %d] %s -> %s\n", id, target, status)

				// Random sleep to simulate human/app behavior
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}
