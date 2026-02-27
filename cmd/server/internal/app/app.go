package app

import (
	"context"
	"fmt"
	"server/internal/app/adapters/primary/http-adapter"
	"server/internal/app/adapters/primary/os-signal-adapter"
	fileMinioRepository "server/internal/app/adapters/secondary/repositories/minio/file"
	accountPostgresRepository "server/internal/app/adapters/secondary/repositories/postgrtes/account"
	bankCardPostgresRepository "server/internal/app/adapters/secondary/repositories/postgrtes/bank_card"
	filePostgresRepository "server/internal/app/adapters/secondary/repositories/postgrtes/file"
	textPostgresRepository "server/internal/app/adapters/secondary/repositories/postgrtes/text"
	userPostgresReporitory "server/internal/app/adapters/secondary/repositories/postgrtes/user"
	accountUsecase "server/internal/app/usecases/account"
	bankCardUsecase "server/internal/app/usecases/bank_card"
	fileUsecase "server/internal/app/usecases/file"
	textUsecase "server/internal/app/usecases/text"
	userUsecase "server/internal/app/usecases/user"
	"server/internal/pkg/graceful"
	"server/internal/pkg/minio"
	postgres "server/internal/pkg/postgres"
)

type App struct {
	HttpAdapter     *http_adapter.HttpAdapter
	OSSignalAdapter *os_signal_adapter.OsSignalAdapter
	PostgresAdapter *postgres.DatabaseAdapter
	MinioAdapter    *minio.MinioAdapter
}

func New() (*App, error) {

	// postgres
	p, err := postgres.New()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to p: %v", err)
	}

	// minio
	m, err := minio.New()
	if err != nil {
		return nil, fmt.Errorf("failed to setup minio connection: %v", err)
	}

	// os signals
	osSignalAdapter := os_signal_adapter.New()

	// http
	httpAdapter := http_adapter.New(&http_adapter.Srv{
		UserUseCase:     userUsecase.New(userPostgresReporitory.New(p.DB)),
		AccountUseCase:  accountUsecase.New(accountPostgresRepository.New(p.DB)),
		BankCardUseCase: bankCardUsecase.New(bankCardPostgresRepository.New(p.DB)),
		TextUseCase:     textUsecase.New(textPostgresRepository.New(p.DB)),
		FileUseCase:     fileUsecase.New(filePostgresRepository.New(p.DB), fileMinioRepository.New(m.CL)),
	})

	return &App{
		HttpAdapter:     httpAdapter,
		OSSignalAdapter: osSignalAdapter,
		PostgresAdapter: p,
		MinioAdapter:    m,
	}, nil
}

func (a App) Start() error {
	gr := graceful.New(
		graceful.NewProcess(a.OSSignalAdapter),
		graceful.NewProcess(a.HttpAdapter),
		graceful.NewProcess(a.PostgresAdapter),
		graceful.NewProcess(a.MinioAdapter),
	)

	err := gr.Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}
