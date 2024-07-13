package telemetry

import "time"

// RequestMetricPayload represent payload fetch from middleware by end of request
type RequestMetricPayload struct {
	Endpoint    string
	Latency     float64
	UserAgent   string
	RequestedAt time.Time
}

// RequestMetric represent domain response for Telemetry Request Metric entity
type RequestMetric struct {
	Endpoint        string  `json:"endpoint" example:"GET /ponds"`
	MinLatency      float64 `json:"min_latency" example:"1"`
	AvgLatency      float64 `json:"avg_latency" example:"10"`
	MaxLatency      float64 `json:"max_latency" example:"100"`
	UniqueUserAgent int64   `json:"unique_user_agent" example:"5"`
	Count           int64   `json:"count" example:"25"`
}

// ListRequestMetricResponse represent domain response for bulk Requet Metric entities
type ListRequestMetricResponse struct {
	Metrics []*RequestMetric `json:"request_metrics"`
}
