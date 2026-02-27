package account_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	"server/internal/app/adapters/primary/http-adapter/handlers/account"
	"testing"

	domain "server/internal/app/domain/account"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type mockAccountService struct {
	getAccountsListFn func(ctx context.Context, userId int64) ([]*domain.Account, error)
}

func (m *mockAccountService) GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error) {
	if m.getAccountsListFn == nil {
		return nil, errors.New("getAccountsListFn is nil")
	}
	return m.getAccountsListFn(ctx, userId)
}

// Остальные методы интерфейса handler.service — не используются в этих тестах, но нужны чтобы мок компилился.
func (m *mockAccountService) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAccountService) CreateNewAccount(ctx context.Context, account *domain.Account) (int64, error) {
	return 0, errors.New("not implemented")
}
func (m *mockAccountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return errors.New("not implemented")
}

func TestHttpHandler_GetAccountList(t *testing.T) {
	// иначе будет panic на logger.Log.Error(...)
	logger.Log = zap.NewNop()

	t.Run("missing userID in context -> 422", func(t *testing.T) {
		svc := &mockAccountService{
			getAccountsListFn: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				t.Fatalf("service must NOT be called when userID missing")
				return nil, nil
			},
		}
		h := account.New(svc)

		req := httptest.NewRequest(http.MethodGet, "/list", nil) // без userId в контексте
		rr := httptest.NewRecorder()

		h.GetAccountList(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
		}
	})

	t.Run("service returns error -> non-200", func(t *testing.T) {
		serviceErr := errors.New("boom")

		called := false
		var gotUserID int64

		svc := &mockAccountService{
			getAccountsListFn: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				called = true
				gotUserID = userId
				return nil, serviceErr
			},
		}
		h := account.New(svc)

		req := httptest.NewRequest(http.MethodGet, "/list", nil).
			WithContext(context.WithValue(context.Background(), constants.UserIDKey, int64(77)))

		rr := httptest.NewRecorder()
		h.GetAccountList(rr, req)

		if !called {
			t.Fatalf("expected service.GetAccountsList to be called")
		}
		if gotUserID != 77 {
			t.Fatalf("expected userId=77, got %d", gotUserID)
		}

		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200, got 200 body=%s", rr.Body.String())
		}
	})

	t.Run("ok -> empty list -> 200 and []", func(t *testing.T) {
		svc := &mockAccountService{
			getAccountsListFn: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				return []*domain.Account{}, nil
			},
		}
		h := account.New(svc)

		req := httptest.NewRequest(http.MethodGet, "/list", nil).
			WithContext(context.WithValue(context.Background(), constants.UserIDKey, int64(1)))

		rr := httptest.NewRecorder()
		h.GetAccountList(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp []account.Account
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json: %v body=%s", err, rr.Body.String())
		}
		if len(resp) != 0 {
			t.Fatalf("expected empty array, got: %+v", resp)
		}
	})

	t.Run("ok -> maps items -> 200", func(t *testing.T) {
		svc := &mockAccountService{
			getAccountsListFn: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				if userId != 5 {
					t.Fatalf("expected userId=5, got %d", userId)
				}
				return []*domain.Account{
					{AccountId: 10, ServiceName: "github", UserName: "stas"},
					{AccountId: 11, ServiceName: "gmail", UserName: "stas.z"},
				}, nil
			},
		}
		h := account.New(svc)

		req := httptest.NewRequest(http.MethodGet, "/list", nil).
			WithContext(context.WithValue(context.Background(), constants.UserIDKey, int64(5)))

		rr := httptest.NewRecorder()
		h.GetAccountList(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp []account.Account
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json: %v body=%s", err, rr.Body.String())
		}

		if len(resp) != 2 {
			t.Fatalf("expected 2 items, got %d: %+v", len(resp), resp)
		}

		if resp[0].AccountID != 10 || resp[0].ServiceName != "github" || resp[0].Username != "stas" {
			t.Fatalf("unexpected item0: %+v", resp[0])
		}
		if resp[1].AccountID != 11 || resp[1].ServiceName != "gmail" || resp[1].Username != "stas.z" {
			t.Fatalf("unexpected item1: %+v", resp[1])
		}
	})
}
