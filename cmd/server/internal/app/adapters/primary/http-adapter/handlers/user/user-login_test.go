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

type mockService struct {
	loginFn func(ctx context.Context, username, password string) (*token.Tokens, error)

	loginCalls int
	lastU      string
	lastP      string
}

func (m *mockService) RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error) {
	panic("not used")
}
func (m *mockService) Login(ctx context.Context, username, password string) (*token.Tokens, error) {
	m.loginCalls++
	m.lastU = username
	m.lastP = password
	return m.loginFn(ctx, username, password)
}
func (m *mockService) Authenticate(tk string) (int64, error) { panic("not used") }
func (m *mockService) RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error) {
	panic("not used")
}

func TestHttpHandler_LoginHandler(t *testing.T) {
	// чтобы не падало на logger.Log.Error(...)
	logger.Log = zap.NewNop()

	t.Run("bad json -> 422 and service NOT called", func(t *testing.T) {
		ms := &mockService{
			loginFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return &token.Tokens{}, nil
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.LoginHandler(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.loginCalls != 0 {
			t.Fatalf("service must NOT be called, calls=%d", ms.loginCalls)
		}
	})

	t.Run("service returns error -> status != 200 and error response", func(t *testing.T) {
		ms := &mockService{
			loginFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return nil, errors.New("invalid credentials")
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(LoginRequest{Username: "john", Password: "bad"})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.LoginHandler(rr, req)

		if ms.loginCalls != 1 {
			t.Fatalf("expected service called once, got %d", ms.loginCalls)
		}
		if ms.lastU != "john" || ms.lastP != "bad" {
			t.Fatalf("unexpected args: u=%q p=%q", ms.lastU, ms.lastP)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected status != 200, got 200, body=%s", rr.Body.String())
		}
		if rr.Body.Len() == 0 {
			t.Fatalf("expected non-empty error body")
		}
	})

	t.Run("ok -> 200 and returns tokens", func(t *testing.T) {
		ms := &mockService{
			loginFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return &token.Tokens{
					JWTToken:     "jwt-123",
					RefreshToken: "rt-456",
				}, nil
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(LoginRequest{Username: "john", Password: "pass"})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.LoginHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var got LoginResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v, body=%s", err, rr.Body.String())
		}

		if got.Token != "jwt-123" {
			t.Fatalf("expected token=jwt-123, got %q", got.Token)
		}
		if got.RefreshToken != "rt-456" {
			t.Fatalf("expected refresh_token=rt-456, got %q", got.RefreshToken)
		}
	})
}
