package sfoxapi

import (
	// "crypto/tls"
	"net/http"
	"time"
)

type SFOXAPI struct {
	Key        string
	HttpClient *http.Client
	URL        string
}

func NewSFOXAPI(apiKey string) *SFOXAPI {
	tr := http.Transport{
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
