package request

type InvitationToOrganization struct {
	Role string `json:"Role" validate:"required|in:admin,editor,viewer"`
	OrganizationID string `json:"organizationID" validate:"required"`
	Invitee string `json:"invitee" validate:"required"`
}
