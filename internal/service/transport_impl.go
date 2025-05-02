package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"
	"transport/internal/domain"
	"transport/internal/metrics"
	"transport/internal/repository"
)

const segmentSize = 300 // 300 bytes
var scanPeriod = 1 * time.Second

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
	dataLink    repository.DataLink
	application repository.Application
	storage     map[time.Time]messageSegments
	mu          sync.Mutex
}

func NewTransportImpl(
	dataLink repository.DataLink,
	application repository.Application,
) *TransportImpl {
	return &TransportImpl{
		dataLink:    dataLink,
		application: application,
		storage:     make(map[time.Time]messageSegments),
		mu:          sync.Mutex{},
	}
}

func (s *TransportImpl) Run(ctx context.Context) {
	ticker := time.NewTicker(scanPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.scanStorage(ctx)
		}
	}
}

func (s *TransportImpl) scanStorage(ctx context.Context) {
	ticker := time.NewTicker(scanPeriod)
	defer ticker.Stop()

	for _, message := range s.storage {
		if time.Since(message.LastReceived) < 3*scanPeriod {
			continue
		}

		log.Printf("segment timeout for message")

		delete(s.storage, message.SendTime)

		err := s.application.SendMessage(ctx, domain.Message{
			Message:  "",
			SendTime: message.SendTime,
			Username: message.Username,
			Error:    "lost",
		})
		if err != nil {
			log.Printf("error sending message: %s", err.Error())
		}
	}
}

func (s *TransportImpl) ReceiveSegment(ctx context.Context, segment domain.Segment) error {
	log.Printf("received %s segment from data link layer", segment.Payload)

	message, err := s.processSegment(segment)
	if err != nil {
		log.Printf("error processing segment: %s", err.Error())
		return nil
	}
	if message.SegmentsLeft != 0 {
		return nil
	}

	averageTime := getAverageTime(message)
	fmt.Println("avg", averageTime)
	metrics.SegmentProcessMilliseconds.Set(averageTime)

	b := strings.Builder{}
	for _, segment := range message.Segments {
		b.WriteString(segment.Payload)
	}
	return s.application.SendMessage(ctx, domain.Message{
		Message:  b.String(),
		SendTime: message.SendTime,
		Username: message.Username,
	})
}

func getAverageTime(message messageSegments) float64 {
	var averageTime int64
	for _, segment := range message.Segments {
		averageTime += segment.ReceivedAt.Sub(message.SegmentsSentTime).Milliseconds()
	}

	averageTime /= int64(len(message.Segments))

	return float64(averageTime) / 1000
}

func (s *TransportImpl) processSegment(segment domain.Segment) (messageSegments, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message, ok := s.storage[segment.SendTime]
	if !ok {
		return messageSegments{}, domain.ErrUnknownMessage
	}
	message.SegmentsLeft--
	message.Segments[segment.Number-1] = messageSegmentsSegment{
		Payload:    segment.Payload,
		ReceivedAt: time.Now(),
	}
	message.LastReceived = time.Now()
	s.storage[segment.SendTime] = message

	if message.SegmentsLeft == 0 {
		delete(s.storage, message.SendTime)
	}

	return message, nil
}

func (s *TransportImpl) SendMessage(ctx context.Context, message domain.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rawSegments := splitMessage(message.Message, segmentSize)

	s.storage[message.SendTime] = messageSegments{
		SegmentsLeft:     uint32(len(rawSegments)),
		Username:         message.Username,
		SendTime:         message.SendTime,
		Segments:         make([]messageSegmentsSegment, len(rawSegments)),
		SegmentsSentTime: time.Now(),
		LastReceived:     time.Now(),
	}

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
