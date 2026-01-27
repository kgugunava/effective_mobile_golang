package service

import (
	"context"

	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription postgres.SubscriptionEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (postgres.SubscriptionEntity, error)
	UpdatePut(ctx context.Context, sub postgres.SubscriptionEntity, id uuid.UUID) (postgres.SubscriptionEntity, error)
	UpdatePatch(ctx context.Context, id uuid.UUID, changes map[string]interface{}) (postgres.SubscriptionEntity, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}