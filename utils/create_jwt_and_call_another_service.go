package utils

import (
	"encoding/json"
	"organization/internal"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
)

func CreateJWTAndCallAnotherService(allMemberIDs []string,audience string, subject string) (map[string]response.UserDetails, error) {
	jwtToken, err := middleware.CreateJWT(audience, subject)
	if err != nil {
		return nil, err
	}
	userIDs := request.UserIDs{
		UserIDs: allMemberIDs,
	}

	userIDsByte, _ := json.MarshalIndent(userIDs, "", "  ")
	body, err := internal.CallAnotherService(jwtToken, "http://localhost:8000/users/details", userIDsByte, "POST")
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
