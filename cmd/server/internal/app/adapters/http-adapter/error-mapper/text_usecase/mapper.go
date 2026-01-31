package text_usecase

import (
	"errors"
	"net/http"
	domain "server/internal/app/domain/text_obj"
)

// Process return httpStatus (200, 400 ...) and ErrMsg according to custom error type.
func Process(err error) (int, string) {
	switch {

	case errors.Is(err, domain.ErrEmptyTextsList):
		return http.StatusNoContent, err.Error()

	case errors.Is(err, domain.ErrTextNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, domain.ErrInvalidUserID),
		errors.Is(err, domain.ErrInvalidTextID),
		errors.Is(err, domain.ErrEmptyTitle),
		errors.Is(err, domain.ErrEmptyText):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, domain.ErrFailedCreateText),
		errors.Is(err, domain.ErrFailedUpdateText):
		return http.StatusInternalServerError, err.Error()

	default:
		return http.StatusInternalServerError, "internal error"
	}
}
