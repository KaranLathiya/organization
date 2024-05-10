package request 

type Body struct {
	Content string `json:"content"`
}

type ChannnelMessage struct {
	Body Body `json:"body"`
}