package bank_card_usecase

import (
	"errors"
	"net/http"
	"testing"

	domain "server/internal/app/domain/bank_card_obj"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "ErrInvalidUserID -> 400",
			err:        domain.ErrInvalidUserID,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrInvalidUserID.Error(),
		},
		{
			name:       "ErrInvalidCardID -> 400",
			err:        domain.ErrInvalidCardID,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrInvalidCardID.Error(),
		},
		{
			name:       "ErrEmptyBankName -> 400",
			err:        domain.ErrEmptyBankName,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrEmptyBankName.Error(),
		},
		{
			name:       "ErrEmptyPID -> 400",
			err:        domain.ErrEmptyPID,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrEmptyPID.Error(),
		},
		{
			name:       "ErrFaildeCreateBankCardObject -> 500",
			err:        domain.ErrFaildeCreateBankCardObject,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    domain.ErrFaildeCreateBankCardObject.Error(),
		},
		{
			name:       "ErrFailedUpdateBankCard -> 500",
			err:        domain.ErrFailedUpdateBankCard,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    domain.ErrFailedUpdateBankCard.Error(),
		},
		{
			name:       "unknown error -> 500 internal error",
			err:        errors.New("boom"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := Process(tt.err)

			if status != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, status)
			}
			if msg != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, msg)
			}
		})
	}
}
