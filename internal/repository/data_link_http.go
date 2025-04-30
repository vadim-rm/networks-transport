package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"transport/internal/domain"
)

type HTTPDataLink struct {
	baseUrl string
}

func NewHTTPDataLink(baseUrl string) *HTTPDataLink {
	return &HTTPDataLink{baseUrl: baseUrl}
}

type httpSegmentMessage struct {
	Payload       string    `json:"payload"`
	Username      string    `json:"username"`
	SendTime      time.Time `json:"sendTime"`
	Number        uint32    `json:"number"`
	TotalSegments uint32    `json:"totalSegments"`
}

func (r *HTTPDataLink) Send(ctx context.Context, segment domain.Segment) error {
	return nil
	message, err := json.Marshal(httpSegmentMessage{
		Payload:       segment.Payload,
		Username:      segment.Username,
		SendTime:      segment.SendTime,
		Number:        segment.Number,
		TotalSegments: segment.TotalSegments,
	})
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Post(
		fmt.Sprintf("%s/segment", r.baseUrl),
		"application/json",
		bytes.NewBuffer(message),
	)
	return err
}
