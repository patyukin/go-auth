package service

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type (
	Service interface {
		Init(ctx context.Context) error
		Run(ctx context.Context) error
		Stop()
	}
	Services interface {
		AddService(service ...Service)
		Run(ctx context.Context) error
	}
	Manager struct {
		log      *slog.Logger
		services []Service
	}
)

func NewManager(log *slog.Logger) Services {
	return &Manager{log: log}
}

func (s *Manager) AddService(service ...Service) {
	s.services = append(s.services, service...)
}

func (s *Manager) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.services))

	s.log.Info("going to start services")
	for _, service := range s.services {
		err := service.Init(ctx)
		if err != nil {
			s.log.ErrorContext(
				ctx,
				"service initialization failed",
				slog.String("err_msg", err.Error()),
			)
			continue
		}
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			if err := svc.Run(ctx); err != nil {
				errChan <- err
			}
		}(service)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case <-c:
		s.stop()
	case <-ctx.Done():
		s.stop()
	case err := <-errChan:
		s.stop()
		return err
	}

	return nil
}

func (s *Manager) stop() {
	s.log.Info("going to stop")
	for _, service := range s.services {
		service.Stop()
	}
}
