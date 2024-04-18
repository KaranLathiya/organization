package response

type AllOrganizationDetailsOfUser struct {
	UserID        string         `json:"userID"`
	Organizations []Organization `json:"organizations"`
}

type Organization struct {
	OrganizationID      string               `json:"organizationID"`
	OwnerID             string               `json:"ownerID"`
	Name                string               `json:"name"`
	Privacy             string               `json:"privacy"`
	CreatedAt           string               `json:"createdAt" `
	UpdatedAt           *string              `json:"updatedAt" `
	UpdatedBy           *string              `json:"updatedBy" `
	OrganizationMembers *[]*OrganizationMember `json:"organizationMembers"`
}

type OrganizationMember struct {
	UserID    string `json:"userID" `
	Role      string `json:"role" `
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Fullname  string `json:"fullname" `
	Username  string `json:"username" `
}
