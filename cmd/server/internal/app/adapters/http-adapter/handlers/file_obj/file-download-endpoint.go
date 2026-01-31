package file_obj

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	domain "server/internal/app/domain/file_obj"
	"server/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func (h *FileHandler) DownloadByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
		return
	}

	meta, reader, err := h.uc.GetFileStream(r.Context(), userId, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			codec.WriteErrorJSON(w, http.StatusNotFound, "file not found")
		case errors.Is(err, domain.ErrInvalidFileID):
			codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid id")
		default:
			codec.WriteErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	defer reader.Close()

	if meta.SizeBytes > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.SizeBytes))
	}

	if meta.ETag != "" {
		w.Header().Set("ETag", meta.ETag)
	}

	filename := meta.Title
	if filename == "" {
		filename = fmt.Sprintf("file-%d", id)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	w.WriteHeader(http.StatusOK)

	logger.Log.Info(fmt.Sprintf("filename: %s fileID %d userID %d", filename, id, userId))
	if _, err := io.Copy(w, reader); err != nil {
		logger.Log.Error(err.Error())
		return
	}
}
