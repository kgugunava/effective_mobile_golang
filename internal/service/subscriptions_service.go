package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

type SubscriptionService struct {
	subscriptionRepo SubscriptionRepository
	logger           *slog.Logger
}

func NewSubscriptionService(repo SubscriptionRepository, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: repo,
		logger:           logger,
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
	s.logger.Debug("creating new subscription")

	subscription.SubscriptionID = uuid.New()

	if !isSubscriptionValid(subscription) {
		s.logger.Warn("invalid subscription data",
			slog.String("subscription_id", subscription.SubscriptionID.String()),
		)
	}

	if err := s.subscriptionRepo.Create(ctx, transferServiceDomainToPostgresEntity(*subscription)); err != nil {
		s.logger.Error("failed to create subscription in repository",
			slog.String("subscription_id", subscription.SubscriptionID.String()),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("repository create failed: %w", err)
	}

	s.logger.Info("subscription created successfully",
		slog.String("subscription_id", subscription.SubscriptionID.String()),
	)
	return subscription, nil
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	s.logger.Debug("getting subscription by ID",
		slog.String("subscription_id", id.String()),
	)

	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription from repository",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		return domain.Subscription{}, err
	}

	s.logger.Debug("subscription retrieved successfully",
		slog.String("subscription_id", id.String()),
	)
	return transferPostgresEntityToServiceDomain(subscription), nil
}

func (s *SubscriptionService) UpdateSubscriptionPut(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error) {
	s.logger.Debug("updating subscription with PUT",
		slog.String("subscription_id", id.String()),
	)

	var updatedSubscription domain.Subscription

	if isSubscriptionValid(newSubscription) {
		updatedSubscriptionPostgresEntity, err := s.subscriptionRepo.UpdatePut(ctx, transferServiceDomainToPostgresEntity(*newSubscription), id)
		if err != nil {
			s.logger.Error("failed to update subscription (PUT) in repository",
				slog.String("subscription_id", id.String()),
				slog.Any("error", err),
			)
			return nil, err
		}
		updatedSubscription = transferPostgresEntityToServiceDomain(updatedSubscriptionPostgresEntity)
	} else {
		s.logger.Warn("invalid subscription data in PUT update",
			slog.String("subscription_id", id.String()),
		)
	}

	s.logger.Info("subscription updated (PUT) successfully",
		slog.String("subscription_id", id.String()),
	)
	return &updatedSubscription, nil
}

func (s *SubscriptionService) UpdateSubscriptionPatch(ctx context.Context, id uuid.UUID, newSubscription *domain.Subscription) (*domain.Subscription, error) {
	s.logger.Debug("updating subscription with PATCH",
		slog.String("subscription_id", id.String()),
	)

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

	if len(changes) == 0 {
		s.logger.Warn("PATCH request with no changes",
			slog.String("subscription_id", id.String()),
		)
	}

	updatedSubscriptionPostgresEntity, err := s.subscriptionRepo.UpdatePatch(ctx, id, changes)
	if err != nil {
		s.logger.Error("failed to patch subscription in repository",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
			slog.Any("changes", changes),
		)
		return nil, err
	}

	transferedUpdatedSubscription := transferPostgresEntityToServiceDomain(updatedSubscriptionPostgresEntity)

	s.logger.Info("subscription patched successfully",
		slog.String("subscription_id", id.String()),
		slog.Any("changes", changes),
	)
	return &transferedUpdatedSubscription, nil
}

func (s *SubscriptionService) DeleteSubscriptionByID(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("deleting subscription",
		slog.String("subscription_id", id.String()),
	)

	if err := s.subscriptionRepo.DeleteByID(ctx, id); err != nil {
		s.logger.Error("failed to delete subscription in repository",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("subscription deleted successfully",
		slog.String("subscription_id", id.String()),
	)
	return nil
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, serviceName string, userID uuid.UUID, startDate time.Time, endDate time.Time) ([]domain.Subscription, error) {
	postgresEntities, err := s.subscriptionRepo.GetSubscriptionsList(ctx, serviceName, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("failed to get subscriptions list in repository",
			slog.String("service_name", serviceName),
			slog.String("user_id", userID.String()),
			slog.String("start_date", startDate.String()),
			slog.String("end_date", endDate.String()),
			slog.Any("error", err),
		)
		return []domain.Subscription{}, err
	}

	domainSubscriptionsList := transferPostgresEntityListsToServiceDomainList(postgresEntities)

	return domainSubscriptionsList, nil
}