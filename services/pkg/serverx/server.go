package serverx

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Run() error
	Shutdown(context.Context) error
}

func RunServer(server Server, gracePeriod ...time.Duration) error {
	errCh := make(chan error)

	go func(s Server) {
		if err := s.Run(); err != nil {
			err = fmt.Errorf("error starting REST server: %w", err)
			errCh <- err
		}
	}(server)

	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Println("Shutting down due to server error:", err)
		return shutdownServer(server, gracePeriod...)
	case <-notifyCh:
		log.Println("Shutting down gracefully...")
		_ = shutdownServer(server, gracePeriod...)
		return nil
	}
}

func shutdownServer(server Server, gracePeriod ...time.Duration) error {
	timeout := 20 * time.Second
	if len(gracePeriod) > 0 {
		timeout = gracePeriod[0]
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("Error during server shutdown:", err)
		return err
	}

	return nil
}
