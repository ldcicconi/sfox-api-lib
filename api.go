package sfoxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SFOXErrorMsg struct {
	errorMessage string `json:"error"`
}

func (e *SFOXErrorMsg) Error() string {
	return e.errorMessage
}

var ErrHttp = fmt.Errorf("error with http")
var ErrReading = fmt.Errorf("error readind body")
var ErrJsonMarshalling = fmt.Errorf("error marshalling JSON")
var ErrJsonUnmarshalling = fmt.Errorf("error unmarshalling JSON")

// TODO: better errors
func (api *SFOXAPI) doRequest(action, path string, body interface{}, result interface{}, logRawData bool) (bodyBytes []byte, statusCode int, err error) {
	// build request
	var reqBodyBytes []byte
	if body != nil && action == "POST" {
		reqBodyBytes, err = json.Marshal(body)
		if err != nil {
			err = ErrJsonMarshalling
			return
		}
	}
	bytesReader := bytes.NewReader(reqBodyBytes)
	req, err := http.NewRequest(action, api.URL+path, bytesReader)
	if err != nil {
		err = ErrHttp
		return
	}
	// attach header for auth
	req.Header.Add("Authorization", "Bearer "+api.Key)
	req.Header.Add("Content-Type", "application/json")
	// send request
	resp, err := api.HttpClient.Do(req)
	if err != nil && resp == nil {
		err = ErrHttp
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	// read body
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = ErrReading
		return
	}
	if logRawData {
		fmt.Println(string(bodyBytes))
	}
	// check status code
	if statusCode >= 400 {
		// try to marshal into error message
		var errorMsg SFOXErrorMsg
		err = json.Unmarshal(bodyBytes, errorMsg)
		if err != nil {
			err = &ResponseBodyError{
				Underlying:   nil,
				ResponseBody: string(bodyBytes),
			}
			return
		} else {
			err = fmt.Errorf("statusCode: %d error_message: %s", statusCode, errorMsg.Error)
		}
	}
	// try to unmarshal
	if result != nil {
		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			err = &ResponseBodyError{
				Underlying:   ErrJsonUnmarshalling,
				ResponseBody: string(bodyBytes),
			}
		}
	}
	// return
	return
}
