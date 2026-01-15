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
	default:
		httpStatus = http.StatusInternalServerError
		responseMessage = err.Error()
	}
	return httpStatus, responseMessage
}
