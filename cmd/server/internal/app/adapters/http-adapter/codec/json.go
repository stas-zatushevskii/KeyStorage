package codec

import (
	"encoding/json"
	"net/http"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			logger.Log.Warn("json encode error", zap.Error(err))
		}
	}
	return
}

func WriteErrorJSON(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
	return
}
