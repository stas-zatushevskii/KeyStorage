package app

import (
	"context"
	"fmt"
	"server/internal/app/adapters/http-adapter"
	os_signal_adapter "server/internal/app/adapters/os-signal-adapter"
	accountRepository "server/internal/app/repository/accout_obj"
	bankCardRepository "server/internal/app/repository/bank_card_obj"
	fileRepository "server/internal/app/repository/file_obj"
	textRepository "server/internal/app/repository/text_obj"
	userReporitory "server/internal/app/repository/user"
	accountUsecase "server/internal/app/usecases/account_obj"
	bankCardUsecase "server/internal/app/usecases/bank_card_obj"
	fileUsecase "server/internal/app/usecases/file_obj"
	textUsecase "server/internal/app/usecases/text_obj"
	userUsecase "server/internal/app/usecases/user"
	"server/internal/pkg/graceful"
	"server/internal/pkg/minio"
	postgres "server/internal/pkg/postgres"
)

type App struct {
	HttpAdapter     *http_adapter.HttpAdapter
	OSSignalAdapter *os_signal_adapter.OsSignalAdapter
}

func New() (*App, error) {

	// postgres
	database, err := postgres.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err = postgres.SetupDatabase(database); err != nil {
		return nil, fmt.Errorf("failed to setup database: %v", err)
	}

	// minio
	fileStorage, err := minio.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to setup minio connection: %v", err)
	}
	if err = fileStorage.InitMinio(); err != nil {
		return nil, fmt.Errorf("failed to setup minio connection: %v", err)
	}

	// os signals
	osSignalAdapter := os_signal_adapter.New()

	// http
	httpAdapter := http_adapter.New(&http_adapter.Srv{
		UserUseCase:        userUsecase.New(userReporitory.New(database)),
		AccountObjUseCase:  accountUsecase.New(accountRepository.New(database)),
		BankCardObjUseCase: bankCardUsecase.New(bankCardRepository.New(database)),
		TextObjUseCase:     textUsecase.New(textRepository.New(database)),
		FileObjUseCase:     fileUsecase.New(fileRepository.New(database), fileStorage),
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
