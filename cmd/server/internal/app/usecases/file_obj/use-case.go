package file_obj

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	domain "server/internal/app/domain/file_obj"
)

type Repository interface {
	Create(ctx context.Context, f *domain.File) (domain.FileID, error)
	GetByID(ctx context.Context, id domain.FileID) (*domain.File, error)
	ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.File, error)
	Delete(ctx context.Context, id domain.FileID) error
}

type ObjectStorage interface {
	PutObject(ctx context.Context, bucket, key string, body *bytes.Reader, size int64, contentType string) (etag string, err error)
	DeleteObject(ctx context.Context, bucket, key string) error
}

type UseCase struct {
	repo    Repository
	storage ObjectStorage
}

func New(repo Repository, storage ObjectStorage) *UseCase {
	return &UseCase{repo: repo, storage: storage}
}

func (u *UseCase) GetByID(ctx context.Context, fileID int64) (*domain.File, error) {
	id := domain.FileID(fileID)
	if id <= 0 {
		return nil, domain.ErrInvalidFileID
	}

	file, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get file by id=%d: %w", fileID, err)
	}

	return file, nil
}

func (u *UseCase) GetFileList(ctx context.Context, userID int64) ([]*domain.File, error) {
	uid := domain.UserID(userID)
	if uid <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	list, err := u.repo.ListByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("list files by user_id=%d: %w", userID, err)
	}

	if len(list) == 0 {
		return nil, domain.ErrEmptyFilesList
	}

	return list, nil
}

// UploadAndCreate
// 1) PutObject в MinIO
// 2) Create метаданные в Postgres
// 3) Если Create упал -> DeleteObject (rollback)
func (u *UseCase) UploadAndCreate(ctx context.Context, file *domain.File, data []byte) (int64, error) {
	if file == nil {
		return 0, fmt.Errorf("file is nil")
	}
	if u.storage == nil {
		return 0, fmt.Errorf("storage is nil")
	}

	// size_bytes из данных (если не заполнен/или ты хочешь всегда доверять факту)
	if file.SizeBytes == 0 {
		file.SizeBytes = int64(len(data))
	}
	if file.SizeBytes < 0 {
		return 0, fmt.Errorf("size_bytes must be >= 0")
	}

	// content-type если не задан
	if file.ContentType == "" {
		file.ContentType = http.DetectContentType(data)
	}

	// 1) Upload to MinIO (so3)
	etag, err := u.storage.PutObject(
		ctx,
		file.Storage.BucketName,
		file.Storage.ObjectKey,
		bytes.NewReader(data),
		file.SizeBytes,
		file.ContentType,
	)
	if err != nil {
		return 0, fmt.Errorf("upload to storage bucket=%s key=%s: %w",
			file.Storage.BucketName, file.Storage.ObjectKey, err)
	}

	// сохраняем etag в домен
	file.ETag = etag

	// 2) Save metadata to Postgres
	id, err := u.repo.Create(ctx, file)
	if err != nil {
		// 3) rollback storage (best effort)
		_ = u.storage.DeleteObject(ctx, file.Storage.BucketName, file.Storage.ObjectKey)
		return 0, fmt.Errorf("create file meta: %w", err)
	}

	return int64(id), nil
}

func (u *UseCase) Create(ctx context.Context, file *domain.File) (int64, error) {
	if file == nil {
		return 0, fmt.Errorf("file is nil")
	}

	id, err := u.repo.Create(ctx, file)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	return int64(id), nil
}

func (u *UseCase) Delete(ctx context.Context, fileID int64) error {
	id := domain.FileID(fileID)
	if id <= 0 {
		return domain.ErrInvalidFileID
	}

	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return err
		}
		return fmt.Errorf("delete file id=%d: %w", fileID, err)
	}

	return nil
}
