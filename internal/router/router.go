package router

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/timac11/musthave-metrics-collector/internal/audit"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/handler"
	"github.com/timac11/musthave-metrics-collector/internal/handler/middleware"
	persistentstorage "github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	dbstorage "github.com/timac11/musthave-metrics-collector/internal/repository/db"
	memorystorage "github.com/timac11/musthave-metrics-collector/internal/repository/memory"
	"github.com/timac11/musthave-metrics-collector/internal/service"
)

// InitRouter return *chi.Mux for metrics server
func InitRouter(serverConfig *config.ServerConfig) (*chi.Mux, error) {
	// init service instance
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

	// init auditors instance
	auditor := audit.NewAuditor(serverConfig.AuditFile, serverConfig.AuditURL)
	handlers := handler.NewApplicationAPIContainer(*serviceInstance, *auditor)
	router := chi.NewRouter()

	router.Route("/debug", func(router chi.Router) {
		router.Mount("/", chimiddleware.Profiler())
	})

	router.Route("/", func(router chi.Router) {
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
		router.Get("/", handlers.GetMetricsPage)
	})

	return router, nil
}
