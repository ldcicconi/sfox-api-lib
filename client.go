package sfoxapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	url         string
	restyClient *resty.Client
}

type Option func(*Client)

func NewClient(apiKey string, options ...Option) *Client {
	client := Client{
		url:         "https://api.sfox.com",
		restyClient: resty.New().SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)),
	}

	for _, option := range options {
		option(&client)
	}

	return &client
}

type Logger interface {
	Infof(format string, v ...interface{})
}

func WithLogger(logger Logger) Option {
	return func(client *Client) {
		client.restyClient.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Infof("url: %v raw_response: %v status_code: %d", r.Request.URL, string(r.Body()), r.StatusCode())
			return nil
		})
	}
}
