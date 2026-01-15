package file_obj

import (
	"context"
	"fmt"
	domain "server/internal/app/domain/file_obj"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.File, error)
	GetByID(ctx context.Context, fileId int64) (*domain.File, error)
	Create(ctx context.Context, file *domain.File) (int64, error)
	Delete(ctx context.Context, fileId int64) error
}

type FiledObj struct {
	repo Repository
}

func New(repo Repository) *FiledObj {
	return &FiledObj{repo: repo}
}

func (f *FiledObj) GetByID(ctx context.Context, fileId int64) (*domain.File, error) {
	file, err := f.repo.GetByID(ctx, fileId)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return file, nil
}

func (f *FiledObj) GetFileList(ctx context.Context, userID int64) ([]*domain.File, error) {
	list, err := f.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}
	if len(list) == 0 {
		return nil, domain.ErrEmptyFilesList
	}
	return list, nil
}

func (f *FiledObj) Create(ctx context.Context, file *domain.File) (int64, error) {
	if file.UserId == 0 {
		return 0, fmt.Errorf("user id is zero")
	}

	if file.FileName == "" {
		return 0, fmt.Errorf("file name is empty")
	}

	id, err := f.repo.Create(ctx, file)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}

	return id, nil
}

func (f *FiledObj) Delete(ctx context.Context, fileId int64) error {
	err := f.repo.Delete(ctx, fileId)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
