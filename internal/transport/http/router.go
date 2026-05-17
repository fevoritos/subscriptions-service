package transporthttp

import (
	"net/http"
	httphandlers "subs-service/internal/transport/http/handlers"
)

func NewRouter(subHandler *httphandlers.SubscriptionHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /subs", subHandler.Create)
	router.HandleFunc("GET /subs/{id}", subHandler.GetByID)
	router.HandleFunc("PUT /subs/{id}", subHandler.Update)
	router.HandleFunc("DELETE /subs/{id}", subHandler.Delete)
	router.HandleFunc("GET /subs", subHandler.List)
	router.HandleFunc("GET /subs/total", subHandler.TotalCost)

	return router
}
