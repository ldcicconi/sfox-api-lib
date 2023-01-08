package sfoxapi

import (
	"errors"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Message string `json:"error"`
}

func (client *Client) doRequest(action, path string, body, result any, logRawData, requiresAuth bool) (err error) {
	var (
		errResponse errorResponse
		url         = fmt.Sprintf("%s%s", client.url, path)
	)

	switch action {
	case http.MethodPost:
		_, err := client.restyClient.R().
			SetBody(body).
			SetResult(result).
			SetError(&errResponse).
			Post(url)
		return err
	case http.MethodGet:
		_, err := client.restyClient.R().
			SetResult(result).
			SetError(&errResponse).
			Get(url)
		return err
	default:
		return errors.New("unsupported method")
	}
}
