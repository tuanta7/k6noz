package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/tuanta7/k6noz/services/internal/ingestion"
	"github.com/tuanta7/k6noz/services/pkg/otelx"
	"github.com/tuanta7/k6noz/services/pkg/serverx"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := ingestion.LoadConfig()
	panicOnErr(err)

	logger, err := zapx.NewLogger()
	panicOnErr(err)
	defer slient.Close(logger)

	prometheus, err := otelx.NewPrometheusProvider()
	panicOnErr(err)

	monitor, err := otelx.NewMonitor(
		"github.com/tuanta7/k6noz/ingestion",
		"",
		prometheus,
	)
	panicOnErr(err)
	defer slient.CloseWithContext(monitor, ctx)

	err = monitor.SetupOtelSDK(ctx)
	panicOnErr(err)

	handler := ingestion.NewHandler(logger)
	server := ingestion.NewServer(cfg, handler, prometheus)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	if err = serverx.RunServer(server); err != nil {
		logger.Error("failed to run server", zap.Error(err))
		return
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
