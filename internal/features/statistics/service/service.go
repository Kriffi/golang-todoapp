package statistics_service

import (
	"context"
	"time"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetPublicAnalytics(ctx context.Context) (domain.PublicAnalytics, error)
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func (s *StatisticsService) GetPublicAnalytics(ctx context.Context) (domain.PublicAnalytics, error) {
	return s.statisticsRepository.GetPublicAnalytics(ctx)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
