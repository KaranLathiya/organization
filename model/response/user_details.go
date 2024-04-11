package response

type UserID struct {
	UserID string `json:"userID"  validate:"required" `
}

type Username struct {
	Username string `json:"username,omitempty"  validate:"required" `
}

type UserDetails struct {
	UserID      string `json:"userID" `
	Firstname   string `json:"firstname" `
	Lastname    string `json:"lastname" `
	Fullname    string `json:"fullname" `
	Username    string `json:"username" `
	Email       *string `json:"email" `
	PhoneNumber *string `json:"phoneNumber" `
	CountryCode *string `json:"countryCode" `
	Privacy     string `json:"privacy" `
	CreatedAt   string `json:"createdAt" `
	UpdatedAt   *string `json:"updatedAt" `
}

type User struct {
	UserID      string `json:"userID" `
	Firstname   string `json:"firstname" `
	Lastname    string `json:"lastname" `
	Fullname    string `json:"fullname" `
	Username    string `json:"username" `
	Email       *string `json:"email" `
	PhoneNumber *string `json:"phoneNumber" `
	CountryCode *string `json:"countryCode" `
}

