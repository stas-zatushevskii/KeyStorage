package app

import (
	"context"
	"fmt"
	"server/internal/app/adapters/http-adapter"
	os_signal_adapter "server/internal/app/adapters/os-signal-adapter"
	accountRepository "server/internal/app/repository/accout_obj"
	userReporitory "server/internal/app/repository/user"
	accountUsecase "server/internal/app/usecases/account_obj"
	userUsecase "server/internal/app/usecases/user"
	"server/internal/pkg/graceful"
	db "server/internal/pkg/postgres"
)

type App struct {
	HttpAdapter     *http_adapter.HttpAdapter
	OSSignalAdapter *os_signal_adapter.OsSignalAdapter
}

func New() (*App, error) {

	database, err := db.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err = db.SetupDatabase(database); err != nil {
		return nil, fmt.Errorf("failed to setup database: %v", err)
	}

	// os signals
	osSignalAdapter := os_signal_adapter.New()

	// http
	httpAdapter := http_adapter.New(&http_adapter.Srv{
		UserUseCase:       userUsecase.New(userReporitory.New(database)),
		AccountObjUseCase: accountUsecase.New(accountRepository.New(database)),
	})

	// todo add grpc

	return &App{
		HttpAdapter:     httpAdapter,
		OSSignalAdapter: osSignalAdapter,
	}, nil
}

func (a App) Start() error {
	gr := graceful.New(
		graceful.NewProcess(a.OSSignalAdapter),
		graceful.NewProcess(a.HttpAdapter),
	)

	err := gr.Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}
