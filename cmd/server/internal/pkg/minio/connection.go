package minio

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/app/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/sync/errgroup"
)

type MinioAdapter struct {
	CL        *minio.Client
	transport *http.Transport
}

func New() (*MinioAdapter, error) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxConnsPerHost:     10,
		MaxIdleConnsPerHost: 10,
	}

	cl, err := minio.New(
		config.App.GetMinioEndpoint(),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.App.GetMinioAccessKey(),
				config.App.GetMinioSecretKey(),
				"",
			),
			Secure:    config.App.GetMinioUseSSL(),
			Transport: transport,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	return &MinioAdapter{CL: cl, transport: transport}, nil
}

func (mc *MinioAdapter) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {

		err := mc.InitMinio()
		if err != nil {
			return fmt.Errorf("failed to setup buckets: %v", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		mc.transport.CloseIdleConnections()
		return nil
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}
