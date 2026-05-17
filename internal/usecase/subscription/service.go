package subscription

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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

func (s *SubscriptionService) Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	var op = "service.create"
	normalized, err := validateCreateInput(s.log, *sub)
	if err != nil {
		return &subdomain.Subscription{}, fmt.Errorf("%s: %w %w", op, err, ErrInvalidInput)
	}
	created, err := s.repo.Create(ctx, &normalized)
	if err != nil {
		return &subdomain.Subscription{}, err
	}

	return created, nil
}
func (s *SubscriptionService) GetByID(ctx context.Context, id subdomain.SubscriptionID) (*subdomain.Subscription, error) {
	return &subdomain.Subscription{}, nil
}
func (s *SubscriptionService) Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	return &subdomain.Subscription{}, nil
}
func (s *SubscriptionService) Delete(ctx context.Context, id subdomain.SubscriptionID) error {
	return nil
}
func (s *SubscriptionService) List(ctx context.Context, filter ListFilter) ([]subdomain.Subscription, error) {
	return []subdomain.Subscription{}, nil
}
func (s *SubscriptionService) TotalCost(ctx context.Context, filter TotalCostFilter) (int, error) {
	return 0, nil
}

func validateCreateInput(log *slog.Logger, input subdomain.Subscription) (subdomain.Subscription, error) {
	log.Info(input.ServiceName)
	input.ServiceName = strings.TrimSpace(input.ServiceName)
	log.Info(input.ServiceName)
	if input.ServiceName == "" {
		return subdomain.Subscription{}, fmt.Errorf("service_name is required")
	}

	return input, nil
}
