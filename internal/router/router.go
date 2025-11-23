package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/handler/middleware"
	"github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	"github.com/timac11/musthave-metrics-collector/internal/repository"
	"github.com/timac11/musthave-metrics-collector/internal/service"
	"net/http"
)

func InitRouter(serverConfig *config.ServerConfig) *chi.Mux {
	persistentStorage := persistentstorage.NewPersistentStorage(serverConfig.FileStoragePath)
	memStorage := repository.NewMemStorage(persistentStorage, serverConfig.Restore)
	service := service.NewService(memStorage)
	handlers := handler.NewApplicationAPIContainer(*service)

	m := middleware.NewMiddleware()
	middlewares := []func(http.Handler) http.Handler{
		m.GzipMiddleware,
		m.RequestLoggerMiddleware,
	}

	router := chi.NewRouter()
	// setup middlewares
	router.Use(middlewares...)

	router.Post("/update", handlers.UpdateMetricV2)
	router.Post("/update/", handlers.UpdateMetricV2)
	router.Post("/update/{metricType}/{metricName}/{value}", handlers.UpdateMetric)
	router.Post("/value", handlers.GetFullMetricInfo)
	router.Post("/value/", handlers.GetFullMetricInfo)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetric)

	router.Get(`/`, handlers.GetMetricsPage)

	return router
}
