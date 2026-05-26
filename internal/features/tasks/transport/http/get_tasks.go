package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/crackeroon/golang_todo/internal/core/logger"
	core_http_request "github.com/crackeroon/golang_todo/internal/core/transport/http/request"
	core_http_response "github.com/crackeroon/golang_todo/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user limit offset query params")
		return
	}
	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks")
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (userId, limit, offset *int, err error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userId, err = core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `userId` query param: %w", err)
	}

	limit, err = core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `limit` query param: %w", err)
		/*return nil, nil, fmt.Errorf("get limit query param: %v: %w", err, core_errors.ErrInvalidArgument)*/
	}

	offset, err = core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `offset` query param: %w", err)
		/*		return nil, nil, fmt.Errorf("get offset query param: %v: %w", err, core_errors.ErrInvalidArgument)*/
	}

	return userId, limit, offset, nil
}
