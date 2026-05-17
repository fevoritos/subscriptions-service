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
	return &subdomain.Subscription{}, nil
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
	return nil
}

func (s *SubscriptionService) List(ctx context.Context, filter ListFilter) ([]subdomain.Subscription, error) {
	return []subdomain.Subscription{}, nil
}

func (s *SubscriptionService) TotalCost(ctx context.Context, filter TotalCostFilter) (int, error) {
	return 0, nil
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
