package bank_card

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	errorMapper "server/internal/app/adapters/primary/http-adapter/error-mapper/bank_card_usecase"
	"testing"

	domain "server/internal/app/domain/bank_card"
	"server/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

/* ---------- mock service ---------- */

type mockServicE struct {
	updateFn func(ctx context.Context, card *domain.BankCard) error

	calledUpdate bool
	lastCard     *domain.BankCard
}

func (m *mockServicE) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	return nil, nil
}
func (m *mockServicE) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	return nil, nil
}
func (m *mockServicE) CreateNewBankCard(ctx context.Context, card *domain.BankCard) (int64, error) {
	return 0, nil
}
func (m *mockServicE) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {
	m.calledUpdate = true
	m.lastCard = card
	if m.updateFn == nil {
		return nil
	}
	return m.updateFn(ctx, card)
}

/* ---------- helpers ---------- */

func newReqWithChiID(t *testing.T, method, path, id string, body []byte) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

/* ---------- tests ---------- */

func TestHttpHandler_UpdateBankCardObj(t *testing.T) {
	// чтобы logger.Log.Error не паниковал в тестах
	old := logger.Log
	logger.Log = zap.NewNop()
	t.Cleanup(func() { logger.Log = old })

	t.Run("json decode error -> 422 (service NOT called)", func(t *testing.T) {
		ms := &mockServicE{
			updateFn: func(ctx context.Context, card *domain.BankCard) error {
				t.Fatalf("service must NOT be called on decode error")
				return nil
			},
		}
		h := New(ms)

		req := newReqWithChiID(t, http.MethodPut, "/update/10", "10", []byte("{bad json"))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateBankCardObj(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.calledUpdate {
			t.Fatalf("service must NOT be called on decode error")
		}
	})

	t.Run("invalid card id -> 400 (service NOT called)", func(t *testing.T) {
		ms := &mockServicE{
			updateFn: func(ctx context.Context, card *domain.BankCard) error {
				t.Fatalf("service must NOT be called on invalid id")
				return nil
			},
		}
		h := New(ms)

		body, _ := json.Marshal(UpdateBankCardRequest{Pid: "123", BankName: "MAIB"})
		req := newReqWithChiID(t, http.MethodPut, "/update/nope", "nope", body)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateBankCardObj(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.calledUpdate {
			t.Fatalf("service must NOT be called on invalid id")
		}
	})

	t.Run("missing userID in context -> 422 (service NOT called)", func(t *testing.T) {
		ms := &mockServicE{
			updateFn: func(ctx context.Context, card *domain.BankCard) error {
				t.Fatalf("service must NOT be called when userID missing")
				return nil
			},
		}
		h := New(ms)

		body, _ := json.Marshal(UpdateBankCardRequest{Pid: "123", BankName: "MAIB"})
		req := newReqWithChiID(t, http.MethodPut, "/update/10", "10", body) // no userID

		rr := httptest.NewRecorder()
		h.UpdateBankCardObj(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if ms.calledUpdate {
			t.Fatalf("service must NOT be called when userID missing")
		}
	})

	t.Run("service returns error -> mapped response", func(t *testing.T) {
		svcErr := domain.ErrBankCardNotFound
		wantStatus, wantMsg := errorMapper.Process(svcErr)

		ms := &mockServicE{
			updateFn: func(ctx context.Context, card *domain.BankCard) error {
				return svcErr
			},
		}
		h := New(ms)

		body, _ := json.Marshal(UpdateBankCardRequest{Pid: "123", BankName: "MAIB"})
		req := newReqWithChiID(t, http.MethodPut, "/update/10", "10", body)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.UpdateBankCardObj(rr, req)

		if rr.Code != wantStatus {
			t.Fatalf("expected %d, got %d, body=%s", wantStatus, rr.Code, rr.Body.String())
		}
		if wantMsg != "" && !bytes.Contains(rr.Body.Bytes(), []byte(wantMsg)) {
			t.Fatalf("expected body to contain %q, got %s", wantMsg, rr.Body.String())
		}
		if !ms.calledUpdate {
			t.Fatalf("expected service to be called")
		}
	})

	t.Run("ok -> 200 and service called with correct domain", func(t *testing.T) {
		ms := &mockServicE{
			updateFn: func(ctx context.Context, card *domain.BankCard) error {
				return nil
			},
		}
		h := New(ms)

		body, _ := json.Marshal(UpdateBankCardRequest{Pid: "PID-999", BankName: "MICB"})
		req := newReqWithChiID(t, http.MethodPut, "/update/55", "55", body)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(777)))

		rr := httptest.NewRecorder()
		h.UpdateBankCardObj(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}
		if !ms.calledUpdate {
			t.Fatalf("expected service to be called")
		}
		if ms.lastCard == nil {
			t.Fatalf("expected lastCard, got nil")
		}
		if ms.lastCard.CardId != 55 {
			t.Fatalf("expected CardId=55, got %d", ms.lastCard.CardId)
		}
		if ms.lastCard.UserId != 777 {
			t.Fatalf("expected UserId=777, got %d", ms.lastCard.UserId)
		}
		if ms.lastCard.Pid != "PID-999" {
			t.Fatalf("expected Pid=PID-999, got %q", ms.lastCard.Pid)
		}
		if ms.lastCard.Bank != "MICB" {
			t.Fatalf("expected Bank=MICB, got %q", ms.lastCard.Bank)
		}
	})
}
