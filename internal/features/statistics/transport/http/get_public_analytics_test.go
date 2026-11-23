package statistics_transport_http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
	core_logger "github.com/Kriffi/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

func TestPublicAnalyticsOnlyReturnsAggregateSections(t *testing.T) {
	handler := NewStatisticsHTTPHandler(analyticsStub{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/statistics/public?user_id=7", nil)
	request = request.WithContext(core_logger.ToContext(request.Context(), &core_logger.Logger{Logger: zap.NewNop()}))
	handler.GetPublicAnalytics(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	expected := []string{"overview", "status", "aging", "throughput", "due_soon"}
	if len(body) != len(expected) {
		t.Fatalf("unexpected response fields: %v", body)
	}
	for _, name := range expected {
		if _, exists := body[name]; !exists {
			t.Errorf("missing %s", name)
		}
	}
}

type analyticsStub struct{}

func (analyticsStub) GetStatistics(context.Context, *int, *time.Time, *time.Time) (domain.Statistics, error) {
	return domain.Statistics{}, nil
}

func (analyticsStub) GetPublicAnalytics(context.Context) (domain.PublicAnalytics, error) {
	return domain.PublicAnalytics{}, nil
}
