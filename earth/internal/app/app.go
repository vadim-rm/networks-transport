package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"transport/internal/config"
	"transport/internal/delivery/http/engine"
	"transport/internal/delivery/http/handlers"
	"transport/internal/repository"
	"transport/internal/service"
)

func Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	dataLinkRepository := repository.NewHTTPDataLink(cfg.Services.DataLinkLevelBaseUrl)

	transportService := service.NewTransportImpl(dataLinkRepository)

	router := engine.Initialize()
	engine.InitializeExternalRoutes(
		router,
		handlers.NewMessage(transportService),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
