package handlers

import (
	"context"
	"encoding/json"
	"time"
	"transport/internal/domain"
	"transport/internal/service"
)

type Segments struct {
	transport service.Transport
}

func NewDataLinkSegments(transport service.Transport) *Segments {
	return &Segments{transport: transport}
}

type kafkaSegmentMessage struct {
	Payload       string    `json:"payload"`
	Username      string    `json:"username"`
	SendTime      time.Time `json:"sendTime"`
	Number        uint32    `json:"segmentNumber"`
	TotalSegments uint32    `json:"totalSegments"`
}

func (h *Segments) Receive(ctx context.Context, messageValue []byte) error {
	var message kafkaSegmentMessage
	err := json.Unmarshal(messageValue, &message)
	if err != nil {
		return err
	}

	return h.transport.ReceiveSegment(ctx, domain.Segment{
		Payload:       message.Payload,
		Username:      message.Username,
		SendTime:      message.SendTime,
		Number:        message.Number,
		TotalSegments: message.TotalSegments,
	})
}
