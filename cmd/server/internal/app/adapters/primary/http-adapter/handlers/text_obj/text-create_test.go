package text_obj

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	"testing"

	domain "server/internal/app/domain/text_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

// mock that implements the whole service interface
type mockService struct {
	// counters
	getTextCalls       int
	getTextListCalls   int
	createTextCalls    int
	updateTextCalls    int
	lastCreateArg      *domain.Text
	lastCreateUserID   int64
	lastCreateTitle    string
	lastCreateTextBody string

	// stubs
	getTextFn     func(ctx context.Context, id int64) (*domain.Text, error)
	getTextListFn func(ctx context.Context, userID int64) ([]*domain.Text, error)
	createFn      func(ctx context.Context, t *domain.Text) (int64, error)
	updateFn      func(ctx context.Context, t *domain.Text) error
}

func (m *mockService) GetText(ctx context.Context, cardId int64) (*domain.Text, error) {
	m.getTextCalls++
	if m.getTextFn == nil {
		return nil, errors.New("GetText not stubbed")
	}
	return m.getTextFn(ctx, cardId)
}

func (m *mockService) GetTextList(ctx context.Context, userId int64) ([]*domain.Text, error) {
	m.getTextListCalls++
	if m.getTextListFn == nil {
		return nil, errors.New("GetTextList not stubbed")
	}
	return m.getTextListFn(ctx, userId)
}

func (m *mockService) CreateNewTextObj(ctx context.Context, card *domain.Text) (int64, error) {
	m.createTextCalls++
	m.lastCreateArg = card
	if card != nil {
		m.lastCreateUserID = card.UserId
		m.lastCreateTitle = card.Title
		m.lastCreateTextBody = card.Text
	}
	if m.createFn == nil {
		return 0, errors.New("CreateNewTextObj not stubbed")
	}
	return m.createFn(ctx, card)
}

func (m *mockService) UpdateText(ctx context.Context, card *domain.Text) error {
	m.updateTextCalls++
	if m.updateFn == nil {
		return errors.New("UpdateText not stubbed")
	}
	return m.updateFn(ctx, card)
}

func TestHttpHandler_CreateText(t *testing.T) {
	logger.Log = zap.NewNop()

	t.Run("bad json -> 422 and service NOT called", func(t *testing.T) {
		svc := &mockService{
			createFn: func(ctx context.Context, t *domain.Text) (int64, error) {
				return 1, nil
			},
		}
		h := New(svc)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("{bad json"))
		rr := httptest.NewRecorder()

		h.CreateText(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d; body=%s", rr.Code, rr.Body.String())
		}
		if svc.createTextCalls != 0 {
			t.Fatalf("service must NOT be called on json decode error")
		}
	})

	t.Run("missing userID in context -> 422 and service NOT called", func(t *testing.T) {
		svc := &mockService{
			createFn: func(ctx context.Context, t *domain.Text) (int64, error) {
				return 1, nil
			},
		}
		h := New(svc)

		body, _ := json.Marshal(CreateTextRequest{Title: "t", Text: "x"})
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		h.CreateText(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d; body=%s", rr.Code, rr.Body.String())
		}
		if svc.createTextCalls != 0 {
			t.Fatalf("service must NOT be called when userID missing")
		}
	})

	t.Run("service returns error -> mapped error response (status != 200)", func(t *testing.T) {
		svc := &mockService{
			createFn: func(ctx context.Context, t *domain.Text) (int64, error) {
				// любая доменная ошибка, которую ваш errorMapper.Process умеет маппить
				return 0, domain.ErrFailedCreateText
			},
		}
		h := New(svc)

		body, _ := json.Marshal(CreateTextRequest{Title: "t", Text: "x"})
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.CreateText(rr, req)

		if rr.Code == http.StatusOK {
			t.Fatalf("expected non-200 status on service error, got 200; body=%s", rr.Body.String())
		}
		if svc.createTextCalls != 1 {
			t.Fatalf("expected service called once, got %d", svc.createTextCalls)
		}
	})

	t.Run("ok -> 200 and returns text_id; passes userId into domain", func(t *testing.T) {
		svc := &mockService{
			createFn: func(ctx context.Context, t *domain.Text) (int64, error) {
				return 123, nil
			},
		}
		h := New(svc)

		body, _ := json.Marshal(CreateTextRequest{Title: "Hello", Text: "World"})
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(42)))

		rr := httptest.NewRecorder()
		h.CreateText(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d; body=%s", rr.Code, rr.Body.String())
		}

		// ensure service received properly built domain object
		if svc.createTextCalls != 1 {
			t.Fatalf("expected service called once, got %d", svc.createTextCalls)
		}
		if svc.lastCreateArg == nil {
			t.Fatalf("expected create arg, got nil")
		}
		if svc.lastCreateUserID != 42 {
			t.Fatalf("expected userId=42 in domain, got %d", svc.lastCreateUserID)
		}
		if svc.lastCreateTitle != "Hello" || svc.lastCreateTextBody != "World" {
			t.Fatalf("unexpected domain: title=%q text=%q", svc.lastCreateTitle, svc.lastCreateTextBody)
		}

		// check response JSON shape
		var resp CreateTextResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("response is not valid json: %v; body=%s", err, rr.Body.String())
		}
		if resp.TextID != 123 {
			t.Fatalf("expected text_id=123, got %d", resp.TextID)
		}
	})
}
