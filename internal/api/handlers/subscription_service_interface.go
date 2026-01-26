package handlers

import (
	"context"
	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error)
}