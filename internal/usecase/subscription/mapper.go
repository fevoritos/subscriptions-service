package subscription

import (
	"fmt"
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
