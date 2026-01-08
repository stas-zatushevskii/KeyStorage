package account_usecase

import (
	"errors"
	"net/http"
	domain "server/internal/app/domain/account_obj"
)

// Process return httpStatus (200, 400 ...) and ErrMsg according to custom error type.
func Process(err error) (int, string) {
	var (
		httpStatus      int
		responseMessage string
	)

	switch {
	case errors.Is(err, domain.ErrEmptyAccountsList):
		httpStatus = http.StatusNoContent
		responseMessage = err.Error()
	case errors.Is(err, domain.ErrAccountInformationNotFound):
		httpStatus = http.StatusNotFound
		responseMessage = err.Error()
	default:
		httpStatus = http.StatusInternalServerError
		responseMessage = err.Error()
	}
	return httpStatus, responseMessage
}
