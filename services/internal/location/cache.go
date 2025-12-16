package location

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/redis"
)

type Cache struct {
	redis redis.Cache
}

func NewRepository(cache redis.Cache) *Cache {
	return &Cache{redis: cache}
}

func (c *Cache) Set(ctx context.Context, location *domain.Location) error {
	return c.redis.Set(ctx, location.TripID, location, 0)
}
