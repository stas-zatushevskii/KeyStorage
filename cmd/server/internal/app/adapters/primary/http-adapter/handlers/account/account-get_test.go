package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "server/internal/app/domain/account"
	"server/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mockService struct {
	getAccountFn func(ctx context.Context, accountId int64) (*domain.Account, error)
}

func (m *mockService) GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error) {
	panic("not used")
}
func (m *mockService) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	if m.getAccountFn == nil {
		panic("getAccountFn is nil")
	}
	return m.getAccountFn(ctx, accountId)
}
func (m *mockService) CreateNewAccount(ctx context.Context, account *domain.Account) (int64, error) {
	panic("not used")
}
func (m *mockService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	panic("not used")
}

func newChiReq(method, path, routePattern, paramKey, paramValue string) *http.Request {
	req := httptest.NewRequest(method, path, nil)

	rctx := chi.NewRouteContext()
	rctx.RoutePatterns = append(rctx.RoutePatterns, routePattern)
	rctx.URLParams.Add(paramKey, paramValue)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestHttpHandler_GetAccountObj(t *testing.T) {
	logger.Log = zap.NewNop()

	t.Run("invalid id -> 400", func(t *testing.T) {
		svc := &mockService{
			getAccountFn: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				t.Fatalf("service must NOT be called on invalid id")
				return nil, nil
			},
		}
		h := New(svc)

		req := newChiReq(http.MethodGet, "/list/abc", "/list/{id}", "id", "abc")
		rr := httptest.NewRecorder()

		h.GetAccountObj(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}

		var body map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &body)
	})

	t.Run("service returns error -> mapped status != 200", func(t *testing.T) {
		serviceErr := domain.ErrAccountNotFound

		called := false
		gotID := int64(0)

		svc := &mockService{
			getAccountFn: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				called = true
				gotID = accountId
				return nil, serviceErr
			},
		}
		h := New(svc)

		req := newChiReq(http.MethodGet, "/list/10", "/list/{id}", "id", "10")
		rr := httptest.NewRecorder()

		h.GetAccountObj(rr, req)

		if !called {
			t.Fatalf("expected service.GetAccount to be called")
		}
		if gotID != 10 {
			t.Fatalf("expected id=10, got %d", gotID)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200, got 200 body=%s", rr.Body.String())
		}
	})

	t.Run("ok -> 200 + json response", func(t *testing.T) {
		svc := &mockService{
			getAccountFn: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				if accountId != 7 {
					t.Fatalf("expected id=7, got %d", accountId)
				}
				return &domain.Account{
					AccountId:   7,
					ServiceName: "github",
					UserName:    "stas",
					Password:    "secret",
				}, nil
			},
		}
		h := New(svc)

		req := newChiReq(http.MethodGet, "/list/7", "/list/{id}", "id", "7")
		rr := httptest.NewRecorder()

		h.GetAccountObj(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp AccountResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json: %v body=%s", err, rr.Body.String())
		}

		if resp.ServiceName != "github" || resp.Username != "stas" || resp.Password != "secret" {
			t.Fatalf("unexpected resp: %+v", resp)
		}
	})

	t.Run("service returns unexpected error -> still mapped (non-200) OR internal", func(t *testing.T) {
		dbErr := errors.New("db down")

		svc := &mockService{
			getAccountFn: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				return nil, dbErr
			},
		}
		h := New(svc)

		req := newChiReq(http.MethodGet, "/list/5", "/list/{id}", "id", "5")
		rr := httptest.NewRecorder()

		h.GetAccountObj(rr, req)

		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200 for error, got 200 body=%s", rr.Body.String())
		}
	})
}
