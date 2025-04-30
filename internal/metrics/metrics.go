package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	SegmentProcessMilliseconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "segment_process_avg_milliseconds_gauge",
		Help: "Среднее время обработки сегмента",
	})
)

func init() {
	prometheus.MustRegister(SegmentProcessMilliseconds)
}
