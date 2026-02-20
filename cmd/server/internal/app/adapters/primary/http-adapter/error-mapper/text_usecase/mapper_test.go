package text_usecase

import (
	"errors"
	"net/http"
	"testing"

	domain "server/internal/app/domain/text_obj"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "ErrEmptyTextsList -> 204",
			err:        domain.ErrEmptyTextsList,
			wantStatus: http.StatusNoContent,
			wantMsg:    domain.ErrEmptyTextsList.Error(),
		},
		{
			name:       "ErrTextNotFound -> 404",
			err:        domain.ErrTextNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    domain.ErrTextNotFound.Error(),
		},
		{
			name:       "ErrInvalidUserID -> 400",
			err:        domain.ErrInvalidUserID,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrInvalidUserID.Error(),
		},
		{
			name:       "ErrInvalidTextID -> 400",
			err:        domain.ErrInvalidTextID,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrInvalidTextID.Error(),
		},
		{
			name:       "ErrEmptyTitle -> 400",
			err:        domain.ErrEmptyTitle,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrEmptyTitle.Error(),
		},
		{
			name:       "ErrEmptyText -> 400",
			err:        domain.ErrEmptyText,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrEmptyText.Error(),
		},
		{
			name:       "ErrFailedCreateText -> 500",
			err:        domain.ErrFailedCreateText,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    domain.ErrFailedCreateText.Error(),
		},
		{
			name:       "ErrFailedUpdateText -> 500",
			err:        domain.ErrFailedUpdateText,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    domain.ErrFailedUpdateText.Error(),
		},
		{
			name:       "unknown error -> 500 internal error",
			err:        errors.New("something bad happened"),
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
