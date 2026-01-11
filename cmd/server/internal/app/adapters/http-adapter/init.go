package http_adapter

import (
	"context"
	account_router "server/internal/app/adapters/http-adapter/handlers/account_obj"
	bankCard_router "server/internal/app/adapters/http-adapter/handlers/bank_card_obj"
	user_router "server/internal/app/adapters/http-adapter/handlers/user"

	"server/internal/app/adapters/http-adapter/middlewares"
	account "server/internal/app/usecases/account_obj"
	bankCard "server/internal/app/usecases/bank_card_obj"
	"server/internal/app/usecases/user"
	http_server "server/internal/pkg/http-server"

	"github.com/go-chi/chi/v5"
)

type HttpAdapter struct {
	server *http_server.Server
}

type Srv struct {
	UserUseCase        *user.User
	AccountObjUseCase  *account.AccountObj
	BankCardObjUseCase *bankCard.BankCardObj
}

func New(svc *Srv) *HttpAdapter {
	router := newRouter(svc)

	s := http_server.New(router)

	return &HttpAdapter{
		server: s,
	}
}

func newRouter(srv *Srv) *chi.Mux {
	// user router
	userRouter := user_router.New(srv.UserUseCase)

	// account router
	accountRouter := account_router.New(srv.AccountObjUseCase)

	// bank card router
	bankCardRouter := bankCard_router.New(srv.BankCardObjUseCase)

	// create router
	r := chi.NewRouter()

	// mount user router
	r.Mount("/user", userRouter)

	// mount account object router with jwt authentification middleware
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/account", accountRouter)
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/card", bankCardRouter)

	return r
}

func (a HttpAdapter) Start(ctx context.Context) error {
	return a.server.Start(ctx)
}
