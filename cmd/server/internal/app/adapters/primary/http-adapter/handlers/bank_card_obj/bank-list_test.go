package bank_card_obj

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	"testing"

	domain "server/internal/app/domain/bank_card_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type mockServiceSS struct {
	getFn    func(ctx context.Context, cardId int64) (*domain.BankCard, error)
	listFn   func(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	createFn func(ctx context.Context, card *domain.BankCard) (int64, error)
	updateFn func(ctx context.Context, card *domain.BankCard) error

	calledList bool
}

func (m *mockServiceSS) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	return nil, nil
}

func (m *mockServiceSS) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	m.calledList = true
	return m.listFn(ctx, userId)
}

func (m *mockServiceSS) CreateNewBankCardObj(ctx context.Context, card *domain.BankCard) (int64, error) {
	return 0, nil
}

func (m *mockServiceSS) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {
	return nil
}

func TestHttpHandler_GetBankCardList(t *testing.T) {

	old := logger.Log
	logger.Log = zap.NewNop()
	t.Cleanup(func() { logger.Log = old })

	t.Run("missing userID in context -> 422 (service NOT called)", func(t *testing.T) {
		ms := &mockServiceSS{
			listFn: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				t.Fatalf("service must NOT be called when userID missing")
				return nil, nil
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		rr := httptest.NewRecorder()

		h.GetBankCardList(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.calledList {
			t.Fatalf("service help NOT be called when userID missing")
		}
	})

	t.Run("service returns error -> mapped error response", func(t *testing.T) {
		svcErr := domain.ErrBankCardNotFound
		wantStatus, wantMsg := errorMapper.Process(svcErr)

		ms := &mockServiceSS{
			listFn: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				if userId != 7 {
					t.Fatalf("expected userId=7, got %d", userId)
				}
				return nil, svcErr
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))
		rr := httptest.NewRecorder()

		h.GetBankCardList(rr, req)

		if rr.Code != wantStatus {
			t.Fatalf("expected %d, got %d, body=%s", wantStatus, rr.Code, rr.Body.String())
		}
		if wantMsg != "" && !bytes.Contains(rr.Body.Bytes(), []byte(wantMsg)) {
			t.Fatalf("expected body to contain %q, got %s", wantMsg, rr.Body.String())
		}
		if !ms.calledList {
			t.Fatalf("expected service to be called")
		}
	})

	t.Run("ok -> 200 + list mapped", func(t *testing.T) {
		ms := &mockServiceSS{
			listFn: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				return []*domain.BankCard{
					{CardId: 1, Bank: "MAIB"},
					{CardId: 2, Bank: "Moldindconbank"},
				}, nil
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))
		rr := httptest.NewRecorder()

		h.GetBankCardList(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var got []Card
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("expected 2 items, got %d: %+v", len(got), got)
		}
		if got[0].CardID != 1 || got[0].BankName != "MAIB" {
			t.Fatalf("unexpected item[0]: %+v", got[0])
		}
		if got[1].CardID != 2 || got[1].BankName != "Moldindconbank" {
			t.Fatalf("unexpected item[1]: %+v", got[1])
		}
	})
}
