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

func (s *SubscriptionService) UpdateSubscriptionPut(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error) {

	var updatedSubscription domain.Subscription

	if isSubscriptionValid(newSubscription) {
		updatedSubscriptionPostgresEntity, err := s.subscriptionRepo.UpdatePut(ctx, transferServiceDomainToPostgresEntity(*newSubscription), id)
		if err != nil {
			return nil, err
		}
		updatedSubscription = transferPostgresEntityToServiceDomain(updatedSubscriptionPostgresEntity)
	}

	return &updatedSubscription, nil

}

func (s *SubscriptionService) UpdateSubscriptionPatch(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error) {
	changes := make(map[string]interface{})
	if newSubscription.ServiceName != "" {
		changes["service_name"] = newSubscription.ServiceName
	}
	if newSubscription.Price != 0 {
		changes["price"] = newSubscription.Price
	}
	if newSubscription.EndDate != nil {
		changes["end_date"] = newSubscription.EndDate
	}

	updatedSubscriptionPostgresEntity, err := s.subscriptionRepo.UpdatePatch(ctx, id, changes)
	if err != nil {
		return nil, err
	}

	transferedUpdatedSubscription := transferPostgresEntityToServiceDomain(updatedSubscriptionPostgresEntity)

	return &transferedUpdatedSubscription, nil
}

func (s *SubscriptionService) DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error {
	if err := s.subscriptionRepo.DeleteByID(ctx, id); err != nil {
		return err
	}
	return nil
}