package bank_card

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	"testing"

	domain "server/internal/app/domain/bank_card"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

type mockService struct {
	getFn        func(ctx context.Context, cardId int64) (*domain.BankCard, error)
	listFn       func(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	createFn     func(ctx context.Context, card *domain.BankCard) (int64, error)
	updateFn     func(ctx context.Context, card *domain.BankCard) error
	calledGet    bool
	calledList   bool
	calledCreate bool
	calledUpdate bool
}

func (m *mockService) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	m.calledGet = true
	if m.getFn == nil {
		return nil, errors.New("getFn is nil")
	}
	return m.getFn(ctx, cardId)
}

func (m *mockService) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	m.calledList = true
	if m.listFn == nil {
		return nil, errors.New("listFn is nil")
	}
	return m.listFn(ctx, userId)
}

func (m *mockService) CreateNewBankCard(ctx context.Context, card *domain.BankCard) (int64, error) {
	m.calledCreate = true
	if m.createFn == nil {
		return 0, errors.New("createFn is nil")
	}
	return m.createFn(ctx, card)
}

func (m *mockService) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {
	m.calledUpdate = true
	if m.updateFn == nil {
		return errors.New("updateFn is nil")
	}
	return m.updateFn(ctx, card)
}

func TestHttpHandler_CreateBankCard(t *testing.T) {
	old := logger.Log
	logger.Log = zap.NewNop()
	t.Cleanup(func() { logger.Log = old })

	t.Run("bad json -> 422", func(t *testing.T) {
		ms := &mockService{
			createFn: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				return 0, nil
			},
		}
		h := New(ms)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.CreateBankCard(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
		}
		if ms.calledCreate {
			t.Fatalf("service must NOT be called on json decode error")
		}
	})

	t.Run("missing userID in context -> 422 (and service NOT called)", func(t *testing.T) {
		ms := &mockService{
			createFn: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				return 0, nil
			},
		}
		h := New(ms)

		body, _ := json.Marshal(map[string]any{
			"bank_name": "MAIB",
			"pid":       "A123",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.CreateBankCard(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
		}
		if ms.calledCreate {
			t.Fatalf("service must NOT be called when userID missing")
		}
	})

	t.Run("service returns error -> mapped status/message", func(t *testing.T) {
		svcErr := domain.ErrInvalidUserID
		wantStatus, wantMsg := errorMapper.Process(svcErr)

		ms := &mockService{
			createFn: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				return 0, svcErr
			},
		}
		h := New(ms)

		body, _ := json.Marshal(map[string]any{
			"bank_name": "MAIB",
			"pid":       "A123",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(777)))

		rr := httptest.NewRecorder()
		h.CreateBankCard(rr, req)

		if rr.Code != wantStatus {
			t.Fatalf("expected %d, got %d, body=%s", wantStatus, rr.Code, rr.Body.String())
		}
		if wantMsg != "" && !bytes.Contains(rr.Body.Bytes(), []byte(wantMsg)) {
			t.Fatalf("expected body to contain %q, got body=%s", wantMsg, rr.Body.String())
		}
		if !ms.calledCreate {
			t.Fatalf("expected service to be called")
		}
	})

	t.Run("ok -> 200 + returns card_id and passes userId into domain", func(t *testing.T) {
		var got *domain.BankCard

		ms := &mockService{
			createFn: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				got = card
				return 55, nil
			},
		}
		h := New(ms)

		body, _ := json.Marshal(map[string]any{
			"bank_name": "MAIB",
			"pid":       "A123",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(777)))

		rr := httptest.NewRecorder()
		h.CreateBankCard(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var resp CreateBankCardResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v, body=%s", err, rr.Body.String())
		}
		if resp.CardID != 55 {
			t.Fatalf("expected card_id=55, got %d", resp.CardID)
		}

		if !ms.calledCreate {
			t.Fatalf("expected service to be called")
		}
		if got == nil {
			t.Fatalf("expected service to receive domain card, got nil")
		}
		if got.UserId != 777 {
			t.Fatalf("expected domain.UserId=777, got %d", got.UserId)
		}
		if got.Bank != "MAIB" || got.Pid != "A123" {
			t.Fatalf("unexpected domain: %+v", got)
		}
	})
}
