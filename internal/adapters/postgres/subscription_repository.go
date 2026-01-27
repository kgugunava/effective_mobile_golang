package postgres

import (
	"context"
	"fmt"
	"errors"
	"strings"

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

func (r *SubscriptionRepository) UpdatePut(ctx context.Context, sub SubscriptionEntity, id uuid.UUID) (SubscriptionEntity, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6
		WHERE subscription_id = $1
		RETURNING subscription_id, service_name, price, user_id, start_date, end_date`

	var updated SubscriptionEntity
	err := r.pool.QueryRow(ctx, query,
		id,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(
		&updated.SubscriptionID,
		&updated.ServiceName,
		&updated.Price,
		&updated.UserID,
		&updated.StartDate,
		&updated.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SubscriptionEntity{}, fmt.Errorf("subscription not found")
		}
		return SubscriptionEntity{}, fmt.Errorf("update failed: %w", err)
	}

	return updated, nil
}

func (r *SubscriptionRepository) UpdatePatch(ctx context.Context, id uuid.UUID, changes map[string]interface{}) (SubscriptionEntity, error) {
	if len(changes) == 0 {
		return r.GetByID(ctx, id)
	}

	allowedFields := map[string]bool{
		"service_name": true,
		"price":        true,
		"end_date":     true,
	}

	var setClauses []string
	var args []interface{}
	argIndex := 2 // $1 = id, $2+ = значения

	for field, value := range changes {
		if !allowedFields[field] {
			return SubscriptionEntity{}, fmt.Errorf("field '%s' cannot be updated", field)
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	query := fmt.Sprintf(`
		UPDATE subscriptions 
		SET %s 
		WHERE subscription_id = $1 
		RETURNING subscription_id, user_id, service_name, price, start_date, end_date`,
		strings.Join(setClauses, ", "),
	)

	args = append([]interface{}{id}, args...)

	var updated SubscriptionEntity
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&updated.SubscriptionID,
		&updated.UserID,
		&updated.ServiceName,
		&updated.Price,
		&updated.StartDate,
		&updated.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SubscriptionEntity{}, err
		}
		return SubscriptionEntity{}, fmt.Errorf("failed to patch subscription: %w", err)
	}

	return updated, nil
}

func (r *SubscriptionRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE subscription_id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no deleted rows")
	}

	return nil
}