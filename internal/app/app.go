package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/http"
	"github.com/timac11/musthave-metrics-collector/internal/http/router"
	"github.com/timac11/musthave-metrics-collector/internal/logger"

	"github.com/timac11/musthave-metrics-collector/cmd/version"

	"golang.org/x/sync/errgroup"
)

func RunApplication() {
	conf := config.InitServerConfig()
	logger.Initialize("INFO")

	mux, err := router.InitRouter(conf)
	if err != nil {
		log.Fatal(err)
	}

	g, appCtx := errgroup.WithContext(context.Background())
	defer appCtx.Err()

	version.Print()
	logger.Info("Starting server on address: ", conf.Address)

	server := http.NewServer(conf.Address, mux)
	g.Go(func() error {
		return server.Start()
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
