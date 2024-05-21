package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"organization/config"
	"organization/constant"
	error_handling "organization/error"
	"organization/model/response"
	"organization/utils"
	"strings"
)

func (c *UserController) MicrosoftAuth(w http.ResponseWriter, r *http.Request) {
	scopes := "https://graph.microsoft.com/Chat.Read https://graph.microsoft.com/Chat.ReadWrite https://graph.microsoft.com/Chat.ReadBasic https://graph.microsoft.com/ChatMessage.Read https://graph.microsoft.com/ChatMessage.Send https://graph.microsoft.com/Channel.ReadBasic.All https://graph.microsoft.com/ChannelMessage.Send https://graph.microsoft.com/User.Read https://graph.microsoft.com/email https://graph.microsoft.com/openid https://graph.microsoft.com/profile"
	authURL := "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + config.ConfigVal.MicrosoftAuth.ClientID + "&response_type=code&redirect_uri=" + config.ConfigVal.MicrosoftAuth.RedirectURI + "&response_mode=query&scope=" + scopes
	microsoftAuthURL := response.MicrosoftAuthURL{AuthURL: authURL}
	utils.SuccessMessageResponse(w, http.StatusOK, microsoftAuthURL)
}

func (c *UserController) GetMicrosoftTokens(w http.ResponseWriter, r *http.Request) {
	data := url.Values{}
	data.Set("code", r.FormValue("code"))
	data.Set("client_id", config.ConfigVal.MicrosoftAuth.ClientID)
	data.Set("redirect_uri", config.ConfigVal.MicrosoftAuth.RedirectURI)
	data.Set("client_secret", config.ConfigVal.MicrosoftAuth.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", "https://graph.microsoft.com/Chat.Read https://graph.microsoft.com/Chat.ReadWrite https://graph.microsoft.com/Chat.ReadBasic https://graph.microsoft.com/ChatMessage.Read https://graph.microsoft.com/ChatMessage.Send https://graph.microsoft.com/Channel.ReadBasic.All https://graph.microsoft.com/ChannelMessage.Send https://graph.microsoft.com/User.Read https://graph.microsoft.com/email https://graph.microsoft.com/openid https://graph.microsoft.com/profile offline_access")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, constant.MICROSOFT_AUTH_URL, strings.NewReader(data.Encode()))
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.InternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil || res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusInternalServerError || res.StatusCode == http.StatusUnauthorized {
		error_handling.ErrorMessageResponse(w, error_handling.InternalServerError)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.InternalServerError)
		return
	}

	var microsoftAuthToken response.MicrosoftAuthToken
	err = json.Unmarshal(resBody, &microsoftAuthToken)
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.UnmarshalError)
		return
	}

	err = c.repo.StoreMicrosoftRefreshToken(microsoftAuthToken.RefreshToken)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.TOKENS_CREATED_AND_STORED})
}
