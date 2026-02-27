package bank_card_obj

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	"testing"

	domain "server/internal/app/domain/bank_card_obj"
	"server/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

/* ---------- mock service ---------- */

type mockServiceS struct {
	getFn     func(ctx context.Context, cardId int64) (*domain.BankCard, error)
	listFn    func(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	createFn  func(ctx context.Context, card *domain.BankCard) (int64, error)
	updateFn  func(ctx context.Context, card *domain.BankCard) error
	calledGet bool
}

func (m *mockServiceS) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	m.calledGet = true
	return m.getFn(ctx, cardId)
}
func (m *mockServiceS) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	return nil, nil
}
func (m *mockServiceS) CreateNewBankCardObj(ctx context.Context, card *domain.BankCard) (int64, error) {
	return 0, nil
}
func (m *mockServiceS) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {
	return nil
}

/* ---------- tests ---------- */

func TestHttpHandler_GetBankCardObj(t *testing.T) {
	old := logger.Log
	logger.Log = zap.NewNop()
	t.Cleanup(func() { logger.Log = old })

	t.Run("invalid id -> 400 (service NOT called)", func(t *testing.T) {
		ms := &mockServiceS{
			getFn: func(ctx context.Context, cardId int64) (*domain.BankCard, error) {
				t.Fatalf("service must NOT be called")
				return nil, nil
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list/abc", nil)
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetBankCardObj(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.calledGet {
			t.Fatalf("service must NOT be called on invalid id")
		}
	})

	t.Run("service returns error -> mapped response", func(t *testing.T) {
		svcErr := domain.ErrBankCardNotFound
		wantStatus, wantMsg := errorMapper.Process(svcErr)

		ms := &mockServiceS{
			getFn: func(ctx context.Context, cardId int64) (*domain.BankCard, error) {
				if cardId != 10 {
					t.Fatalf("expected cardId=10, got %d", cardId)
				}
				return nil, svcErr
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list/10", nil)
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "10")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetBankCardObj(rr, req)

		if rr.Code != wantStatus {
			t.Fatalf("expected %d, got %d, body=%s", wantStatus, rr.Code, rr.Body.String())
		}
		if wantMsg != "" && !bytes.Contains(rr.Body.Bytes(), []byte(wantMsg)) {
			t.Fatalf("expected body to contain %q, got %s", wantMsg, rr.Body.String())
		}
		if !ms.calledGet {
			t.Fatalf("expected service to be called")
		}
	})

	t.Run("ok -> 200 + returns card", func(t *testing.T) {
		ms := &mockServiceS{
			getFn: func(ctx context.Context, cardId int64) (*domain.BankCard, error) {
				return &domain.BankCard{
					Bank: "MAIB",
					Pid:  "A123",
				}, nil
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list/7", nil)
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "7")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetBankCardObj(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var resp BankCardResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if resp.BankName != "MAIB" || resp.PID != "A123" {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})
}
