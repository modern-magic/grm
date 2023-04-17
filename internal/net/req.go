package net

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/modern-magic/grm/internal/logger"
)

type RequestMessage struct {
	Err   error
	Path  string
	Alias string
	Sec   string
}

func MakeRequest(urls map[string]string, fn func(message RequestMessage)) {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}
	var wg sync.WaitGroup
	results := make(chan RequestMessage, len(urls))
	for alias, url := range urls {
		wg.Add(1)
		go func(alias, url string) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(url)
			if err != nil {
				results <- RequestMessage{
					Err:   err,
					Path:  "",
					Alias: alias,
					Sec:   "",
				}
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
			results <- RequestMessage{
				Err:   nil,
				Path:  url,
				Alias: alias,
				Sec:   fmt.Sprintf("%s%dms%s", color, duration.Milliseconds(), logger.TerminalColors.Reset),
			}
		}(alias, url)
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
