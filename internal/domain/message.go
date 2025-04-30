package domain

import "time"

type Message struct {
	Message  string
	SendTime time.Time
	Username string
	Error    string
}
