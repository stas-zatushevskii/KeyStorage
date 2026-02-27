package text_obj

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/app/adapters/primary/http-adapter/constants"
	"testing"

	textDomain "server/internal/app/domain/text_obj"
)

type mockServices struct {
	getListFn func(ctx context.Context, userId int64) ([]*textDomain.Text, error)
	calls     int
	lastUID   int64
}

func (m *mockServices) GetText(ctx context.Context, cardId int64) (*textDomain.Text, error) {
	panic("not used in these tests")
}
func (m *mockServices) CreateNewTextObj(ctx context.Context, card *textDomain.Text) (int64, error) {
	panic("not used in these tests")
}
func (m *mockServices) UpdateText(ctx context.Context, card *textDomain.Text) error {
	panic("not used in these tests")
}

func (m *mockServices) GetTextList(ctx context.Context, userId int64) ([]*textDomain.Text, error) {
	m.calls++
	m.lastUID = userId
	return m.getListFn(ctx, userId)
}

func TestHttpHandler_GetTextList(t *testing.T) {
	t.Run("missing userID in context -> 422 (BUG: service is still called)", func(t *testing.T) {
		ms := &mockServices{
			getListFn: func(ctx context.Context, userId int64) ([]*textDomain.Text, error) {
				// handler вызовет сервис даже при отсутствии userID, userId будет 0
				return nil, errors.New("service called unexpectedly")
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		rr := httptest.NewRecorder()

		h.GetTextList(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status 422, got %d, body=%s", rr.Code, rr.Body.String())
		}

		// ВАЖНО: это демонстрирует текущий баг — сервис реально вызван.
		if ms.calls != 1 {
			t.Fatalf("expected service called once due to BUG, got %d", ms.calls)
		}
		if ms.lastUID != 0 {
			t.Fatalf("expected userId=0 due to missing context, got %d", ms.lastUID)
		}
	})

	t.Run("service returns error -> status != 200 (mapped)", func(t *testing.T) {
		ms := &mockServices{
			getListFn: func(ctx context.Context, userId int64) ([]*textDomain.Text, error) {
				return nil, errors.New("boom")
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.GetTextList(rr, req)

		if ms.calls != 1 {
			t.Fatalf("expected service called once, got %d", ms.calls)
		}
		if rr.Code == http.StatusOK {
			t.Fatalf("expected status != 200, got 200, body=%s", rr.Body.String())
		}

		// Тело зависит от codec/error-mapper, поэтому тут проверяем только что не пустое.
		if rr.Body.Len() == 0 {
			t.Fatalf("expected non-empty error body")
		}
	})

	t.Run("ok -> 200 and json list", func(t *testing.T) {
		ms := &mockServices{
			getListFn: func(ctx context.Context, userId int64) ([]*textDomain.Text, error) {
				return []*textDomain.Text{
					{TextId: 1, Title: "A"},
					{TextId: 2, Title: "B"},
				}, nil
			},
		}
		h := &HttpHandler{service: ms}

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDKey, int64(7)))

		rr := httptest.NewRecorder()
		h.GetTextList(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
		}

		var got []Text
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v, body=%s", err, rr.Body.String())
		}

		if len(got) != 2 {
			t.Fatalf("expected 2 items, got %d: %+v", len(got), got)
		}
		if got[0].TextID != 1 || got[0].Title != "A" {
			t.Fatalf("unexpected item[0]: %+v", got[0])
		}
		if got[1].TextID != 2 || got[1].Title != "B" {
			t.Fatalf("unexpected item[1]: %+v", got[1])
		}
	})
}
