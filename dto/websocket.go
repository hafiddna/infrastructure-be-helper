package dto

import "time"

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ChatPayload struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type NotificationPayload struct {
	Icon     string    `json:"icon"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	NotifyAt time.Time `json:"notify_at"`
}
