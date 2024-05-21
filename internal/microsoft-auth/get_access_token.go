package microsoftauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"organization/config"
	"organization/constant"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
	"strings"
)

func GetAccessTokenUsingRefreshToken(refreshToken string) (response.MicrosoftAuthToken, error) {
	var microsoftAuthToken response.MicrosoftAuthToken
	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", config.ConfigVal.MicrosoftAuth.ClientID)
	data.Set("redirect_uri", config.ConfigVal.MicrosoftAuth.RedirectURI)
	data.Set("client_secret", config.ConfigVal.MicrosoftAuth.ClientSecret)
	data.Set("grant_type", "refresh_token")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, constant.MICROSOFT_AUTH_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return microsoftAuthToken, error_handling.InternalServerError
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil || res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusInternalServerError || res.StatusCode == http.StatusUnauthorized {
		return microsoftAuthToken, error_handling.InternalServerError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return microsoftAuthToken, err
	}

	err = json.Unmarshal(resBody, &microsoftAuthToken)
	if err != nil {
		return microsoftAuthToken, error_handling.UnmarshalError
	}

	return microsoftAuthToken, nil
}

func SendMessageToChannel(message string, accessToken string) error {
	apiURL := constant.MICROSOFT_GRAPH_API_BASE_URL + "/v1.0/teams/" + config.ConfigVal.MicrosoftAuth.TeamID + "/channels/" + config.ConfigVal.MicrosoftAuth.ChannelID + "/messages"
	client := &http.Client{}

	messageBody := request.ChannnelMessage{
		Body: request.Body{
			Content: message,
		},
	}
	jsonBody, err := json.MarshalIndent(messageBody, "", " ")
	if err != nil {
		return error_handling.MarshalError	
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return error_handling.InternalServerError
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil || res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusInternalServerError || res.StatusCode == http.StatusUnauthorized {
		return error_handling.InternalServerError
	}

	return nil
}
