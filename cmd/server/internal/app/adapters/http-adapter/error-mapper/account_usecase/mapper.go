package account_usecase

import (
	"errors"
	"net/http"
	domain "server/internal/app/domain/account_obj"
)

// Process return httpStatus (200, 400 ...) and ErrMsg according to custom error type.
func Process(err error) (int, string) {

	switch {
	case errors.Is(err, domain.ErrEmptyAccountsList):
		return http.StatusNoContent, err.Error()

	case errors.Is(err, domain.ErrAccountNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, domain.ErrInvalidUserID),
		errors.Is(err, domain.ErrInvalidAccountID),
		errors.Is(err, domain.ErrEmptyServiceName):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, domain.ErrFailedCreateAccount),
		errors.Is(err, domain.ErrFailedUpdateAccount):
		return http.StatusInternalServerError, err.Error()

	default:
		return http.StatusInternalServerError, "internal error"
	}
}
