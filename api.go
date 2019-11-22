package sfoxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ErrJsonUnmarshalling = fmt.Errorf("error unmarshalling JSON body")

type SFOXErrorMsg struct {
	Error string `json:"error"`
}

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
	req.Header.Add("Content-Type", "application/json")
	// send request
	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	// read body
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// check status code
	if statusCode >= 400 {
		// try to marshal into error message
		var errorMsg SFOXErrorMsg
		err = json.Unmarshal(bodyBytes, errorMsg)
		if err != nil {
			err = ErrJsonUnmarshalling
			return
		} else {
			err = fmt.Errorf("statusCode: %d error_message: %s", statusCode, errorMsg.Error)
		}
	}
	// try to unmarshal
	if result != nil {
		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			err = ErrJsonUnmarshalling
			return
		}
	}
	// return
	return
}
