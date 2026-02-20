package file_obj

import (
	"io"
	"net/http"
	"server/internal/app/adapters/primary/http-adapter/codec"
	"server/internal/app/adapters/primary/http-adapter/constants"
	"server/internal/app/config"
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

	userId, ok := r.Context().Value(constants.UserIDKey).(int64)
	if !ok {
		codec.WriteErrorJSON(w, http.StatusUnprocessableEntity, "user ID not found in context")
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

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "failed to read file")
		return
	}

	// define MIME type
	detectedType := http.DetectContentType(buf[:n])

	// return cursor to the start
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		codec.WriteErrorJSON(w, http.StatusInternalServerError, "failed to reset file")
		return
	}

	if _, ok = config.App.AllowedMimeSet()[detectedType]; !ok {
		codec.WriteErrorJSON(w, http.StatusBadRequest, "unsupported file type")
		return
	}

	sizeBytes := fileHeader.Size

	// add bucket name and unique key
	ref, err := domain.NewStorageRef("user-files", objectKey)
	if err != nil {
		codec.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	f, err := domain.NewFile(
		userId,
		title,
		ref,
		sizeBytes,
		detectedType,
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
