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
	if err != nil {
		return FetchState{
			isTimeout: true,
		}
	}
	defer resp.Body.Close()

	return FetchState{
		status:    resp.Status,
		time:      float64(time.Since(start).Seconds()),
		isTimeout: false,
	}

}
