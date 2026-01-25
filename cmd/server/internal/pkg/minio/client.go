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

func NewConnection(cfg Config) (*Client, error) {
	cl, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	return &Client{mc: cl}, nil
}

func (c *Client) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.App.GetMinioEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(config.App.GetMinioRootUser(), config.App.GetMinioRootPassword(), ""),
		Secure: config.App.GetMinioUseSSL(),
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	c.mc = client

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
