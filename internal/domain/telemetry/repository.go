package telemetry

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type RequestMetricRepository interface {
	GetSummary(context.Context) ([]*RequestLogMetricSummaryType, error)
	Store(context.Context, *RequestLogType) error
}

type requestMetricRepository struct {
	db *sqlx.DB
}

// NewRepository return an instance of requestMetric containing interface to DB layer
func NewRepository(db *sqlx.DB) RequestMetricRepository {
	return &requestMetricRepository{db: db}
}

var pgSquirrel = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (repo *requestMetricRepository) GetSummary(ctx context.Context) (res []*RequestLogMetricSummaryType, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("endpoint", "round(min(latency)::numeric, 2) min_latency", "round(avg(latency)::numeric, 2) avg_latency", "round(max(latency)::numeric, 2) max_latency",
		"count(*) req_count", "count(distinct user_agent) unique_ua").
		From("request_logs").
		GroupBy("endpoint").ToSql()

	res = []*RequestLogMetricSummaryType{}
	rows, err := repo.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	}

	for rows.Next() {
		col := &RequestLogMetricSummaryType{}

		if err = rows.StructScan(col); err != nil {
			logger.Error().Err(err).Msg("failed to map result")
			return
		}

		res = append(res, col)
	}

	return
}

func (repo *requestMetricRepository) Store(ctx context.Context, payload *RequestLogType) (err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Insert("request_logs").Columns("endpoint", "latency", "user_agent", "requested_at").
		Values(payload.Endpoint, payload.Latency, payload.UserAgent, payload.RequestedAt).ToSql()

	fmt.Println(stmt, args)

	_, err = repo.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to save request log metrics")
		return
	}

	return
}
