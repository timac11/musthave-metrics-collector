package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/repository"
	"github.com/timac11/musthave-metrics-collector/internal/service"
)

func InitRouter() *chi.Mux {
	service := service.NewService(repository.NewMemStorage())
	handlers := handler.NewApplicationAPIContainer(*service)

	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{value}", handlers.UpdateMetric)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	router.Get(`/`, handlers.GetMetricsPage)

	return router
}
