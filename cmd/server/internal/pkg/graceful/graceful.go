package graceful

import (
	"context"
	"server/internal/pkg/logger"

	"golang.org/x/sync/errgroup"
)

type graceful struct {
	processes []process
}

func New(processes ...process) *graceful {
	return &graceful{
		processes: processes,
	}
}

func (gr graceful) Start(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, proc := range gr.processes {
		if proc.disabled {
			continue
		}

		f := func() error {
			err := proc.starter.Start(gCtx)
			if err != nil {
				logger.Log.Error(err.Error())
				logger.Log.Info("GracefulShutDown started")
				return err
			}
			return nil
		}

		g.Go(f)
	}


	err := g.Wait()
	if err != nil {
		logger.Log.Error("Application stopped gracefully")

		return err
	}

	logger.Log.Info("Every process stopped by itself with no error")

	return nil
}
