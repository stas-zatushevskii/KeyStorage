package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	domain "server/internal/app/domain/file"
)

type Repository interface {
	Create(ctx context.Context, f *domain.File) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.File, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.File, error)
	Delete(ctx context.Context, id int64) error
}

type ObjectStorage interface {
	PutObject(ctx context.Context, bucket string, key string, body io.Reader, size int64, contentType string) (string, error)
	DeleteObject(ctx context.Context, bucket string, key string) error
	GetObjectReader(
		ctx context.Context,
		bucket string,
		objectKey string,
	) (io.ReadCloser, error)
}

type File struct {
	repo    Repository
	storage ObjectStorage
}

func New(repo Repository, storage ObjectStorage) *File {
	return &File{repo: repo, storage: storage}
}

func (u *File) GetByID(ctx context.Context, fileID int64) (*domain.File, error) {
	id := fileID
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

func (u *File) GetFileList(ctx context.Context, userID int64) ([]*domain.File, error) {
	uid := userID
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

func (u *File) UploadAndCreate(ctx context.Context, file *domain.File, data []byte) (int64, error) {
	if file == nil {
		return 0, fmt.Errorf("file is nil")
	}
	if u.storage == nil {
		return 0, fmt.Errorf("storage is nil")
	}
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

	file.ETag = etag

	id, err := u.repo.Create(ctx, file)
	if err != nil {
		// rollback storage
		_ = u.storage.DeleteObject(ctx, file.Storage.BucketName, file.Storage.ObjectKey)
		return 0, fmt.Errorf("create file meta: %w", err)
	}

	return id, nil
}

func (u *File) GetFileStream(ctx context.Context, userID, fileID int64) (*domain.File, io.ReadCloser, error) {

	f, err := u.repo.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return nil, nil, domain.ErrFileNotFound
		}
	}

	if f.UserID != userID {
		return nil, nil, domain.ErrInvalidUserID
	}

	rc, err := u.storage.GetObjectReader(ctx, f.Storage.BucketName, f.Storage.ObjectKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get object: %w", err)
	}

	return f, rc, nil
}
