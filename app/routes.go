package app

import (
	"net/http"

	"github.com/fabiuhp/orders-api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRoutes)

	return router
}

func loadOrderRoutes(router chi.Router) {
	order := &handler.Order{}

	router.Post("/", order.Create)
	router.Get("/", order.List)
	router.Get("/{id}", order.GetById)
	router.Put("/{id}", order.UpdateById)
	router.Delete("/{id}", order.DeleteById)
}
