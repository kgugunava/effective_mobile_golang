package service

import (
	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"
	"github.com/kgugunava/effective_mobile_golang/internal/domain"
)

func transferServiceDomainToPostgresEntity(subscription domain.Subscription) postgres.SubscriptionEntity {
	entity := postgres.SubscriptionEntity{
		SubscriptionID: subscription.SubscriptionID,
		UserID:      subscription.UserID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		StartDate:   subscription.StartDate,
	}

	if subscription.EndDate != nil {
		entity.EndDate = subscription.EndDate
	}

	return entity

}

func transferPostgresEntityToServiceDomain(entity postgres.SubscriptionEntity) domain.Subscription {
	domain := domain.Subscription{
		SubscriptionID: entity.SubscriptionID,
		ServiceName: entity.ServiceName,
		Price: entity.Price,
		UserID: entity.UserID,
		StartDate: entity.StartDate,
	}

	if entity.EndDate != nil {
		domain.EndDate = entity.EndDate
	}

	return domain
}

func transferPostgresEntityListsToServiceDomainList(entities []postgres.SubscriptionEntity) []domain.Subscription {
	var domainSubscriptions []domain.Subscription

	for _, entity := range(entities) {
		domain := domain.Subscription{
			SubscriptionID: entity.SubscriptionID,
			ServiceName: entity.ServiceName,
			Price: entity.Price,
			UserID: entity.UserID,
			StartDate: entity.StartDate,
		}
		domainSubscriptions = append(domainSubscriptions, domain)
	}

	return domainSubscriptions
}