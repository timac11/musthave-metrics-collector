package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/handler/middleware"
	"github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	"github.com/timac11/musthave-metrics-collector/internal/repository/db"
	"github.com/timac11/musthave-metrics-collector/internal/repository/memory"
	"github.com/timac11/musthave-metrics-collector/internal/service"
)

func InitRouter(serverConfig *config.ServerConfig) (*chi.Mux, error) {
	var serviceInstance *service.Service
	serviceConfig := service.ServiceConfig{Attempts: serverConfig.RetryAttempts, AttemptsInterval: serverConfig.RetryInterval}

	if serverConfig.DatabaseDsn != "" {
		dbClient, err := dbstorage.NewPgClient(serverConfig.DatabaseDsn)

		if err != nil {
			return nil, err
		}

		serviceInstance = service.NewService(dbClient, serviceConfig)
	} else {
		persistentStorage := persistentstorage.NewPersistentStorage(serverConfig.FileStoragePath)
		memStorage := memorystorage.NewMemStorage(persistentStorage, serverConfig.Restore)
		serviceInstance = service.NewService(memStorage, serviceConfig)
	}

	handlers := handler.NewApplicationAPIContainer(*serviceInstance)
	router := chi.NewRouter()
	m := middleware.NewMiddleware(serverConfig.SigningKey)

	// setup middlewares
	middlewares := []func(http.Handler) http.Handler{
		m.GzipMiddleware,
		m.RequestLoggerMiddleware,
	}
	router.Use(middlewares...)

	if serverConfig.SigningKey != "" {
		hashMiddleware := []func(http.Handler) http.Handler{
			m.CheckSignatureMiddleware,
			m.SetSignatureMiddleware,
		}
		router.Use(hashMiddleware...)
	}

	// setup routes
	router.Post("/update", handlers.UpdateMetricV2)
	router.Post("/update/", handlers.UpdateMetricV2)
	router.Post("/updates", handlers.UpdateMetrics)
	router.Post("/updates/", handlers.UpdateMetrics)
	router.Post("/update/{metricType}/{metricName}/{value}", handlers.UpdateMetric)
	router.Post("/value", handlers.GetFullMetricInfo)
	router.Post("/value/", handlers.GetFullMetricInfo)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	router.Get("/ping", handlers.DBPing)

	router.Get(`/`, handlers.GetMetricsPage)

	return router, nil
}
