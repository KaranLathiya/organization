package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"organization/constant"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
)

func CallAnotherService(jwtToken string, url string, input []byte, method string) ([]byte, error) {
	var buffer *bytes.Buffer
	if input != nil {
		buffer = bytes.NewBuffer(input)
	}
	var req *http.Request
	var err error
	if method == http.MethodGet || method == http.MethodDelete {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, buffer)
	}
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
	if res.StatusCode != http.StatusOK {
		var customError error_handling.CustomError
		err = json.Unmarshal(body, &customError)
		if err != nil {
			return nil, error_handling.UnmarshalError
		}
		return nil, customError
	}
	return body, nil
}

func CreateJWTAndCallUserServiceForUserDetails(allMemberIDs []string, audience string, subject string) (map[string]response.UserDetails, error) {
	jwtToken, err := CreateJWT(audience, subject)
	if err != nil {
		return nil, err
	}
	userIDs := request.UserIDs{
		UserIDs: allMemberIDs,
	}
	userIDsByte, err := json.MarshalIndent(userIDs, "", "  ")
	if err != nil {
		return nil, error_handling.MarshalError
	}
	body, err := CallAnotherService(jwtToken, constant.USER_SERVICE_BASE_URL+"internal/users/details", userIDsByte, http.MethodPost)
	if err != nil {
		return nil, err
	}
	userDetails := make(map[string]response.UserDetails)
	err = json.Unmarshal(body, &userDetails)
	if err != nil {
		return nil, err
	}
	return userDetails, nil
}
