package service

import (
	"context"
	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription postgres.Subscription) error
}