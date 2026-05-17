package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	subdomain "subs-service/internal/domain"
)

type SubscriptionService struct {
	log  *slog.Logger
	repo Repository
}

func NewSubscriptionService(log *slog.Logger, repo Repository) *SubscriptionService {
	return &SubscriptionService{
		log:  log,
		repo: repo,
	}

}

func (s *SubscriptionService) Create(ctx context.Context, in CreateInput) (*subdomain.Subscription, error) {
	var op = "service.create"

	sub, err := toSubscriptionFromCreate(in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(err, ErrInvalidInput))
	}

	if err := validateSubscription(sub); err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(err, ErrInvalidInput))
	}

	created, err := s.repo.Create(ctx, &sub)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, in GetByIDInput) (*subdomain.Subscription, error) {
	const op = "service.get_by_id"

	id, err := parseSubscriptionID(in.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(err, ErrInvalidInput))
	}

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, subdomain.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, subdomain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (s *SubscriptionService) Update(ctx context.Context, in UpdateInput) (*subdomain.Subscription, error) {
	const op = "service.update"

	sub, err := toSubscriptionFromUpdate(in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(err, ErrInvalidInput))
	}

	if err := validateSubscription(sub); err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(err, ErrInvalidInput))
	}

	updated, err := s.repo.Update(ctx, &sub)
	if err != nil {
		if errors.Is(err, subdomain.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, subdomain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updated, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, in DeleteInput) error {
	const op = "service.delete"

	id, err := parseSubscriptionID(in.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionService) List(ctx context.Context, in ListInput) ([]subdomain.Subscription, error) {
	const op = "service.list"

	q, err := toListFilter(in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidInput)
	}
	return s.repo.List(ctx, q)
}

func (s *SubscriptionService) TotalCost(ctx context.Context, in TotalCostInput) (int, error) {
	const op = "service.total_cost"

	q, err := toTotalCostQuery(in)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidInput)
	}
	if q.PeriodFrom.After(q.PeriodTo) {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidInput)
	}
	return s.repo.TotalCost(ctx, q)
}

func validateSubscription(sub subdomain.Subscription) error {
	if sub.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if sub.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return fmt.Errorf("end_date must be before start_date")
	}
	return nil
}
