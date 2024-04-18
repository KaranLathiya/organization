package response

type InvitationDetails struct {
	ID             string `json:"id" `
	Role           string `json:"role" `
	OrganizationID string `json:"organizationID"`
	InvitedBy      string `json:"invitedBy"`
	InvitedAt      string `json:"invitedAt"`
}
