package subscription

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	subdomain "subs-service/internal/domain"
)

func toSubscriptionFromCreate(in CreateInput) (subdomain.Subscription, error) {
	userID, err := subdomain.ParseID(strings.TrimSpace(in.UserID))
	if err != nil {
		return subdomain.Subscription{}, fmt.Errorf("user_id: %w", err)
	}
	startDate, err := parseMonthYear(in.StartDate)
	if err != nil {
		return subdomain.Subscription{}, fmt.Errorf("start_date: %w", err)
	}
	var endDate *time.Time
	if in.EndDate != nil {
		parsed, err := parseMonthYear(*in.EndDate)
		if err != nil {
			return subdomain.Subscription{}, fmt.Errorf("end_date: %w", err)
		}
		endDate = &parsed
	}
	return subdomain.Subscription{
		UserID:      userID,
		ServiceName: strings.TrimSpace(in.ServiceName),
		Price:       in.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func toSubscriptionFromUpdate(in UpdateInput) (subdomain.Subscription, error) {
	id, err := subdomain.ParseSubscriptionID(strings.TrimSpace(in.ID))
	if err != nil {
		return subdomain.Subscription{}, fmt.Errorf("id: %w", err)
	}
	sub, err := toSubscriptionFromCreate(CreateInput{
		UserID:      in.UserID,
		ServiceName: in.ServiceName,
		Price:       in.Price,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if err != nil {
		return subdomain.Subscription{}, err
	}
	sub.ID = id
	return sub, nil
}

func toListQuery(f ListFilter) (subdomain.ListQuery, error) {
	q := subdomain.ListQuery{
		ServiceName: f.ServiceName,
		Limit:       f.Limit,
		Offset:      f.Offset,
	}
	if f.UserID != nil {
		uid, err := subdomain.ParseID(strings.TrimSpace(*f.UserID))
		if err != nil {
			return subdomain.ListQuery{}, fmt.Errorf("user_id: %w", err)
		}
		q.UserID = &uid
	}
	return q, nil
}

func parseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse(monthYearLayout, strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("expected MM-YYYY: %w", err)
	}
	return t, nil
}

func parseSubscriptionID(id string) (subdomain.SubscriptionID, error) {
	parsed, err := subdomain.ParseSubscriptionID(strings.TrimSpace(id))
	if err != nil {
		return subdomain.SubscriptionID{}, fmt.Errorf("id: %w", err)
	}
	return parsed, nil
}

func toListFilter(in ListInput) (subdomain.ListQuery, error) {
	q := subdomain.ListQuery{}

	if in.UserID != nil {
		uid, err := subdomain.ParseID(strings.TrimSpace(*in.UserID))
		if err != nil {
			return subdomain.ListQuery{}, fmt.Errorf("user_id: %w", err)
		}
		q.UserID = &uid
	}

	if in.ServiceName != nil {
		name := strings.TrimSpace(*in.ServiceName)
		if name != "" {
			q.ServiceName = &name
		}
	}

	limit, err := parseNonNegativeInt(in.Limit, defaultListLimit)
	if err != nil {
		return subdomain.ListQuery{}, fmt.Errorf("limit: %w", err)
	}
	offset, err := parseNonNegativeInt(in.Offset, 0)
	if err != nil {
		return subdomain.ListQuery{}, fmt.Errorf("offset: %w", err)
	}

	q.Limit = clampLimit(limit)
	q.Offset = offset

	return q, nil
}

const defaultListLimit = 20
const maxListLimit = 100

func parseNonNegativeInt(raw string, defaultVal int) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("must be a non-negative integer")
	}
	return n, nil
}

func clampLimit(n int) int {
	if n <= 0 {
		return defaultListLimit
	}
	if n > maxListLimit {
		return maxListLimit
	}
	return n
}

func toTotalCostQuery(in TotalCostInput) (subdomain.TotalCostQuery, error) {
	if strings.TrimSpace(in.PeriodFrom) == "" {
		return subdomain.TotalCostQuery{}, fmt.Errorf("period_from is required")
	}
	if strings.TrimSpace(in.PeriodTo) == "" {
		return subdomain.TotalCostQuery{}, fmt.Errorf("period_to is required")
	}

	from, err := parseMonthYear(in.PeriodFrom)
	if err != nil {
		return subdomain.TotalCostQuery{}, fmt.Errorf("period_from: %w", err)
	}
	to, err := parseMonthYear(in.PeriodTo)
	if err != nil {
		return subdomain.TotalCostQuery{}, fmt.Errorf("period_to: %w", err)
	}

	q := subdomain.TotalCostQuery{
		PeriodFrom: from,
		PeriodTo:   to,
	}

	if in.UserID != nil {
		uid, err := subdomain.ParseID(strings.TrimSpace(*in.UserID))
		if err != nil {
			return subdomain.TotalCostQuery{}, fmt.Errorf("user_id: %w", err)
		}
		q.UserID = &uid
	}

	if in.ServiceName != nil {
		name := strings.TrimSpace(*in.ServiceName)
		if name != "" {
			q.ServiceName = &name
		}
	}

	return q, nil
}
