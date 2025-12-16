package location

import (
	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
)

type Handler struct {
	kafka kafka.Broker
}

func (h *Handler) UpdateDriverLocation(ctx echo.Context) error {
	return nil
}
