package subscription

import "time"

type ListQuery struct {
	UserID      *UserID
	ServiceName *string
	Limit       int
	Offset      int
}
type TotalCostQuery struct {
	UserID      UserID
	ServiceName *string
	PeriodFrom  time.Time
	PeriodTo    time.Time
}
