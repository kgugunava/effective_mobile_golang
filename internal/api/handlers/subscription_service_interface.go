package handlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	UpdateSubscriptionPut(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error)
	UpdateSubscriptionPatch(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error) 
	DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error
}