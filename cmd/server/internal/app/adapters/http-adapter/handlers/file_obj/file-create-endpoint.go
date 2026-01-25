package file_obj

import (
	"encoding/json"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	domain "server/internal/app/domain/file_obj"
)

type createFileRequest struct {
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	BucketName  string `json:"bucket_name"`
	ObjectKey   string `json:"object_key"`
	SizeBytes   int64  `json:"size_bytes"`
	ContentType string `json:"content_type"`
	ETag        string `json:"etag"`
}

type createFileResponse struct {
	ID int64 `json:"id"`
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid json")
		return
	}

	ref, err := domain.NewStorageRef(req.BucketName, req.ObjectKey)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	file, err := domain.NewFile(
		domain.UserID(req.UserID),
		req.Title,
		ref,
		req.SizeBytes,
		req.ContentType,
		req.ETag,
	)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.uc.Create(r.Context(), file)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusInternalServerError, "internal error")
		return
	}

	codec.WriteJSON(w, http.StatusCreated, createFileResponse{ID: id})
}
