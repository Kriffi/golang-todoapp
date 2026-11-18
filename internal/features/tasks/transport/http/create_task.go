package tasks_transport_http

import (
	"net/http"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
	core_logger "github.com/Kriffi/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Kriffi/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Kriffi/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	DueDate      *string `json:"due_date" format:"date" example:"2026-12-01"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

// CreateTask godoc
// @Summary Создать задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Создать задачу"
// @Success 201 {object} TaskDTOResponse
// @Failure 400 {object} core_http_response.ErrorResponseBody
// @Failure 404 {object} core_http_response.ErrorResponseBody
// @Failure 500 {object} core_http_response.ErrorResponseBody
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain.DueDate = request.DueDate

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := taskDTOFromDomain(taskDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
