package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"transport/internal/domain"
	"transport/internal/metrics"
	"transport/internal/repository"
)

var scanPeriod = 1 * time.Second

type messageSegments struct {
	SegmentsLeft uint32
	Username     string
	SendTime     time.Time
	LastReceived time.Time
	Segments     []messageSegmentsSegment
}

type messageSegmentsSegment struct {
	Payload    string
	ReceivedAt time.Time
}

type TransportImpl struct {
	application repository.Application
	storage     map[time.Time]messageSegments
	mu          sync.Mutex
}

func NewTransportImpl(
	application repository.Application,
) *TransportImpl {
	return &TransportImpl{
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
		averageTime += segment.ReceivedAt.Sub(message.SendTime).Milliseconds()
	}

	averageTime /= int64(len(message.Segments))

	return float64(averageTime) / 1000
}

func (s *TransportImpl) processSegment(segment domain.Segment) (messageSegments, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message, ok := s.storage[segment.SendTime]
	if !ok {
		message = messageSegments{
			SegmentsLeft: segment.TotalSegments,
			Username:     segment.Username,
			SendTime:     segment.SendTime,
			Segments:     make([]messageSegmentsSegment, segment.TotalSegments),
			LastReceived: time.Now(),
		}
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
