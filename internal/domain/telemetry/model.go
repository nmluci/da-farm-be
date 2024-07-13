package telemetry

import "time"

type RequestLogType struct {
	ID          int64     `db:"id"`
	Endpoint    string    `db:"endpoint"`
	Latency     float64   `db:"latency"`
	UserAgent   string    `db:"user_agent"`
	RequestedAt time.Time `db:"requested_at"`
}

type RequestLogMetricSummaryType struct {
	Endpoint   string  `db:"endpoint"`
	AvgLatency float64 `db:"avg_latency"`
	MinLatency float64 `db:"min_latency"`
	MaxLatency float64 `db:"max_latency"`
	Count      int64   `db:"req_count"`
	UniqueUA   int64   `db:"unique_ua"`
}
