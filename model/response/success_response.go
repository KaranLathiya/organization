package response

type SuccessResponse struct {
	Message string `json:"message"`
}

type OrganizationID struct {
	OrganizationID string `json:"organizationID"  `
}

type JWTToken struct {
	JWTToken string `json:"jwtToken"  `
}