package account_obj_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	handler "server/internal/app/adapters/primary/http-adapter/handlers/account_obj"
	"testing"

	domain "server/internal/app/domain/account_obj"
	"server/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mockAccountServiceS struct {
	updateAccountFn func(ctx context.Context, account *domain.Account) error
}

func (m *mockAccountServiceS) UpdateAccount(ctx context.Context, account *domain.Account) error {
	if m.updateAccountFn == nil {
		return errors.New("updateAccountFn is nil")
	}
	return m.updateAccountFn(ctx, account)
}

// Остальные методы интерфейса handler.service (не используются тут)
func (m *mockAccountServiceS) GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAccountServiceS) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAccountServiceS) CreateNewAccountObj(ctx context.Context, account *domain.Account) (int64, error) {
	return 0, errors.New("not implemented")
}

func withChiURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestHttpHandler_UpdateAccountObj(t *testing.T) {
	// иначе в ветке ошибки будет panic на logger.Log.Error(...)
	logger.Log = zap.NewNop()

	t.Run("invalid json -> 422", func(t *testing.T) {
		svc := &mockAccountServiceS{
			updateAccountFn: func(ctx context.Context, account *domain.Account) error {
				t.Fatalf("service must NOT be called when json invalid")
				return nil
			},
		}
		h := handler.New(svc)

		req := httptest.NewRequest(http.MethodPut, "/update/10", bytes.NewBufferString("{bad json"))
		req = withChiURLParam(req, "id", "10")
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateAccountObj(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
		}
	})

	t.Run("invalid url id -> 400", func(t *testing.T) {
		svc := &mockAccountServiceS{
			updateAccountFn: func(ctx context.Context, account *domain.Account) error {
				t.Fatalf("service must NOT be called when url id invalid")
				return nil
			},
		}
		h := handler.New(svc)

		body, _ := json.Marshal(map[string]any{
			"service_name": "github",
			"user_name":    "stas",
			"password":     "p",
		})

		req := httptest.NewRequest(http.MethodPut, "/update/nope", bytes.NewReader(body))
		req = withChiURLParam(req, "id", "nope")
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateAccountObj(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("missing userID in context -> 422 (and service NOT called)", func(t *testing.T) {
		called := false
		svc := &mockAccountServiceS{
			updateAccountFn: func(ctx context.Context, account *domain.Account) error {
				called = true
				return nil
			},
		}
		h := handler.New(svc)

		body, _ := json.Marshal(map[string]any{
			"service_name": "github",
			"user_name":    "stas",
			"password":     "p",
		})

		req := httptest.NewRequest(http.MethodPut, "/update/10", bytes.NewReader(body))
		req = withChiURLParam(req, "id", "10")
		// userID НЕ кладем

		rr := httptest.NewRecorder()
		h.UpdateAccountObj(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
		}
		if called {
			t.Fatalf("service must NOT be called when userID missing")
		}
	})

	t.Run("service returns error -> non-200", func(t *testing.T) {
		serviceErr := errors.New("boom")

		var got *domain.Account
		svc := &mockAccountServiceS{
			updateAccountFn: func(ctx context.Context, account *domain.Account) error {
				got = account
				return serviceErr
			},
		}
		h := handler.New(svc)

		body, _ := json.Marshal(map[string]any{
			"service_name": "github",
			"user_name":    "stas",
			"password":     "p",
		})

		req := httptest.NewRequest(http.MethodPut, "/update/55", bytes.NewReader(body))
		req = withChiURLParam(req, "id", "55")
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateAccountObj(rr, req)

		if got == nil {
			t.Fatalf("expected service to be called with account")
		}
		if got.AccountId != 55 {
			t.Fatalf("expected AccountId=55, got %d", got.AccountId)
		}
		if got.UserId != 7 {
			t.Fatalf("expected UserId=7, got %d", got.UserId)
		}
		if got.ServiceName != "github" || got.UserName != "stas" || got.Password != "p" {
			t.Fatalf("unexpected account payload: %+v", got)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200, got 200 body=%s", rr.Body.String())
		}
	})

	t.Run("ok -> 200", func(t *testing.T) {
		var got *domain.Account

		svc := &mockAccountServiceS{
			updateAccountFn: func(ctx context.Context, account *domain.Account) error {
				got = account
				return nil
			},
		}
		h := handler.New(svc)

		body, _ := json.Marshal(map[string]any{
			"service_name": "gmail",
			"user_name":    "stas.z",
			"password":     "secret",
		})

		req := httptest.NewRequest(http.MethodPut, "/update/99", bytes.NewReader(body))
		req = withChiURLParam(req, "id", "99")
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(100)))

		rr := httptest.NewRecorder()
		h.UpdateAccountObj(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
		}

		if got == nil {
			t.Fatalf("expected account passed to service, got nil")
		}
		if got.AccountId != 99 || got.UserId != 100 {
			t.Fatalf("unexpected ids: %+v", got)
		}
	})
}
