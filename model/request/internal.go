package request

type UserIDs struct {
	UserIDs []string `json:"userIDs" validate:"required"`
}

