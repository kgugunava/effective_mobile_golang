package postgres

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionEntity struct {
	SubscriptionID uuid.UUID `db:"subscription_id"`
	ServiceName string `db:"service_name"`
	Price int `db:"price"`
	UserID uuid.UUID `db:"user_id"`
	StartDate time.Time `db:"start_date"`
	EndDate *time.Time `db:"end_date"`
}