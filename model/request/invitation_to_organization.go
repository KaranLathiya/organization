package request

type InvitationToOrganization struct {
	Role           string `json:"role" validate:"required|in:admin,editor,viewer"`
	OrganizationID string `json:"organizationID" validate:"required"`
	Invitee        string `json:"invitee" validate:"required"`
}

type RespondToInvitation struct {
	InvitationAccept *bool  `json:"invitationAccept" validate:"required|isBool" `
	OrganizationID   string `json:"organizationID" validate:"required"`
}
