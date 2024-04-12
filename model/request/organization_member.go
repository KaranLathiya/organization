package request

type UpdateMemberRole struct {
	Role           string `json:"Role" validate:"required|in:admin,editor,viewer"`
	OrganizationID string `json:"organizationID" validate:"required"`
	MemberID       string `json:"memberID" validate:"required"`
}

type OrganizationID struct {
	OrganizationID string `json:"organizationID" validate:"required"`
}

type RemoveMemberFromOrganization struct {
	OrganizationID string `json:"organizationID" validate:"required"`
	MemberID       string `json:"memberID" validate:"required"`
}