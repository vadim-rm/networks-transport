package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"transport/internal/domain"
)

type HTTPApplication struct {
	baseUrl string
}

func NewHTTPApplication(baseUrl string) *HTTPApplication {
	return &HTTPApplication{baseUrl: baseUrl}
}

type httpMessage struct {
	Message  string    `json:"message"`
	SendTime time.Time `json:"sendTime"`
	Username string    `json:"username"`
	Error    string    `json:"error"`
}

func (r *HTTPApplication) SendMessage(ctx context.Context, message domain.Message) error {
	log.Printf("sending %s to application level", message.Message)

	messageValue, err := json.Marshal(httpMessage{
		Message:  message.Message,
		SendTime: message.SendTime,
		Username: message.Username,
		Error:    message.Error,
	})
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Post(
		fmt.Sprintf("%s/receive", r.baseUrl),
		"application/json",
		bytes.NewBuffer(messageValue),
	)
	return err
}
