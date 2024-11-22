package room

import "github.com/prometheus/client_golang/prometheus"

type ServerMetrics struct {
	ActiveConnections prometheus.Gauge
	Requests          *prometheus.CounterVec
	Responses         *prometheus.CounterVec
	ResponseDuration  *prometheus.HistogramVec
}

func NewServerMetrics() ServerMetrics {
	sm := ServerMetrics{
		ActiveConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "sio_active_connections",
				Help: "Current number of active socket.io connections.",
			},
		),
		Requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sio_requests_total",
				Help: "Total number of processed socket.io requests",
			},
			[]string{"event"},
		),
		Responses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sio_responses_total",
				Help: "Total number of socket.io responses",
			},
			[]string{"event"},
		),
		ResponseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "sio_response_duration_seconds",
				Help:    "Histogram of response durations for socket.io requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"event"},
		),
	}
	prometheus.MustRegister(sm.ActiveConnections, sm.Requests, sm.Responses, sm.ResponseDuration)
	return sm
}
