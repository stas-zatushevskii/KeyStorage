package http_adapter

import (
	"context"
	UserRouter "server/internal/app/adapters/http-adapter/handlers/user"
	"server/internal/app/usecases/user"
	http_server "server/internal/pkg/http-server"

	"github.com/go-chi/chi/v5"
)

type HttpAdapter struct {
	server *http_server.Server
}

type Svc struct {
	UserUseCase *user.User
	// todo: add more
}

func New(svc *Svc) *HttpAdapter {
	router := newRouter(svc)

	s := http_server.New(router)

	return &HttpAdapter{
		server: s,
	}
}

func newRouter(srv *Svc) *chi.Mux {
	// user router
	userRouter := UserRouter.Routes(srv.UserUseCase)
	// todo: add more

	// create router
	r := chi.NewRouter()

	// mount user router
	r.Mount("/user", userRouter)
	// todo: add more

	return r
}

func (a HttpAdapter) Start(ctx context.Context) error {
	return a.server.Start(ctx)
}
