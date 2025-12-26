package slient

import (
	"context"
	"io"

	"github.com/labstack/gommon/log"
)

// Close implements a wrapper for services that require a silent close within a defer statement.
// It should only be used within the main function.
func Close(srv io.Closer) {
	err := srv.Close()
	log.Warnf("Error while closing: %s", err)
}

type CloserWithContext interface {
	Close(ctx context.Context) error
}

func CloseWithContext(srv CloserWithContext, ctx context.Context) {
	err := srv.Close(ctx)
	log.Warnf("Error while closing: %s", err)
}
