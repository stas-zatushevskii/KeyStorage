package http_adapter

import (
	"context"
	accountrouter "server/internal/app/adapters/primary/http-adapter/handlers/account"
	bankCardrouter "server/internal/app/adapters/primary/http-adapter/handlers/bank_card"
	filerouter "server/internal/app/adapters/primary/http-adapter/handlers/file"
	textrouter "server/internal/app/adapters/primary/http-adapter/handlers/text"
	userrouter "server/internal/app/adapters/primary/http-adapter/handlers/user"
	"server/internal/app/adapters/primary/http-adapter/middlewares"
	account "server/internal/app/usecases/account"
	bankCard "server/internal/app/usecases/bank_card"
	file "server/internal/app/usecases/file"
	text "server/internal/app/usecases/text"
	"server/internal/app/usecases/user"
	httpserver "server/internal/pkg/http-server"

	"github.com/go-chi/chi/v5"
)

type HttpAdapter struct {
	server *httpserver.Server
}

type Srv struct {
	UserUseCase     *user.User
	AccountUseCase  *account.Account
	BankCardUseCase *bankCard.BankCard
	TextUseCase     *text.Text
	FileUseCase     *file.File
}

func New(svc *Srv) *HttpAdapter {
	router := newRouter(svc)

	s := httpserver.New(router)

	return &HttpAdapter{
		server: s,
	}
}

func newRouter(srv *Srv) *chi.Mux {
	// user router
	userRouter := userrouter.New(srv.UserUseCase)

	// account router
	accountRouter := accountrouter.New(srv.AccountUseCase)

	// bank card router
	bankCardRouter := bankCardrouter.New(srv.BankCardUseCase)

	// text handler
	textRouter := textrouter.New(srv.TextUseCase)

	// file handler
	fileRouter := filerouter.New(srv.FileUseCase)

	// create router
	r := chi.NewRouter()

	// mount user router
	r.Mount("/user", userRouter.Routes(srv.UserUseCase))

	// mount account ect router with jwt authentification middleware
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/account", accountRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/card", bankCardRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/text", textRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/file", fileRouter.Routes())

	return r
}

func (a HttpAdapter) Start(ctx context.Context) error {
	return a.server.Start(ctx)
}
