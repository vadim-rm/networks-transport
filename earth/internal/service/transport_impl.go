package service

import (
	"context"
	"fmt"
	"math"
	"time"
	"transport/internal/domain"
	"transport/internal/repository"
)

const segmentSize = 300 // 300 bytes

type messageSegments struct {
	SegmentsLeft     uint32
	Username         string
	SendTime         time.Time
	LastReceived     time.Time
	SegmentsSentTime time.Time
	Segments         []messageSegmentsSegment
}

type messageSegmentsSegment struct {
	Payload    string
	ReceivedAt time.Time
}

type TransportImpl struct {
	dataLink repository.DataLink
}

func NewTransportImpl(
	dataLink repository.DataLink,
) *TransportImpl {
	return &TransportImpl{
		dataLink: dataLink,
	}
}

func (s *TransportImpl) SendMessage(ctx context.Context, message domain.Message) error {
	rawSegments := splitMessage(message.Message, segmentSize)

	for i, segment := range rawSegments {
		err := s.dataLink.Send(ctx, domain.Segment{
			Payload:       segment,
			Username:      message.Username,
			SendTime:      message.SendTime,
			Number:        uint32(i + 1),
			TotalSegments: uint32(len(rawSegments)),
		})

		if err != nil {
			return fmt.Errorf("error sending segment: %w", err)
		}
	}

	return nil
}

func splitMessage(payload string, segmentSize int) []string {
	result := make([]string, 0)

	length := len(payload) // длина в байтах
	segmentCount := int(math.Ceil(float64(length) / float64(segmentSize)))

	for i := 0; i < segmentCount; i++ {
		result = append(result, payload[i*segmentSize:min((i+1)*segmentSize, length)]) // срез делается также по байтам
	}

	return result
}
