package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/crackeroon/golang_todo/internal/core/logger"
	core_http_request "github.com/crackeroon/golang_todo/internal/core/transport/http/request"
	core_http_response "github.com/crackeroon/golang_todo/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskId path value",
		)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task domain")
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}
