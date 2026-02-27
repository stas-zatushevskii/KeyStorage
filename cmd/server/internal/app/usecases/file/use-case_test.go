package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	domain "server/internal/app/domain/file_obj"
)

type repoFake struct {
	create       func(ctx context.Context, f *domain.File) (int64, error)
	getByID      func(ctx context.Context, id int64) (*domain.File, error)
	listByUserID func(ctx context.Context, userID int64) ([]*domain.File, error)
	delete       func(ctx context.Context, id int64) error
}

func (r *repoFake) Create(ctx context.Context, f *domain.File) (int64, error) {
	if r.create != nil {
		return r.create(ctx, f)
	}
	return 0, nil
}
func (r *repoFake) GetByID(ctx context.Context, id int64) (*domain.File, error) {
	if r.getByID != nil {
		return r.getByID(ctx, id)
	}
	return nil, nil
}
func (r *repoFake) ListByUserID(ctx context.Context, userID int64) ([]*domain.File, error) {
	if r.listByUserID != nil {
		return r.listByUserID(ctx, userID)
	}
	return nil, nil
}
func (r *repoFake) Delete(ctx context.Context, id int64) error {
	if r.delete != nil {
		return r.delete(ctx, id)
	}
	return nil
}

type storageFake struct {
	putObject       func(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (string, error)
	deleteObject    func(ctx context.Context, bucket, key string) error
	getObjectReader func(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error)
}

func (s *storageFake) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (string, error) {
	if s.putObject != nil {
		return s.putObject(ctx, bucket, key, body, size, contentType)
	}
	return "", nil
}
func (s *storageFake) DeleteObject(ctx context.Context, bucket, key string) error {
	if s.deleteObject != nil {
		return s.deleteObject(ctx, bucket, key)
	}
	return nil
}
func (s *storageFake) GetObjectReader(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error) {
	if s.getObjectReader != nil {
		return s.getObjectReader(ctx, bucket, objectKey)
	}
	return nil, nil
}

type nopCloser struct{ io.Reader }

func (n nopCloser) Close() error { return nil }

func TestFileObj_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid fileID -> ErrInvalidFileID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{}, &storageFake{})
		_, err := uc.GetByID(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidFileID) {
			t.Fatalf("expected ErrInvalidFileID, got: %v", err)
		}
	})

	t.Run("repo returns ErrFileNotFound -> return it as-is", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return nil, domain.ErrFileNotFound
			},
		}, &storageFake{})

		_, err := uc.GetByID(ctx, 10)
		if !errors.Is(err, domain.ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got: %v", err)
		}
	})

	t.Run("repo returns other error -> wrapped with context", func(t *testing.T) {
		t.Parallel()

		dbErr := errors.New("db timeout")

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return nil, dbErr
			},
		}, &storageFake{})

		_, err := uc.GetByID(ctx, 99)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "get file by id=99") {
			t.Fatalf("expected error message to contain context, got: %v", err)
		}
	})

	t.Run("ok -> returns file", func(t *testing.T) {
		t.Parallel()

		want := &domain.File{ID: 1, UserID: 2}

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return want, nil
			},
		}, &storageFake{})

		got, err := uc.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if got != want {
			t.Fatalf("expected same pointer")
		}
	})
}

func TestFileObj_GetFileList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid userID -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{}, &storageFake{})
		_, err := uc.GetFileList(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("repo error -> wrapped with context", func(t *testing.T) {
		t.Parallel()

		dbErr := errors.New("select failed")

		uc := New(&repoFake{
			listByUserID: func(ctx context.Context, userID int64) ([]*domain.File, error) {
				return nil, dbErr
			},
		}, &storageFake{})

		_, err := uc.GetFileList(ctx, 7)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "list files by user_id=7") {
			t.Fatalf("expected context in error, got: %v", err)
		}
	})

	t.Run("empty list -> ErrEmptyFilesList", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			listByUserID: func(ctx context.Context, userID int64) ([]*domain.File, error) {
				return []*domain.File{}, nil
			},
		}, &storageFake{})

		_, err := uc.GetFileList(ctx, 7)
		if !errors.Is(err, domain.ErrEmptyFilesList) {
			t.Fatalf("expected ErrEmptyFilesList, got: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		want := []*domain.File{{ID: 1, UserID: 7}, {ID: 2, UserID: 7}}

		uc := New(&repoFake{
			listByUserID: func(ctx context.Context, userID int64) ([]*domain.File, error) {
				return want, nil
			},
		}, &storageFake{})

		got, err := uc.GetFileList(ctx, 7)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("expected len=%d got len=%d", len(want), len(got))
		}
		if got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("expected same pointers in slice")
		}
	})
}

func TestFileObj_UploadAndCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	baseFile := func() *domain.File {
		return &domain.File{
			UserID:      1,
			SizeBytes:   3,
			ContentType: "text/plain",
			Storage: domain.StorageRef{
				BucketName: "b",
				ObjectKey:  "k",
			},
		}
	}

	t.Run("file nil -> error", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{}, &storageFake{})
		_, err := uc.UploadAndCreate(ctx, nil, []byte("abc"))
		if err == nil || !strings.Contains(err.Error(), "file is nil") {
			t.Fatalf("expected 'file is nil' error, got: %v", err)
		}
	})

	t.Run("storage nil -> error", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{}, nil)
		_, err := uc.UploadAndCreate(ctx, baseFile(), []byte("abc"))
		if err == nil || !strings.Contains(err.Error(), "storage is nil") {
			t.Fatalf("expected 'storage is nil' error, got: %v", err)
		}
	})

	t.Run("PutObject error -> wrapped with bucket/key", func(t *testing.T) {
		t.Parallel()

		putErr := errors.New("s3 down")

		uc := New(&repoFake{}, &storageFake{
			putObject: func(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (string, error) {
				return "", putErr
			},
		})

		_, err := uc.UploadAndCreate(ctx, baseFile(), []byte("abc"))
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, putErr) {
			t.Fatalf("expected wrapped putErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "upload to storage bucket=b key=k") {
			t.Fatalf("expected context in error, got: %v", err)
		}
	})

	t.Run("repo.Create error -> rollback DeleteObject called, error wrapped", func(t *testing.T) {
		t.Parallel()

		createErr := errors.New("insert failed")

		deleted := false

		uc := New(&repoFake{
			create: func(ctx context.Context, f *domain.File) (int64, error) {
				return 0, createErr
			},
		}, &storageFake{
			putObject: func(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (string, error) {
				b, _ := io.ReadAll(body)
				if string(b) != "abc" {
					t.Fatalf("expected body=abc, got %q", string(b))
				}
				return "etag-1", nil
			},
			deleteObject: func(ctx context.Context, bucket, key string) error {
				if bucket != "b" || key != "k" {
					t.Fatalf("expected delete bucket=b key=k, got bucket=%s key=%s", bucket, key)
				}
				deleted = true
				return nil
			},
		})

		f := baseFile()
		_, err := uc.UploadAndCreate(ctx, f, []byte("abc"))
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, createErr) {
			t.Fatalf("expected wrapped createErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "create file meta") {
			t.Fatalf("expected context in error, got: %v", err)
		}
		if !deleted {
			t.Fatalf("expected rollback DeleteObject to be called")
		}
		if f.ETag != "etag-1" {
			t.Fatalf("expected ETag=etag-1, got %q", f.ETag)
		}
	})

	t.Run("ok -> returns id and sets ETag", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, f *domain.File) (int64, error) {
				if f.ETag != "etag-ok" {
					t.Fatalf("expected ETag=etag-ok, got %q", f.ETag)
				}
				return 777, nil
			},
		}, &storageFake{
			putObject: func(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (string, error) {
				if bucket != "b" || key != "k" {
					t.Fatalf("expected bucket=b key=k, got bucket=%s key=%s", bucket, key)
				}
				if size != 3 {
					t.Fatalf("expected size=3, got %d", size)
				}
				if contentType != "text/plain" {
					t.Fatalf("expected contentType=text/plain, got %q", contentType)
				}
				got, _ := io.ReadAll(body)
				if !bytes.Equal(got, []byte("abc")) {
					t.Fatalf("expected body=abc, got %q", string(got))
				}
				return "etag-ok", nil
			},
		})

		f := baseFile()
		id, err := uc.UploadAndCreate(ctx, f, []byte("abc"))
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if id != 777 {
			t.Fatalf("expected id=777, got %d", id)
		}
		if f.ETag != "etag-ok" {
			t.Fatalf("expected ETag=etag-ok, got %q", f.ETag)
		}
	})
}

func TestFileObj_GetFileStream(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("repo returns ErrFileNotFound -> ErrFileNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return nil, domain.ErrFileNotFound
			},
		}, &storageFake{})

		_, _, err := uc.GetFileStream(ctx, 1, 10)
		if !errors.Is(err, domain.ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got: %v", err)
		}
	})

	t.Run("user mismatch -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return &domain.File{
					ID:     id,
					UserID: 999,
					Storage: domain.StorageRef{
						BucketName: "b",
						ObjectKey:  "k",
					},
				}, nil
			},
		}, &storageFake{
			getObjectReader: func(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error) {
				t.Fatalf("storage.GetObjectReader must NOT be called on user mismatch")
				return nil, nil
			},
		})

		_, _, err := uc.GetFileStream(ctx, 1, 10)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("storage error -> wrapped", func(t *testing.T) {
		t.Parallel()

		stErr := errors.New("minio down")

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return &domain.File{
					ID:     id,
					UserID: 1,
					Storage: domain.StorageRef{
						BucketName: "b",
						ObjectKey:  "k",
					},
				}, nil
			},
		}, &storageFake{
			getObjectReader: func(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error) {
				return nil, stErr
			},
		})

		_, _, err := uc.GetFileStream(ctx, 1, 10)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, stErr) {
			t.Fatalf("expected wrapped stErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "get object") {
			t.Fatalf("expected context in error, got: %v", err)
		}
	})

	t.Run("ok -> returns file and reader", func(t *testing.T) {
		t.Parallel()

		f := &domain.File{
			ID:     10,
			UserID: 1,
			Storage: domain.StorageRef{
				BucketName: "b",
				ObjectKey:  "k",
			},
		}

		rc := nopCloser{Reader: strings.NewReader("hello")}

		uc := New(&repoFake{
			getByID: func(ctx context.Context, id int64) (*domain.File, error) {
				return f, nil
			},
		}, &storageFake{
			getObjectReader: func(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error) {
				if bucket != "b" || objectKey != "k" {
					t.Fatalf("expected bucket=b key=k, got bucket=%s key=%s", bucket, objectKey)
				}
				return rc, nil
			},
		})

		gotFile, gotRC, err := uc.GetFileStream(ctx, 1, 10)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if gotFile != f {
			t.Fatalf("expected same file pointer")
		}
		if gotRC == nil {
			t.Fatalf("expected non-nil reader")
		}
		defer gotRC.Close()

		b, _ := io.ReadAll(gotRC)
		if string(b) != "hello" {
			t.Fatalf("expected body=hello, got %q", string(b))
		}
	})
}
