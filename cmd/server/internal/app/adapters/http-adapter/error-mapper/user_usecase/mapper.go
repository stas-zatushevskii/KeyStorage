package user_usecase

import (
	"errors"
	"net/http"
	domain "server/internal/app/domain/user"
)

// Process return httpStatus (200, 400 ...) and ErrMsg according to custom error type.
func Process(err error) (int, string) {
	var (
		httpStatus      int
		responseMessage string
	)

	switch {
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		httpStatus = http.StatusConflict
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrPasswordMismatch):
		httpStatus = http.StatusUnauthorized
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrTokenNotValid):
		httpStatus = http.StatusUnauthorized
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrInvalidRefreshToken):
		httpStatus = http.StatusUnauthorized
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrRefreshTokenNotFound):
		httpStatus = http.StatusBadRequest
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrTokenRevoked):
		httpStatus = http.StatusForbidden
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrUserNotFound):
		httpStatus = http.StatusForbidden
		responseMessage = err.Error()
	default:
		httpStatus = http.StatusInternalServerError
		responseMessage = err.Error()
	}
	return httpStatus, responseMessage
}
