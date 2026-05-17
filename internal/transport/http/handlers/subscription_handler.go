package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	subdomain "subs-service/internal/domain"
	subusecase "subs-service/internal/usecase/subscription"
	"time"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidUserID       = errors.New("invalid user id")
)

type SubscriptionHandler struct {
	log     *slog.Logger
	usecase subusecase.Usecase
}

func NewSubscriptionHandler(log *slog.Logger, usecase subusecase.Usecase) *SubscriptionHandler {
	return &SubscriptionHandler{log: log, usecase: usecase}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handler.create"
	var req CreateSubscriptionRequest
	if err := decodeJSON(r, &req); err != nil {
		h.log.Error(op, "error", err)
		writeUsecaseError(w, h.log, subusecase.ErrInvalidInput)
		return
	}

	userID, err := subdomain.ParseID(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadGateway, ErrInvalidUserID)
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		h.log.Error(op, "error", err)
		writeError(w, http.StatusBadRequest, errors.New("invalid start_date format, expected MM-YYYY"))
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			h.log.Error(op, "error", err)
			writeError(w, http.StatusBadRequest, errors.New("invalid end_date format, expected MM-YYYY"))
			return
		}
		endDate = &parsedEndDate
	}
	sub := &subdomain.Subscription{
		UserID:      userID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	created, err := h.usecase.Create(r.Context(), sub)
	if err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}

	response := toSubscriptionResponse(created)
	writeJSON(w, http.StatusCreated, response)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeUsecaseError(w http.ResponseWriter, log *slog.Logger, err error) {
	switch {
	case errors.Is(err, subusecase.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, subusecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		log.Error("usecase error", slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, ErrInternalServerError)
	}
}
