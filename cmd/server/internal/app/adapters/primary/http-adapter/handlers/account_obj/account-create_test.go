package account_obj

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"server/internal/app/adapters/primary/http-adapter/constants"
	domain "server/internal/app/domain/account_obj"
	"server/internal/pkg/logger"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	logger.Log = zap.NewNop()

	os.Exit(m.Run())
}

type serviceMock struct {
	getAccountsListFn func(ctx context.Context, userId int64) ([]*domain.Account, error)
	getAccountFn      func(ctx context.Context, accountId int64) (*domain.Account, error)
	createFn          func(ctx context.Context, account *domain.Account) (int64, error)
	updateFn          func(ctx context.Context, account *domain.Account) error

	getAccountsListCalled int
	getAccountCalled      int
	createCalled          int
	updateCalled          int

	lastUserID    int64
	lastAccountID int64
	lastAccount   *domain.Account
}

func (m *serviceMock) GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error) {
	m.getAccountsListCalled++
	m.lastUserID = userId
	return m.getAccountsListFn(ctx, userId)
}

func (m *serviceMock) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	m.getAccountCalled++
	m.lastAccountID = accountId
	return m.getAccountFn(ctx, accountId)
}

func (m *serviceMock) CreateNewAccountObj(ctx context.Context, account *domain.Account) (int64, error) {
	m.createCalled++
	m.lastAccount = account
	return m.createFn(ctx, account)
}

func (m *serviceMock) UpdateAccount(ctx context.Context, account *domain.Account) error {
	m.updateCalled++
	m.lastAccount = account
	return m.updateFn(ctx, account)
}

func TestHttpHandler_CreateAccount(t *testing.T) {
	t.Run("bad json -> 422 and service not called", func(t *testing.T) {
		svc := &serviceMock{
			createFn: func(ctx context.Context, a *domain.Account) (int64, error) {
				t.Fatalf("service must not be called")
				return 0, nil
			},
			getAccountsListFn: func(context.Context, int64) ([]*domain.Account, error) { return nil, nil },
			getAccountFn:      func(context.Context, int64) (*domain.Account, error) { return nil, nil },
			updateFn:          func(context.Context, *domain.Account) error { return nil },
		}

		h := New(svc)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.CreateAccount(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d", rr.Code)
		}
		if svc.createCalled != 0 {
			t.Fatalf("expected createCalled=0, got %d", svc.createCalled)
		}
		if !strings.Contains(rr.Body.String(), "json decode error") {
			t.Fatalf("expected json decode error, body=%s", rr.Body.String())
		}
	})

	t.Run("missing userID in context -> 422 and service not called (this test will FAIL until you add return)", func(t *testing.T) {
		svc := &serviceMock{
			createFn: func(ctx context.Context, a *domain.Account) (int64, error) {
				t.Fatalf("BUG: service called even though userID missing in context. account=%+v", a)
				return 0, nil
			},
			getAccountsListFn: func(context.Context, int64) ([]*domain.Account, error) { return nil, nil },
			getAccountFn:      func(context.Context, int64) (*domain.Account, error) { return nil, nil },
			updateFn:          func(context.Context, *domain.Account) error { return nil },
		}

		h := New(svc)

		body, _ := json.Marshal(CreateAccountRequest{
			ServiceName: "gmail",
			UserName:    "john",
			Password:    "pass",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.CreateAccount(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d body=%s", rr.Code, rr.Body.String())
		}
		if svc.createCalled != 0 {
			t.Fatalf("expected createCalled=0, got %d", svc.createCalled)
		}
		if !strings.Contains(rr.Body.String(), "user ID not found in context") {
			t.Fatalf("expected userID error, body=%s", rr.Body.String())
		}
	})

	t.Run("service returns error -> mapped error response (status not 200)", func(t *testing.T) {
		svc := &serviceMock{
			createFn: func(ctx context.Context, a *domain.Account) (int64, error) {
				return 0, domain.ErrEmptyServiceName
			},
			getAccountsListFn: func(context.Context, int64) ([]*domain.Account, error) { return nil, nil },
			getAccountFn:      func(context.Context, int64) (*domain.Account, error) { return nil, nil },
			updateFn:          func(context.Context, *domain.Account) error { return nil },
		}

		h := New(svc)

		body, _ := json.Marshal(CreateAccountRequest{
			UserName: "john",
			Password: "pass",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(42)))

		rr := httptest.NewRecorder()
		h.CreateAccount(rr, req)

		if svc.createCalled != 1 {
			t.Fatalf("expected createCalled=1, got %d", svc.createCalled)
		}
		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200, got 200 body=%s", rr.Body.String())
		}
		if len(strings.TrimSpace(rr.Body.String())) == 0 {
			t.Fatalf("expected non-empty error body")
		}
	})

	t.Run("ok -> 200 and returns account_id, service called with proper domain", func(t *testing.T) {
		svc := &serviceMock{
			createFn: func(ctx context.Context, a *domain.Account) (int64, error) {
				if a.UserId != 101 {
					t.Fatalf("expected UserId=101, got %d", a.UserId)
				}
				if a.ServiceName != "gmail" || a.UserName != "john" || a.Password != "pass" {
					t.Fatalf("unexpected account: %+v", a)
				}
				return 777, nil
			},
			getAccountsListFn: func(context.Context, int64) ([]*domain.Account, error) { return nil, nil },
			getAccountFn:      func(context.Context, int64) (*domain.Account, error) { return nil, nil },
			updateFn:          func(context.Context, *domain.Account) error { return nil },
		}

		h := New(svc)

		body, _ := json.Marshal(CreateAccountRequest{
			ServiceName: "gmail",
			UserName:    "john",
			Password:    "pass",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(101)))

		rr := httptest.NewRecorder()
		h.CreateAccount(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
		}
		if svc.createCalled != 1 {
			t.Fatalf("expected createCalled=1, got %d", svc.createCalled)
		}

		var resp CreateAccountResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("response is not json: %v body=%s", err, rr.Body.String())
		}
		if resp.AccountID != 777 {
			t.Fatalf("expected account_id=777, got %d", resp.AccountID)
		}
	})
}
