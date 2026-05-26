package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/crackeroon/golang_todo/internal/core/domain"
	core_logger "github.com/crackeroon/golang_todo/internal/core/logger"
	core_http_request "github.com/crackeroon/golang_todo/internal/core/transport/http/request"
	core_http_response "github.com/crackeroon/golang_todo/internal/core/transport/http/response"
	core_http_types "github.com/crackeroon/golang_todo/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type PatchTaskResponse TaskDTOResponse

func (h *PatchTaskRequest) Validate() error {
	if h.Title.Set {
		if h.Title.Value == nil {
			return fmt.Errorf("`Title` can not be NULL")
		}
		titleLen := len([]rune(*h.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 characters")
		}
	}
	if h.Description.Set {
		if h.Description.Value != nil {
			descriptionLen := len([]rune(*h.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 characters")
			}
		}
	}

	if h.Completed.Set {
		if h.Completed.Value == nil {
			return fmt.Errorf("`Completed` can not be NULL")
		}
	}
	return nil
}

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	var req PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}
	taskPatch := taskPatchFromRequest(req)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
