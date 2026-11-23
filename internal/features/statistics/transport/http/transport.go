package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
	core_http_server "github.com/Kriffi/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetPublicAnalytics(ctx context.Context) (domain.PublicAnalytics, error)
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
		{
			Method:  http.MethodGet,
			Path:    "/statistics/public",
			Handler: h.GetPublicAnalytics,
		},
	}
}
