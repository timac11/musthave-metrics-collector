package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/handler/middleware"
	"github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	"github.com/timac11/musthave-metrics-collector/internal/repository/db"
	"github.com/timac11/musthave-metrics-collector/internal/repository/memory"
	"github.com/timac11/musthave-metrics-collector/internal/service"
	"net/http"
)

func InitRouter(serverConfig *config.ServerConfig) (*chi.Mux, error) {
	var serviceInstance *service.Service

	if serverConfig.DatabaseDsn != "" {
		dbClient, err := dbstorage.NewPgClient(serverConfig.DatabaseDsn)

		if err != nil {
			return nil, err
		}

		serviceInstance = service.NewService(dbClient)
	} else {
		persistentStorage := persistentstorage.NewPersistentStorage(serverConfig.FileStoragePath)
		memStorage := memorystorage.NewMemStorage(persistentStorage, serverConfig.Restore)
		serviceInstance = service.NewService(memStorage)
	}

	handlers := handler.NewApplicationAPIContainer(*serviceInstance)

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
