package file_obj

import (
	"io"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	domain "server/internal/app/domain/file_obj"

	"github.com/google/uuid"
)

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 10 << 20 // 10 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	userID, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	objectKey := uuid.New().String()

	if title == "" {
		title = fileHeader.Filename
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	sizeBytes := fileHeader.Size

	// add bucket name and unique key
	ref, err := domain.NewStorageRef("user-files", objectKey)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	f, err := domain.NewFile(
		userID,
		title,
		ref,
		sizeBytes,
		contentType,
	)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.uc.UploadAndCreate(r.Context(), f, fileBytes)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	codec.WriteJSON(w, http.StatusCreated, "file uploaded successfully")
}
