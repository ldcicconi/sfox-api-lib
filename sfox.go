package sfoxapi

import (
	// "crypto/tls"
	"net/http"
	"time"
)

type SFOXAPI struct {
	Key          string
	HttpClient   *http.Client
	URL          string
	ErrorMonitor *Monitor
}

func NewSFOXAPI(apiKey string, monitor *Monitor) *SFOXAPI {
	tr := http.Transport{
		MaxIdleConnsPerHost: 20,
	}
	return &SFOXAPI{
		Key: apiKey,
		HttpClient: &http.Client{
			Timeout:   time.Second * 5,
			Transport: &tr,
		},
		URL: "https://api.sfox.com",
	}
}
