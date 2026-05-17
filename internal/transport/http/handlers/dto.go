package handlers

import subdomain "subs-service/internal/domain"

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func toSubscriptionResponse(sub *subdomain.Subscription) SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          sub.ID.String(),
		UserID:      sub.UserID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		StartDate:   sub.StartDate.Format("01-2006"),
	}

	if sub.EndDate != nil {
		formattedEnd := sub.EndDate.Format("01-2006")
		resp.EndDate = &formattedEnd
	}

	return resp
}
