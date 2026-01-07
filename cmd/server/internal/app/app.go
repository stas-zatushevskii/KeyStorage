package app

import (
	"context"
	"fmt"
	"server/internal/app/adapters/http-adapter"
	os_signal_adapter "server/internal/app/adapters/os-signal-adapter"
	userReporitory "server/internal/app/repository/user"
	userUsecase "server/internal/app/usecases/user"
	"server/internal/pkg/graceful"
	db "server/internal/pkg/postgres"
)

type App struct {
	HttpAdapter *http_adapter.HttpAdapter
}

func New() (*App, error) {

	database, err := db.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("[Db] failed to connect to database: %v", err)
	}
	if err := db.SetupDatabase(database); err != nil {
		return nil, fmt.Errorf("[Db] failed to setup database: %v", err)
	}

	// http
	httpAdapter := http_adapter.New(&http_adapter.Svc{UserUseCase: userUsecase.New(userReporitory.New(database))})
	// todo add grpc

	return &App{
		HttpAdapter: httpAdapter,
	}, nil
}

func (a App) Start() error {
	gr := graceful.New(
		graceful.NewProcess(os_signal_adapter.New()),
		graceful.NewProcess(a.HttpAdapter),
	)

	err := gr.Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}
