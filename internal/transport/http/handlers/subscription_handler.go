package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	subdomain "subs-service/internal/domain"
	subusecase "subs-service/internal/usecase/subscription"
)

var (
	ErrBadRequest          = errors.New("invalid request")
	ErrNotFound            = errors.New("resource not found")
	ErrInternalServerError = errors.New("internal server error")
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

	sub := subusecase.CreateInput{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
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
	const op = "handler.get_by_id"
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	sub, err := h.usecase.GetByID(r.Context(), subusecase.GetByIDInput{ID: id})
	if err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toSubscriptionResponse(sub))
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handler.update"

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	var req UpdateSubscriptionRequest
	if err := decodeJSON(r, &req); err != nil {
		h.log.Error(op, "error", err)
		writeUsecaseError(w, h.log, subusecase.ErrInvalidInput)
		return
	}

	updated, err := h.usecase.Update(r.Context(), subusecase.UpdateInput{
		ID:          id,
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusOK, toSubscriptionResponse(updated))
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.delete"

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	if err := h.usecase.Delete(r.Context(), subusecase.DeleteInput{ID: id}); err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	subs, err := h.usecase.List(r.Context(), subusecase.ListInput{
		UserID:      optionalQuery(q.Get("user_id")),
		ServiceName: optionalQuery(q.Get("service_name")),
		Limit:       q.Get("limit"),
		Offset:      q.Get("offset"),
	})
	if err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}
	writeJSON(w, http.StatusOK, toSubscriptionsResponse(subs))
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	total, err := h.usecase.TotalCost(r.Context(), subusecase.TotalCostInput{
		UserID:      optionalQuery(q.Get("user_id")),
		ServiceName: optionalQuery(q.Get("service_name")),
		PeriodFrom:  q.Get("period_from"),
		PeriodTo:    q.Get("period_to"),
	})
	if err != nil {
		writeUsecaseError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusOK, TotalCostResponse{Total: total})
}

func optionalQuery(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
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
	case errors.Is(err, subdomain.ErrNotFound):
		log.Info("not found", slog.Any("err", err))
		writeError(w, http.StatusNotFound, ErrNotFound)
	case errors.Is(err, subusecase.ErrInvalidInput):
		log.Warn("invalid input", slog.Any("err", err))
		writeError(w, http.StatusBadRequest, ErrBadRequest)
	default:
		log.Error("usecase error", slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, ErrInternalServerError)
	}
}
