package user

import (
	"encoding/json"
	"errors"
	"net/http"
	jsonUtils "server/internal/app/adapters/http-adapter/json"
	domain "server/internal/app/domain/user"
	"server/internal/pkg/logger"
)

type ProcessErrorResponse struct {
	ErrMsg     json.RawMessage
	HTTPStatus int
}

// ProcessServiceErrors return httpStatus (200, 400 ...) and errMsg according to custom error type.
func ProcessServiceErrors(err error, HandlerName string) ProcessErrorResponse {
	switch {
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		return ProcessErrorResponse{
			ErrMsg:     jsonUtils.ErrorAsJSON(err),
			HTTPStatus: http.StatusNotFound,
		}
	default:
		// fixme
		logger.Log.Error(HandlerName)
		return ProcessErrorResponse{
			ErrMsg:     jsonUtils.ErrorAsJSON(err),
			HTTPStatus: http.StatusInternalServerError,
		}
	}
}
