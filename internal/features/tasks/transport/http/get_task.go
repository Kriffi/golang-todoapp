package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Kriffi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Kriffi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Kriffi/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получение задачи
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} core_http_response.ErrorResponseBody
// @Failure 404 {object} core_http_response.ErrorResponseBody
// @Failure 500 {object} core_http_response.ErrorResponseBody
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(task))

	responseHandler.JSONResponse(response, http.StatusOK)
}
