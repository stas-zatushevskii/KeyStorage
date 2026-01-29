package account_usecase

import (
	"errors"
	"net/http"
	"testing"

	domain "server/internal/app/domain/account_obj"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantMessage    string
	}{
		{
			name:           "ErrEmptyAccountsList -> 204",
			err:            domain.ErrEmptyAccountsList,
			wantStatusCode: http.StatusNoContent,
			wantMessage:    domain.ErrEmptyAccountsList.Error(),
		},
		{
			name:           "ErrAccountNotFound -> 404",
			err:            domain.ErrAccountNotFound,
			wantStatusCode: http.StatusNotFound,
			wantMessage:    domain.ErrAccountNotFound.Error(),
		},
		{
			name:           "ErrInvalidUserID -> 400",
			err:            domain.ErrInvalidUserID,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    domain.ErrInvalidUserID.Error(),
		},
		{
			name:           "ErrInvalidAccountID -> 400",
			err:            domain.ErrInvalidAccountID,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    domain.ErrInvalidAccountID.Error(),
		},
		{
			name:           "ErrEmptyServiceName -> 400",
			err:            domain.ErrEmptyServiceName,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    domain.ErrEmptyServiceName.Error(),
		},
		{
			name:           "ErrFailedCreateAccount -> 500",
			err:            domain.ErrFailedCreateAccount,
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    domain.ErrFailedCreateAccount.Error(),
		},
		{
			name:           "ErrFailedUpdateAccount -> 500",
			err:            domain.ErrFailedUpdateAccount,
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    domain.ErrFailedUpdateAccount.Error(),
		},
		{
			name:           "unknown error -> 500 internal error",
			err:            errors.New("some random error"),
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := Process(tt.err)

			if status != tt.wantStatusCode {
				t.Fatalf("expected status %d, got %d", tt.wantStatusCode, status)
			}

			if msg != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, msg)
			}
		})
	}
}
