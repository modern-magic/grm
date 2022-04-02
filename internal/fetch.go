package internal

import (
	"net/http"
	"time"
)

type FetchState struct {
	status    string
	time      float64
	isTimeout bool
}

func fetch(uri string) FetchState {
	start := time.Now()
	resp, err := http.Get(uri)
	isTimeout := false
	status := resp.Status
	if err != nil {
		isTimeout = true
	}
	tc := time.Since(start).Seconds()
	meta := FetchState{
		status:    status,
		isTimeout: isTimeout,
		time:      tc,
	}
	return meta

}
