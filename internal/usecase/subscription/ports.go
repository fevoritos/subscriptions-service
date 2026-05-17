package subscription

import (
	"context"
	subdomain "subs-service/internal/domain"
	"time"
)

type Repository interface {
	Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	GetByID(ctx context.Context, id subdomain.SubscriptionID) (*subdomain.Subscription, error)
	Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	Delete(ctx context.Context, id subdomain.SubscriptionID) error
	List(ctx context.Context, q subdomain.ListQuery) ([]subdomain.Subscription, error)
	TotalCost(ctx context.Context, q subdomain.TotalCostQuery) (int, error)
}

type ListFilter struct {
	UserID      *subdomain.UserID
	ServiceName *string
	Limit       int
	Offset      int
}

type TotalCostFilter struct {
	UserID      subdomain.UserID
	ServiceName *string
	PeriodFrom  time.Time
	PeriodTo    time.Time
}

type Usecase interface {
	Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	GetByID(ctx context.Context, id subdomain.SubscriptionID) (*subdomain.Subscription, error)
	Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	Delete(ctx context.Context, id subdomain.SubscriptionID) error
	List(ctx context.Context, filter ListFilter) ([]subdomain.Subscription, error)
	TotalCost(ctx context.Context, filter TotalCostFilter) (int, error)
}
