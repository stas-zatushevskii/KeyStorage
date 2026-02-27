package user_obj

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/pkg/logger"
	"server/internal/pkg/token"

	"go.uber.org/zap"
)

type mockServiceRefresh struct {
	refreshFn func(ctx context.Context, jwt, refresh string) (*token.Tokens, error)

	refreshCalls int
	lastJWT      string
	lastRT       string
}

func (m *mockServiceRefresh) RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error) {
	panic("not used")
}
func (m *mockServiceRefresh) Login(ctx context.Context, username, password string) (*token.Tokens, error) {
	panic("not used")
}
func (m *mockServiceRefresh) Authenticate(tk string) (int64, error) {
	panic("not used")
}
func (m *mockServiceRefresh) RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error) {
	m.refreshCalls++
	m.lastJWT = jwt
	m.lastRT = refreshToken
	return m.refreshFn(ctx, jwt, refreshToken)
}

func TestHttpHandler_RefreshTokenHandler(t *testing.T) {
	// чтобы не падало на logger.Log.Error(...)
	logger.Log = zap.NewNop()

	t.Run("bad json -> 422 and service NOT called", func(t *testing.T) {
		ms := &mockServiceRefresh{
			refreshFn: func(ctx context.Context, jwt, refresh string) (*token.Tokens, error) {
				return &token.Tokens{}, nil
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.RefreshTokenHandler(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.refreshCalls != 0 {
			t.Fatalf("service must NOT be called, calls=%d", ms.refreshCalls)
		}
	})

	t.Run("service returns error -> status != 200 and error response", func(t *testing.T) {
		ms := &mockServiceRefresh{
			refreshFn: func(ctx context.Context, jwt, refresh string) (*token.Tokens, error) {
				return nil, errors.New("refresh failed")
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(RefreshTokenRequest{
			Username: "jwt-123",
			Password: "rt-456",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.RefreshTokenHandler(rr, req)

		if ms.refreshCalls != 1 {
			t.Fatalf("expected service called once, got %d", ms.refreshCalls)
		}
		if ms.lastJWT != "jwt-123" || ms.lastRT != "rt-456" {
			t.Fatalf("unexpected args: jwt=%q rt=%q", ms.lastJWT, ms.lastRT)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected status != 200, got 200, body=%s", rr.Body.String())
		}
		if rr.Body.Len() == 0 {
			t.Fatalf("expected non-empty error body")
		}
	})

	t.Run("ok -> 200 and returns tokens", func(t *testing.T) {
		ms := &mockServiceRefresh{
			refreshFn: func(ctx context.Context, jwt, refresh string) (*token.Tokens, error) {
				return &token.Tokens{
					JWTToken:     "new-jwt",
					RefreshToken: "new-rt",
				}, nil
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(RefreshTokenRequest{
			Username: "old-jwt",
			Password: "old-rt",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.RefreshTokenHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var got RefreshTokenResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v, body=%s", err, rr.Body.String())
		}

		if got.Token != "new-jwt" {
			t.Fatalf("expected token=new-jwt, got %q", got.Token)
		}
		if got.RefreshToken != "new-rt" {
			t.Fatalf("expected refresh_token=new-rt, got %q", got.RefreshToken)
		}
	})
}
