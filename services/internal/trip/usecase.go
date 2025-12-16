package trip

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
)

type UseCase struct {
	cache  Cache
	logger *zapx.Logger
}

func (u *UseCase) GetNearbyDrivers(ctx context.Context, location *domain.Location) ([]*domain.Driver, error) {
	return nil, nil
}

func (u *UseCase) GetDriverLatestLocation(ctx context.Context, driverID string) (*domain.Location, error) {
	return nil, nil
}
