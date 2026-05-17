package subscription

import (
	"context"
	subdomain "subs-service/internal/domain"
)

const monthYearLayout = "01-2006"

type Repository interface {
	Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	GetByID(ctx context.Context, id subdomain.SubscriptionID) (*subdomain.Subscription, error)
	Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	Delete(ctx context.Context, id subdomain.SubscriptionID) error
	List(ctx context.Context, q subdomain.ListQuery) ([]subdomain.Subscription, error)
	TotalCost(ctx context.Context, q subdomain.TotalCostQuery) (int, error)
}

type Usecase interface {
	Create(ctx context.Context, in CreateInput) (*subdomain.Subscription, error)
	GetByID(ctx context.Context, in GetByIDInput) (*subdomain.Subscription, error)
	Update(ctx context.Context, in UpdateInput) (*subdomain.Subscription, error)
	Delete(ctx context.Context, in DeleteInput) error
	List(ctx context.Context, filter ListInput) ([]subdomain.Subscription, error)
	TotalCost(ctx context.Context, in TotalCostInput) (int, error)
}

type CreateInput struct {
	UserID      string
	ServiceName string
	Price       int
	StartDate   string
	EndDate     *string
}

type UpdateInput struct {
	ID          string
	UserID      string
	ServiceName string
	Price       int
	StartDate   string
	EndDate     *string
}

type GetByIDInput struct {
	ID string
}

type DeleteInput struct {
	ID string
}

type ListFilter struct {
	UserID      *string
	ServiceName *string
	Limit       int
	Offset      int
}

type TotalCostFilter struct {
	UserID      string
	ServiceName *string
	PeriodFrom  string
	PeriodTo    string
}

type ListInput struct {
	UserID      *string
	ServiceName *string
	Limit       string
	Offset      string
}

type TotalCostInput struct {
	UserID      *string
	ServiceName *string
	PeriodFrom  string
	PeriodTo    string
}
