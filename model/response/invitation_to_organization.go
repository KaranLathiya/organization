package response

type InvitationDetails struct {
	ID             string `json:"ID" `
	Role           string `json:"Role" `
	OrganizationID string `json:"organizationID"`
	InvitedBy      string `json:"invitedBy"`
	InvitedAt      string `json:"invitedAt"`
}
