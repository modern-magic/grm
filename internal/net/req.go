package net

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/modern-magic/grm/internal/logger"
)

func MakeRequest(urls []string, fn func(re string)) {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}
	var wg sync.WaitGroup
	results := make(chan string, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(url)
			if err != nil {
				results <- fmt.Sprintf("%s-> (%v)", url, err)
				return
			}
			defer resp.Body.Close()
			duration := time.Since(start)
			code := resp.StatusCode
			color := logger.TerminalColors.Dim
			if code >= 200 && code < 300 {
				color = logger.TerminalColors.Green
			}
			if code >= 300 && code < 400 {
				color = logger.TerminalColors.Yellow
			}
			if code >= 400 && code < 600 {
				color = logger.TerminalColors.Red
			}
			results <- fmt.Sprintf("%s-> (%s)", url, fmt.Sprintf("%s%dms%s", color, duration.Milliseconds(), logger.TerminalColors.Reset))
		}(url)
	}
	maxConcurrent := 2
	concurrent := make(chan struct{}, maxConcurrent)

	go func() {
		for i := 0; i < len(urls); i++ {
			concurrent <- struct{}{}
		}
	}()

	wg.Wait()
	close(results)
	for result := range results {
		fn(result)
		<-concurrent
	}
}
