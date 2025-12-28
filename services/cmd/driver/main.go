package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/tuanta7/k6noz/services/internal/driver"
	"github.com/tuanta7/k6noz/services/pkg/mongo"
	"github.com/tuanta7/k6noz/services/pkg/otelx"
	"github.com/tuanta7/k6noz/services/pkg/serverx"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := driver.LoadConfig()
	slient.PanicOnErr(err)

	logger, err := zapx.NewLogger()
	slient.PanicOnErr(err, "failed to create logger")
	defer slient.Close(logger)

	prometheus, err := otelx.NewPrometheusProvider()
	slient.PanicOnErr(err)

	monitor, err := otelx.NewMonitor(cfg.OTelServiceName, cfg.OTelGRPCEndpoint, prometheus)
	slient.PanicOnErr(err)
	defer slient.CloseWithContext(monitor, ctx)

	err = monitor.SetupOtelSDK(ctx)
	slient.PanicOnErr(err)

	mongoClient, err := mongo.NewClient(ctx, &mongo.Config{
		URI:            cfg.MongoConfig.URI,
		Database:       cfg.MongoConfig.Database,
		ConnectTimeout: cfg.MongoConfig.ConnectTimeout,
		QueryTimeout:   cfg.MongoConfig.QueryTimeout,
		Monitor:        true,
	})
	slient.PanicOnErr(err)
	defer slient.CloseWithContext(mongoClient, ctx)

	repo := driver.NewRepository(mongoClient)
	uc := driver.NewUseCase(logger, repo)
	handler := driver.NewHandler(logger, uc)

	server := driver.NewServer(cfg.BindAddress, handler)
	err = serverx.RunServer(server)
	slient.PanicOnErr(err)
	defer func(server *driver.Server, ctx context.Context) {
		err = server.Shutdown(ctx)
		if err != nil {
			fmt.Println("Error during server shutdown:", err)
		}
	}(server, ctx)
}
