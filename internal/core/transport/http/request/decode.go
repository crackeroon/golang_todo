package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/crackeroon/golang_todo/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json request body: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"validate request: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
