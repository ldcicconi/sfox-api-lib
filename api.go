package sfoxapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// TODO: better errors
func (api *SFOXAPI) doRequest(action, path string, body interface{}, result interface{}) (bodyBytes []byte, statusCode int, err error) {
	// build request
	var reqBodyBytes []byte
	if body != nil && action == "POST" {
		reqBodyBytes, err = json.Marshal(body)
		if err != nil {
			return
		}
	}
	bytesReader := bytes.NewReader(reqBodyBytes)
	req, err := http.NewRequest(action, api.URL+path, bytesReader)
	if err != nil {
		return
	}
	// attach header for auth
	req.Header.Add("Authorization", "Bearer "+api.Key)
	// send request
	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	//read body
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return
	}
	//try to unmarshal
	if result != nil {
		err = json.Unmarshal(bodyBytes, result)
	}
	// return
	return
}
