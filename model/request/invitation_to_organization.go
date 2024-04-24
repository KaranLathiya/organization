package request

type InvitationToOrganization struct {
	Role           string `json:"role" validate:"required|in:admin,editor,viewer"`
	OrganizationID string `json:"organizationID" validate:"required"`
	Invitee        string `json:"invitee" validate:"required"`
}

type RespondToInvitation struct {
	Respond        string   `json:"Respond" validate:"required|in:accept,reject"`
	OrganizationID string `json:"organizationID" validate:"required"`
}
