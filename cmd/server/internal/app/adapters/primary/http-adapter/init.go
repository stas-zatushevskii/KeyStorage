package http_adapter

import (
	"context"
	account_router "server/internal/app/adapters/primary/http-adapter/handlers/account_obj"
	bankCard_router "server/internal/app/adapters/primary/http-adapter/handlers/bank_card_obj"
	file_router "server/internal/app/adapters/primary/http-adapter/handlers/file_obj"
	text_router "server/internal/app/adapters/primary/http-adapter/handlers/text_obj"
	user_router "server/internal/app/adapters/primary/http-adapter/handlers/user"
	"server/internal/app/adapters/primary/http-adapter/middlewares"
	account "server/internal/app/usecases/account_obj"
	bankCard "server/internal/app/usecases/bank_card_obj"
	file "server/internal/app/usecases/file_obj"
	text "server/internal/app/usecases/text_obj"
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
	TextObjUseCase     *text.TextObj
	FileObjUseCase     *file.FileObj
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

	// text handler
	textRouter := text_router.New(srv.TextObjUseCase)

	// file handler
	fileRouter := file_router.New(srv.FileObjUseCase)

	// create router
	r := chi.NewRouter()

	// mount user router
	r.Mount("/user", userRouter.Routes(srv.UserUseCase))

	// mount account object router with jwt authentification middleware
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/account", accountRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/card", bankCardRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/text", textRouter.Routes())
	r.With(middlewares.JWTMiddleware(srv.UserUseCase)).Mount("/file", fileRouter.Routes())

	return r
}

func (a HttpAdapter) Start(ctx context.Context) error {
	return a.server.Start(ctx)
}
