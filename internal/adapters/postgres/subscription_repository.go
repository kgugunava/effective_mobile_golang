package postgres

import (
	"context"
	"fmt"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool: pool,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription SubscriptionEntity) error {
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

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (SubscriptionEntity, error) {
	query := `
		SELECT subscription_id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE subscription_id = $1
	`

	var entity SubscriptionEntity
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&entity.SubscriptionID,
		&entity.ServiceName,
		&entity.Price,
		&entity.UserID,
		&entity.StartDate,
		&entity.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SubscriptionEntity{}, fmt.Errorf("subscription not found")
		}
		return SubscriptionEntity{}, fmt.Errorf("failed to get subscription: %w", err)
	}

	return entity, nil
}