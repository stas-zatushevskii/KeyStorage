package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func (c *Client) PutObject(
	ctx context.Context,
	bucket, key string,
	body io.Reader,
	size int64,
	contentType string,
) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := c.mc.PutObject(ctx, bucket, key, body, size, opts)
	if err != nil {
		return "", fmt.Errorf("minio put object bucket=%s key=%s: %w", bucket, key, err)
	}

	return info.ETag, nil
}

func (c *Client) DeleteObject(ctx context.Context, bucket, key string) error {
	opts := minio.RemoveObjectOptions{}

	if err := c.mc.RemoveObject(ctx, bucket, key, opts); err != nil {
		return fmt.Errorf("minio remove object bucket=%s key=%s: %w", bucket, key, err)
	}

	return nil
}
