package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/timac11/musthave-metrics-collector/internal/audit"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/grpc"
	"github.com/timac11/musthave-metrics-collector/internal/http"
	"github.com/timac11/musthave-metrics-collector/internal/http/router"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	persistentstorage "github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	dbstorage "github.com/timac11/musthave-metrics-collector/internal/repository/db"
	memorystorage "github.com/timac11/musthave-metrics-collector/internal/repository/memory"
	"github.com/timac11/musthave-metrics-collector/internal/service"

	"github.com/timac11/musthave-metrics-collector/cmd/version"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

func RunApplication() {
	conf := config.InitServerConfig()
	logger.Initialize("INFO")

	auditor := initAppAuditor(conf)
	appService, err := initAppService(conf)
	if err != nil {
		log.Fatal(err)
	}

	g, appCtx := errgroup.WithContext(context.Background())
	defer appCtx.Err()

	version.Print()
	logger.Info("Starting server on address: ", conf.Address)

	var server Server
	g.Go(func() error {
		if conf.Mode == "http" {
			mux, err := router.InitRouter(conf, appService, auditor)
			if err != nil {
				return err
			}
			server := http.NewServer(conf.Address, mux)
			return server.Start()
		}

		if conf.Mode == "grpc" {
			server := grpc.NewServer(conf.Address, *appService)
			return server.Start()
		}

		return nil
	})

	// graceful shutdown
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	g.Go(func() error {
		<-quitChan
		logger.Info("graceful shutdown signal")

		if err = server.Stop(appCtx); err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err.Error())
	}
}

func initAppService(conf *config.ServerConfig) (*service.Service, error) {
	// init service instance
	var serviceInstance *service.Service
	serviceConfig := service.ServiceConfig{Attempts: conf.RetryAttempts, AttemptsInterval: conf.RetryInterval}

	if conf.DatabaseDsn != "" {
		dbClient, err := dbstorage.NewPgClient(conf.DatabaseDsn)

		if err != nil {
			return nil, err
		}

		serviceInstance = service.NewService(dbClient, serviceConfig)
	} else {
		persistentStorage := persistentstorage.NewPersistentStorage(conf.FileStoragePath)
		memStorage := memorystorage.NewMemStorage(persistentStorage, conf.Restore)
		serviceInstance = service.NewService(memStorage, serviceConfig)
	}

	return serviceInstance, nil
}

func initAppAuditor(conf *config.ServerConfig) *audit.Auditor {
	return audit.NewAuditor(conf.AuditFile, conf.AuditURL)
}
