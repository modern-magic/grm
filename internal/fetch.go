package internal

import (
	"net/http"
	"time"
)

type FetchContext struct {
	Status    string
	Time      float64
	IsTimeout bool
}

func Fetch(uri string) FetchContext {
	start := time.Now()
	resp, err := http.Get(uri)
	if err != nil {
		return FetchContext{
			IsTimeout: true,
		}
	}
	defer resp.Body.Close()

	return FetchContext{
		Status:    resp.Status,
		Time:      float64(time.Since(start).Seconds()),
		IsTimeout: false,
	}

}
