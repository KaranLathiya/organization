package response

type MicrosoftAuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MicrosoftAuthURL struct {
	AuthURL string `json:"authUrl" validate:"required" `
}	