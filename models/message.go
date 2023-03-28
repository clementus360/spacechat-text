package models

type Message struct {
	ID string `json:"id"`
	Payload string `json:"payload"`
	ChatId string `json:"chatId"`
	Receiver string `json:"receiver"`
	Sender string `json:"sender"`
	Timestamp string `json:"timestamp"`
}
