package codec

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON_OK(t *testing.T) {
	rr := httptest.NewRecorder()

	body := map[string]string{
		"hello": "world",
	}

	WriteJSON(rr, http.StatusOK, body)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var decoded map[string]string
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded["hello"] != "world" {
		t.Fatalf("unexpected body: %+v", decoded)
	}
}

func TestWriteJSON_NilBody(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteJSON(rr, http.StatusNoContent, nil)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	// body должен быть пустым
	if rr.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", rr.Body.String())
	}
}

func TestWriteErrorJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteErrorJSON(rr, http.StatusBadRequest, "something went wrong")

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
	}

	var decoded map[string]string
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded["error"] != "something went wrong" {
		t.Fatalf("unexpected error message: %+v", decoded)
	}
}
