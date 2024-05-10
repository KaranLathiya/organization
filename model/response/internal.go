package response

type UserDetails struct {
	UserID    string `json:"userID" `
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Fullname  string `json:"fullname" `
	Username  string `json:"username" `
}

type OrganizationCreatorDetails struct {
	Fullname    string  `json:"fullname" `
	Email       *string `json:"email" `
	PhoneNumber *string `json:"phoneNumber" `	
	CountryCode *string `json:"countryCode" `	
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
