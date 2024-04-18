package response

type UserDetails struct {
	UserID      string `json:"userID" `
	Firstname   string `json:"firstname" `
	Lastname    string `json:"lastname" `
	Fullname    string `json:"fullname" `
	Username    string `json:"username" `
}