package statistics_postgres_repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
)

//go:embed public_analytics.sql
var publicAnalyticsSQL string

func (r *StatisticsRepository) GetPublicAnalytics(ctx context.Context) (domain.PublicAnalytics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var analytics domain.PublicAnalytics
	if err := r.pool.QueryRow(ctx, publicAnalyticsSQL).Scan(&analytics); err != nil {
		return domain.PublicAnalytics{}, fmt.Errorf("select public task analytics: %w", err)
	}
	return analytics, nil
}
