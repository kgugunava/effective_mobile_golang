package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

type SubscriptionService struct {
	subscriptionRepo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: repo,
	}
}

func isSubscriptionValid(subscription *domain.Subscription) bool {
	if subscription.Price <= 0 {
		return false
	}
	if subscription.StartDate.IsZero() {
		return false
	}
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		return false
	}

	return true
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error) {
	subscription.SubscriptionID = uuid.New()
	if isSubscriptionValid(subscription) {
		if err := s.subscriptionRepo.Create(ctx, transferServiceDomainToPostgresEntity(*subscription)); err != nil {
			return nil, fmt.Errorf("repository create failed: %w", err)
		}
	}
	return subscription, nil
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Subscription{}, err
	}
	return transferPostgresEntityToServiceDomain(subscription), nil
}