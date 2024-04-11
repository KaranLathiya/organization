package request

type CreateOrganization struct {
	Name    string `json:"name" validate:"required|min_len:2|max_len:50" `
	Privacy string `json:"Privacy" validate:"required|in:public,private"`
}

type UpdateOrganizationDetails struct {
	OrganizationID string `json:"organizationID" validate:"required"`
	Name           string `json:"name" validate:"required|min_len:2|max_len:50" `
	Privacy        string `json:"Privacy" validate:"required|in:public,private"`
}
