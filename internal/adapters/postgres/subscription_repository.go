package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool: pool,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription Subscription) error {
	query := `
		INSERT INTO subscriptions (subscription_id, user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		subscription.SubscriptionID,
		subscription.UserID,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate, 
	)
	if err != nil {
		return fmt.Errorf("failed to insert subscription: %w", err)
	}

	return nil
}