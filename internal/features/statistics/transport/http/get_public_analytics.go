package statistics_transport_http

import (
	"net/http"

	core_logger "github.com/Kriffi/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Kriffi/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func (h *StatisticsHTTPHandler) GetPublicAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	response := core_http_response.NewHTTPResponseHandler(log, w)
	analytics, err := h.statisticsService.GetPublicAnalytics(ctx)
	if err != nil {
		log.Error("failed to get public task analytics", zap.Error(err))
		response.JSONResponse(map[string]string{"message": "Analytics temporarily unavailable"}, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	response.JSONResponse(analytics, http.StatusOK)
}
