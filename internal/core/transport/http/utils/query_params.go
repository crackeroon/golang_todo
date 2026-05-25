package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/crackeroon/golang_todo/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"parameter='%s' by key='%s' is not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return &val, nil
}

/*func GetTimeQueryParam(r *http.Request, key string) time.Time {

}
*/
