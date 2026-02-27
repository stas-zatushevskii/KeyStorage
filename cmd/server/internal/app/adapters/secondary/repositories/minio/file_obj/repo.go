package file_obj

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func (r *Repository) PutObject(
	ctx context.Context,
	bucket, key string,
	body io.Reader,
	size int64,
	contentType string,
) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := r.mc.PutObject(ctx, bucket, key, body, size, opts)
	if err != nil {
		return "", fmt.Errorf("minio put object bucket=%s key=%s: %w", bucket, key, err)
	}

	return info.ETag, nil
}

func (r *Repository) DeleteObject(ctx context.Context, bucket, key string) error {
	opts := minio.RemoveObjectOptions{}

	if err := r.mc.RemoveObject(ctx, bucket, key, opts); err != nil {
		return fmt.Errorf("minio remove object bucket=%s key=%s: %w", bucket, key, err)
	}
	return nil
}

func (r *Repository) GetObjectReader(
	ctx context.Context,
	bucket string,
	objectKey string,
) (io.ReadCloser, error) {

	obj, err := r.mc.GetObject(
		ctx,
		bucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		return nil, err
	}

	return obj, nil
}
