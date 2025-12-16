package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	goredis "github.com/redis/go-redis/v9"
	"github.com/tuanta7/k6noz/services/pkg/redis"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

func main() {
	ctx := context.Background()

	logger, err := zapx.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	redisClient, err := redis.NewClusterClient(ctx, &goredis.FailoverOptions{})
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	// locationRepo := location.NewRepository(redisClient)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
