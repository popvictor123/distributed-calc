package main

import (
	"log"
	"net/http"

	"github.com/popvictor123/distributed-calc/internal/orchestrator/api"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/calculator"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/repository"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo := repository.NewRepository()
	calc := calculator.NewCalculator()
	svc := service.NewService(repo, calc)
	handler := api.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/calculate", handler.CalculateHandler)
		r.Get("/expressions", handler.GetExpressionsHandler)
		r.Get("/expressions/{id}", handler.GetExpressionByIDHandler)
	})

	r.Route("/internal", func(r chi.Router) {
		r.Get("/task", handler.GetTaskHandler)
		r.Post("/task", handler.SubmitTaskResultHandler)
	})

	log.Println("Starting orchestrator server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
