package http_server

import (
	"context"
	"errors"
	"net/http"
	"server/internal/app/config"
	"server/internal/pkg/logger"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	server *http.Server
}

func New(router http.Handler) *Server {
	server := &http.Server{
		Handler: router,
		Addr:    config.App.GetServerAddr(),
	}

	s := Server{
		server: server,
	}

	return &s
}

func (a *Server) Start(ctx context.Context) error {
	logger.Log.Info("HTTP server starting on addr: " + config.App.GetServerAddr())

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), config.App.GetShutDownTimeout())
		defer cancel()

		err := a.server.Shutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		err := a.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// ok
			} else {
				return err
			}
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}
