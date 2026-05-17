package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	subdomain "subs-service/internal/domain"
	subusecase "subs-service/internal/usecase/subscription"
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
	w.WriteHeader(http.StatusOK)
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
	case errors.Is(err, subdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, subusecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		log.Error("usecase error", slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, ErrInternalServerError)
	}
}
