package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6noz/services/pkg/redis"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

func main() {
	ctx := context.Background()

	logger, err := zapx.NewLogger()
	panicOnErr(err)
	defer logger.Close()

	redisClient, err := redis.NewFailoverClient(ctx, &redis.Config{},
		redis.WithMetrics(),
		redis.WithTraces(),
	)
	panicOnErr(err)
	defer redisClient.Close()

	// locationRepo := location.NewRepository(redisClient)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	defer e.Close()

	err = e.Start(":8080")
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
