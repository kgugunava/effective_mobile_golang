package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriptionRepository(pool *pgxpool.Pool, logger *slog.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool:   pool,
		logger: logger,
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
		r.logger.Error("failed to insert subscription into DB",
			slog.String("subscription_id", subscription.SubscriptionID.String()),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to insert subscription: %w", err)
	}

	r.logger.Info("subscription created successfully",
		slog.String("subscription_id", subscription.SubscriptionID.String()),
	)
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
			r.logger.Warn("subscription not found",
				slog.String("subscription_id", id.String()),
			)
			return SubscriptionEntity{}, fmt.Errorf("subscription not found")
		}
		r.logger.Error("failed to get subscription",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		return SubscriptionEntity{}, fmt.Errorf("failed to get subscription: %w", err)
	}

	r.logger.Debug("subscription retrieved",
		slog.String("subscription_id", id.String()),
	)
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
			r.logger.Warn("subscription not found for update",
				slog.String("subscription_id", id.String()),
			)
			return SubscriptionEntity{}, fmt.Errorf("subscription not found")
		}
		r.logger.Error("update failed",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		return SubscriptionEntity{}, fmt.Errorf("update failed: %w", err)
	}

	r.logger.Info("subscription updated (PUT)",
		slog.String("subscription_id", id.String()),
	)
	return updated, nil
}

func (r *SubscriptionRepository) UpdatePatch(ctx context.Context, id uuid.UUID, changes map[string]interface{}) (SubscriptionEntity, error) {
	if len(changes) == 0 {
		r.logger.Debug("no fields to update, returning current state",
			slog.String("subscription_id", id.String()),
		)
		return r.GetByID(ctx, id)
	}

	allowedFields := map[string]bool{
		"service_name": true,
		"price":        true,
		"end_date":     true,
	}

	var setClauses []string
	var args []interface{}
	argIndex := 2

	for field, value := range changes {
		if !allowedFields[field] {
			r.logger.Error("attempt to update non-allowed field",
				slog.String("field", field),
				slog.String("subscription_id", id.String()),
			)
			return SubscriptionEntity{}, fmt.Errorf("field '%s' cannot be updated", field)
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
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
			r.logger.Warn("subscription not found for patch",
				slog.String("subscription_id", id.String()),
			)
			return SubscriptionEntity{}, fmt.Errorf("subscription not found")
		}
		r.logger.Error("failed to patch subscription",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
			slog.Any("changes", changes),
		)
		return SubscriptionEntity{}, fmt.Errorf("failed to patch subscription: %w", err)
	}

	r.logger.Info("subscription patched",
		slog.String("subscription_id", id.String()),
		slog.Any("changes", changes),
	)
	return updated, nil
}

func (r *SubscriptionRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE subscription_id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete subscription",
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn("delete requested but subscription not found",
			slog.String("subscription_id", id.String()),
		)
		return fmt.Errorf("subscription not found")
	}

	r.logger.Info("subscription deleted",
		slog.String("subscription_id", id.String()),
	)
	return nil
}

func (r *SubscriptionRepository) GetSubscriptionsList(ctx context.Context, serviceName string, userID uuid.UUID, startDate time.Time, endDate time.Time) ([]SubscriptionEntity, error) {
	query := `
		SELECT subscription_id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1`

	var args []interface{}
	argPos := 1

	query += fmt.Sprintf(" AND user_id = $%d", argPos)
	args = append(args, userID)
	argPos++

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, serviceName)
		argPos++
	}

	if !startDate.IsZero() {
		query += fmt.Sprintf(" AND start_date >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}

	if !endDate.IsZero() {
		query += fmt.Sprintf(" AND end_date <= $%d", argPos)
		args = append(args, endDate)
		argPos++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to execute list query",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to fetch subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []SubscriptionEntity
	for rows.Next() {
		var sub SubscriptionEntity
		err := rows.Scan(
			&sub.SubscriptionID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			r.logger.Error("failed to scan subscription row",
				slog.Any("error", err),
			)
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("row iteration error",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("subscription iteration failed: %w", err)
	}

	r.logger.Debug("subscriptions list fetched",
		slog.String("user_id", userID.String()),
		slog.Int("count", len(subscriptions)),
	)

	return subscriptions, nil
}