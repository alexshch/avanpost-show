package entity

type Message struct {
	Message string `json:"message"`
}

func NewMessage(message string) *Message {
	return &Message{message}
}
