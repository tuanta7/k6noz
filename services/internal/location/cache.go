package location

import (
	"context"
	"encoding/json"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/redis"
)

type Cache struct {
	redis redis.Cache
}

func NewCache(cache redis.Cache) *Cache {
	return &Cache{redis: cache}
}

func (c *Cache) GetLocation(ctx context.Context, driverID string) (*domain.Location, error) {
	data, err := c.redis.Get(ctx, driverID)
	if err != nil {
		return nil, err
	}

	var location domain.Location
	err = json.Unmarshal(data, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}
