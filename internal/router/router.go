package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/handler/middleware"
	"github.com/timac11/musthave-metrics-collector/internal/repository"
	"github.com/timac11/musthave-metrics-collector/internal/service"
	"net/http"
)

func InitRouter() *chi.Mux {
	service := service.NewService(repository.NewMemStorage())
	handlers := handler.NewApplicationAPIContainer(*service)
	m := middleware.NewMiddleware()
	middlewares := []func(http.Handler) http.Handler{
		m.RequestLoggerMiddleware,
	}

	router := chi.NewRouter()
	// setup middlewares
	router.Use(middlewares...)

	// setup routes
	router.Post("/update/{metricType}/{metricName}/{value}", handlers.UpdateMetric)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	router.Get(`/`, handlers.GetMetricsPage)

	return router
}
