package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	error_handling "organization/error"
)

func CallAnotherService(jwtToken string, url string, input []byte, method string) ([]byte, error) {
	var buffer *bytes.Buffer
	if input != nil {
		buffer = bytes.NewBuffer(input)
	}
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, error_handling.InternalServerError
	}
	req.Header.Add("Authorization", jwtToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, error_handling.InternalServerError
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, error_handling.InternalServerError
	}
	if res.StatusCode != 200 {
		var customError error_handling.CustomError
		err = json.Unmarshal(body, &customError)
		if err != nil {
			return nil, error_handling.UnmarshalError
		}
		return nil,customError
	}
	return body, nil
}
