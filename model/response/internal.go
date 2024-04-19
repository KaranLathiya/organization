package response

type UserDetails struct {
	UserID    string `json:"userID" `
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Fullname  string `json:"fullname" `
	Username  string `json:"username" `
}

type OrganizationListOfUser struct {
	UserID       string         `json:"userIDs" `
	Organizations *[]OrganizationInfoOfUser `json:"organizations"`
}

type OrganizationInfoOfUser struct {
	Role           string `json:"role"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationID"`
}
