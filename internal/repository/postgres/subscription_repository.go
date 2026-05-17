package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	subdomain "subs-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func New(log *slog.Logger, pool *pgxpool.Pool) *Repository {
	return &Repository{
		log:  log,
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	const op = "repo.Create"

	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	createdSub, err := scanSubscription(row)

	if err != nil {
		r.log.Error("failed to execute insert query", "error", err)
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return createdSub, nil
}

func (r *Repository) GetByID(ctx context.Context, id subdomain.SubscriptionID) (*subdomain.Subscription, error) {
	const op = "repo.GetByID"

	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	sub, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, subdomain.ErrNotFound)
		}
		r.log.Error("failed to execute select query", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (r *Repository) Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	const op = "repo.Update"

	const query = `
		UPDATE subscriptions
		SET
			service_name = $1,
			price        = $2,
			user_id      = $3,
			start_date   = $4,
			end_date     = $5,
			updated_at   = NOW()
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	updated, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, subdomain.ErrNotFound)
		}
		r.log.Error(op, slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id subdomain.SubscriptionID) error {
	const op = "repo.Delete"

	const query = `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.log.Error(op, slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, q subdomain.ListQuery) ([]subdomain.Subscription, error) {
	const op = "repo.List"

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1
	`
	args := make([]any, 0, 4)
	n := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", n)
		args = append(args, *q.UserID)
		n++
	}

	if q.ServiceName != nil && strings.TrimSpace(*q.ServiceName) != "" {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", n)
		args = append(args, "%"+strings.TrimSpace(*q.ServiceName)+"%")
		n++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, q.Limit, q.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		r.log.Error(op, slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	subs := make([]subdomain.Subscription, 0)
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		subs = append(subs, *sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func (r *Repository) TotalCost(ctx context.Context, q subdomain.TotalCostQuery) (int, error) {
	const op = "repo.TotalCost"
	periodToExclusive := q.PeriodTo.AddDate(0, 1, 0)
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date < $1
		  AND (end_date IS NULL OR end_date >= $2)
	`
	args := []any{periodToExclusive, q.PeriodFrom}
	n := 3
	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", n)
		args = append(args, *q.UserID)
		n++
	}
	if q.ServiceName != nil && strings.TrimSpace(*q.ServiceName) != "" {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", n)
		args = append(args, "%"+strings.TrimSpace(*q.ServiceName)+"%")
	}
	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		r.log.Error(op, slog.Any("err", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return total, nil
}

type subScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner subScanner) (*subdomain.Subscription, error) {
	var (
		sub subdomain.Subscription
	)

	if err := scanner.Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	); err != nil {
		return nil, err
	}

	return &sub, nil
}
