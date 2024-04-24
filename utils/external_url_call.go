package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	error_handling "organization/error"
)

func ExternalURLCall(method string, url string, bodyDataRequest interface{}, bodyDataResponse map[string]interface{}) (map[string]interface{}, error) {
	bodyDataByte, _ := json.MarshalIndent(bodyDataRequest, "", " ")
	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyDataByte))
	}
	if err != nil {
		return nil, error_handling.InternalServerError
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusNotFound ||  res.StatusCode == http.StatusInternalServerError {
		return nil, error_handling.InternalServerError
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBody, &bodyDataResponse)
	return bodyDataResponse, err
}
