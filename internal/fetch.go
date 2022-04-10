package internal

import (
	"net/http"
	"time"
)

type FetchContext struct {
	status    string
	time      float64
	isTimeout bool
}

func fetch(uri string) FetchContext {
	start := time.Now()
	resp, err := http.Get(uri)
	if err != nil {
		return FetchContext{
			isTimeout: true,
		}
	}
	defer resp.Body.Close()

	return FetchContext{
		status:    resp.Status,
		time:      float64(time.Since(start).Seconds()),
		isTimeout: false,
	}

}
