package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	"strings"
	"testing"
)

// ---- mock auth service ----

type mockAuthService struct {
	authenticateFn func(token string) (int64, error)

	calls     int
	lastToken string
}

func (m *mockAuthService) Authenticate(token string) (int64, error) {
	m.calls++
	m.lastToken = token
	return m.authenticateFn(token)
}

func TestJWTMiddleware(t *testing.T) {
	t.Run("missing Authorization -> 401 and service NOT called", func(t *testing.T) {
		svc := &mockAuthService{
			authenticateFn: func(token string) (int64, error) {
				return 0, nil
			},
		}

		nextCalled := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled++
			w.WriteHeader(http.StatusOK)
		})

		h := JWTMiddleware(svc)(next)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil) // without Authorization
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "JWT token not found") {
			t.Fatalf("expected body to contain %q, got %q", "JWT token not found", rr.Body.String())
		}
		if svc.calls != 0 {
			t.Fatalf("Authenticate must NOT be called, calls=%d", svc.calls)
		}
		if nextCalled != 0 {
			t.Fatalf("next must NOT be called, nextCalled=%d", nextCalled)
		}
	})

	t.Run("Authenticate returns error -> status != 200 and next NOT called", func(t *testing.T) {
		svc := &mockAuthService{
			authenticateFn: func(token string) (int64, error) {
				return 0, errors.New("invalid token")
			},
		}

		nextCalled := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled++
			w.WriteHeader(http.StatusOK)
		})

		h := JWTMiddleware(svc)(next)

		req := httptest.NewRequest(http.MethodGet, "/any", nil)
		req.Header.Set("Authorization", "Bearer xxx")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if svc.calls != 1 {
			t.Fatalf("expected Authenticate called once, got %d", svc.calls)
		}
		if svc.lastToken != "Bearer xxx" {
			t.Fatalf("unexpected token passed to Authenticate: %q", svc.lastToken)
		}

		// точный статус зависит от errorMapper.Process(err),
		// поэтому проверяем что это НЕ 200 и что тело не пустое
		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200, got 200, body=%s", rr.Body.String())
		}
		if rr.Body.Len() == 0 {
			t.Fatalf("expected non-empty body")
		}
		if nextCalled != 0 {
			t.Fatalf("next must NOT be called, nextCalled=%d", nextCalled)
		}
	})

	t.Run("ok -> calls next and puts userID in context", func(t *testing.T) {
		svc := &mockAuthService{
			authenticateFn: func(token string) (int64, error) {
				if token != "Bearer good" {
					return 0, errors.New("unexpected token")
				}
				return int64(77), nil
			},
		}

		nextCalled := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled++

			v := r.Context().Value(constants.UserIDKey)
			id, ok := v.(int64)
			if !ok {
				t.Fatalf("expected int64 userID in context, got %T (%v)", v, v)
			}
			if id != 77 {
				t.Fatalf("expected userID=77, got %d", id)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		h := JWTMiddleware(svc)(next)

		req := httptest.NewRequest(http.MethodGet, "/any", nil)
		req.Header.Set("Authorization", "Bearer good")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if svc.calls != 1 {
			t.Fatalf("expected Authenticate called once, got %d", svc.calls)
		}
		if nextCalled != 1 {
			t.Fatalf("expected next called once, got %d", nextCalled)
		}
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if rr.Body.String() != "ok" {
			t.Fatalf("expected body %q, got %q", "ok", rr.Body.String())
		}
	})
}
