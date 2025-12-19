package ingestion

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
)

type Handler struct {
	publisher kafka.Publisher
}

func NewHandler() *Handler {
	return &Handler{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

func (h *Handler) HandleWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	return nil
}
