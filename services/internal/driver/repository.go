package driver

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository struct {
	mongo mongo.DB
}

func NewRepository(mongo mongo.DB) *Repository {
	return &Repository{mongo: mongo}
}

func (r *Repository) GetDriverByID(ctx context.Context, driverID string) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.mongo.Get(ctx, "drivers", bson.D{{"id", driverID}}, &driver)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}
