package bank_card_usecase

import (
	"errors"
	"net/http"
	domain "server/internal/app/domain/bank_card_obj"
)

// Process return httpStatus (200, 400 ...) and ErrMsg according to custom error type.
func Process(err error) (int, string) {
	switch {

	case errors.Is(err, domain.ErrInvalidUserID),
		errors.Is(err, domain.ErrInvalidCardID),
		errors.Is(err, domain.ErrEmptyBankName),
		errors.Is(err, domain.ErrEmptyPID):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, domain.ErrFaildeCreateBankCardObject),
		errors.Is(err, domain.ErrFailedUpdateBankCard):
		return http.StatusInternalServerError, err.Error()

	default:
		return http.StatusInternalServerError, "internal error"
	}
}
