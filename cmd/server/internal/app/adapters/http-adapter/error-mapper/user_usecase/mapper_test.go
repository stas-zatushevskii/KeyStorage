package user_usecase

import (
	"errors"
	"net/http"
	"testing"

	domain "server/internal/app/domain/user"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "ErrUsernameAlreadyExists -> 409",
			err:        domain.ErrUsernameAlreadyExists,
			wantStatus: http.StatusConflict,
			wantMsg:    domain.ErrUsernameAlreadyExists.Error(),
		},
		{
			name:       "ErrPasswordMismatch -> 401",
			err:        domain.ErrPasswordMismatch,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    domain.ErrPasswordMismatch.Error(),
		},
		{
			name:       "ErrTokenNotValid -> 401",
			err:        domain.ErrTokenNotValid,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    domain.ErrTokenNotValid.Error(),
		},
		{
			name:       "ErrInvalidRefreshToken -> 401",
			err:        domain.ErrInvalidRefreshToken,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    domain.ErrInvalidRefreshToken.Error(),
		},
		{
			name:       "ErrRefreshTokenNotFound -> 400",
			err:        domain.ErrRefreshTokenNotFound,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrRefreshTokenNotFound.Error(),
		},
		{
			name:       "ErrTokenRevoked -> 403",
			err:        domain.ErrTokenRevoked,
			wantStatus: http.StatusForbidden,
			wantMsg:    domain.ErrTokenRevoked.Error(),
		},
		{
			name:       "ErrUserNotFound -> 403",
			err:        domain.ErrUserNotFound,
			wantStatus: http.StatusForbidden,
			wantMsg:    domain.ErrUserNotFound.Error(),
		},
		{
			name:       "unknown error -> 500",
			err:        errors.New("some internal error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "some internal error",
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
