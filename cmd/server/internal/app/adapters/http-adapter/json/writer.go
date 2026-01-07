package json

import (
	"bytes"
	"encoding/json"
	"net/http"
	"server/internal/pkg/logger"
)

// WriteJSONResponse writing structure (in json format) in response
func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf) // encode data in buffer
	if err := enc.Encode(data); err != nil {
		// if failed to encode data in json format: return status 500
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes()) // write encoded data in response
	if err != nil {
		logger.Log.Error("failed to write response")
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorAsJSON(err error) json.RawMessage {
	if err == nil {
		return nil
	}
	resp := errorResponse{Error: err.Error()}
	response, _ := json.Marshal(resp)
	return response
}
