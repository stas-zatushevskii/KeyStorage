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

type mockServiceRegister struct {
	registerFn func(ctx context.Context, username, password string) (*token.Tokens, error)

	registerCalls int
	lastUsername  string
	lastPassword  string
}

func (m *mockServiceRegister) RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error) {
	m.registerCalls++
	m.lastUsername = username
	m.lastPassword = password
	return m.registerFn(ctx, username, password)
}

func (m *mockServiceRegister) Login(ctx context.Context, username, password string) (*token.Tokens, error) {
	panic("not used")
}
func (m *mockServiceRegister) Authenticate(tk string) (int64, error) {
	panic("not used")
}
func (m *mockServiceRegister) RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error) {
	panic("not used")
}

func TestHttpHandler_RegistrationHandler(t *testing.T) {
	// иначе упадет на logger.Log.Error(...)
	logger.Log = zap.NewNop()

	t.Run("bad json -> 422 and service NOT called", func(t *testing.T) {
		ms := &mockServiceRegister{
			registerFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return &token.Tokens{}, nil
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.RegistrationHandler(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.registerCalls != 0 {
			t.Fatalf("service must NOT be called, calls=%d", ms.registerCalls)
		}
	})

	t.Run("service returns error -> status != 200 and error response", func(t *testing.T) {
		ms := &mockServiceRegister{
			registerFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return nil, errors.New("register failed")
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(RegisterNewUserRequest{
			Username: "john",
			Password: "secret",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.RegistrationHandler(rr, req)

		if ms.registerCalls != 1 {
			t.Fatalf("expected service called once, got %d", ms.registerCalls)
		}
		if ms.lastUsername != "john" || ms.lastPassword != "secret" {
			t.Fatalf("unexpected args: username=%q password=%q", ms.lastUsername, ms.lastPassword)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected status != 200, got 200, body=%s", rr.Body.String())
		}
		if rr.Body.Len() == 0 {
			t.Fatalf("expected non-empty error body")
		}
	})

	t.Run("ok -> 200 and returns tokens", func(t *testing.T) {
		ms := &mockServiceRegister{
			registerFn: func(ctx context.Context, username, password string) (*token.Tokens, error) {
				return &token.Tokens{
					JWTToken:     "jwt-1",
					RefreshToken: "rt-1",
				}, nil
			},
		}
		h := &HttpHandler{service: ms}

		body, _ := json.Marshal(RegisterNewUserRequest{
			Username: "john",
			Password: "secret",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.RegistrationHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var got RegisterNewUserResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v, body=%s", err, rr.Body.String())
		}

		if got.Token != "jwt-1" {
			t.Fatalf("expected token=jwt-1, got %q", got.Token)
		}
		if got.RefreshToken != "rt-1" {
			t.Fatalf("expected refresh_token=rt-1, got %q", got.RefreshToken)
		}
	})
}
