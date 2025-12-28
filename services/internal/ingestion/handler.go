package ingestion

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/kafka"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

type Handler struct {
	publisher *kafka.Publisher
	logger    *zapx.Logger
}

func NewHandler(logger *zapx.Logger, publisher *kafka.Publisher) *Handler {
	return &Handler{
		logger:    logger,
		publisher: publisher,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func(ws *websocket.Conn) {
		if err := ws.Close(); err != nil {
			h.logger.Warn("websocket close error", zap.Error(err))
		}
	}(ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			h.logger.Error("websocket read error", zap.Error(err))
			break
		}

		var location domain.DriverLocation
		if err := json.Unmarshal(msg, &location); err != nil {
			h.logger.Error("invalid payload", zap.Error(err))
			continue
		}

		key := []byte(location.DriverID)
		if err := h.publisher.PublishSync(r.Context(), domain.DriverLocationTopic, key, msg); err != nil {
			h.logger.Error("publish error", zap.Error(err))
			continue
		}
	}
}
