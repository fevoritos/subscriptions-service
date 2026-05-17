package transporthttp

import (
	"net/http"

	_ "subs-service/internal/transport/http/docs"
	httphandlers "subs-service/internal/transport/http/handlers"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(subHandler *httphandlers.SubscriptionHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /subs", subHandler.Create)
	router.HandleFunc("GET /subs/total", subHandler.TotalCost)
	router.HandleFunc("GET /subs/{id}", subHandler.GetByID)
	router.HandleFunc("PUT /subs/{id}", subHandler.Update)
	router.HandleFunc("DELETE /subs/{id}", subHandler.Delete)
	router.HandleFunc("GET /subs", subHandler.List)

	router.Handle("GET /docs/{any...}", httpSwagger.WrapHandler)

	return router
}
