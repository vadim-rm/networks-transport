package domain

import "time"

type Segment struct {
	Payload         string
	Username        string
	SendTime        time.Time
	SegmentSentTime time.Time
	Number          uint32
	TotalSegments   uint32
}
