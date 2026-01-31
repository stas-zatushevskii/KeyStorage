package minio

import (
	"context"
	"fmt"
	"server/internal/app/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc *minio.Client
}

func NewConnection() (*Client, error) {
	cl, err := minio.New(
		config.App.GetMinioEndpoint(),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.App.GetMinioAccessKey(),
				config.App.GetMinioSecretKey(),
				"",
			),
			Secure: config.App.GetMinioUseSSL(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	return &Client{mc: cl}, nil
}

func (c *Client) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Проверка наличия бакета и его создание, если не существует
	exists, err := c.mc.BucketExists(ctx, config.App.GetMinioBucketName())
	if err != nil {
		return err
	}
	if !exists {
		err := c.mc.MakeBucket(ctx, config.App.GetMinioBucketName(), minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
