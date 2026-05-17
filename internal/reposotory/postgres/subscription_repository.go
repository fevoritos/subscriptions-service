package postgres

import (
	"context"
	"fmt"
	"log/slog"

	subdomain "subs-service/internal/domain"

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
	return &subdomain.Subscription{}, nil
}
func (r *Repository) Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	return &subdomain.Subscription{}, nil
}

func (r *Repository) Delete(ctx context.Context, id subdomain.SubscriptionID) error {
	return nil
}

func (r *Repository) List(ctx context.Context, q subdomain.ListQuery) ([]subdomain.Subscription, error) {
	return []subdomain.Subscription{}, nil
}
func (r *Repository) TotalCost(ctx context.Context, q subdomain.TotalCostQuery) (int, error) {
	return 0, nil
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
