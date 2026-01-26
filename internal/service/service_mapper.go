package service

import (
	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"
	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

func transferServiceDomainToPostgresEntity(subscription domain.Subscription) postgres.Subscription {
	entity := postgres.Subscription{
		SubscriptionID:          subscription.SubscriptionID,
		UserID:      subscription.UserID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		StartDate:   subscription.StartDate,
	}

	if subscription.EndDate != nil {
		entity.EndDate = *subscription.EndDate
	}

	return entity

}