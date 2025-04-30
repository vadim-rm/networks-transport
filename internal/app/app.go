package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	"transport/internal/config"
	"transport/internal/delivery/http/engine"
	"transport/internal/delivery/http/handlers"
	"transport/internal/delivery/kafka/consumer"
	kafkahandlers "transport/internal/delivery/kafka/handlers"
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

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          cfg.Kafka.Topic,
		Partition:      0,
		GroupID:        cfg.Kafka.GroupId,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})

	applicationRepository := repository.NewHTTPApplication(cfg.Services.ApplicationLevelBaseUrl)
	dataLinkRepository := repository.NewHTTPDataLink(cfg.Services.DataLinkLevelBaseUrl)

	transportService := service.NewTransportImpl(dataLinkRepository, applicationRepository)
	go transportService.Run(ctx)

	segmentsConsumer := consumer.NewConsumer(
		reader,
		kafkahandlers.NewDataLinkSegments(transportService).Receive,
	)
	go func() {
		err := segmentsConsumer.Run(ctx)
		if err != nil {
			log.Fatalf("error running segments consumer: %s", err.Error())
		}
	}()

	router := engine.Initialize()
	engine.InitializeExternalRoutes(
		router,
		handlers.NewMessage(transportService),
		handlers.NewSegment(writer),
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

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.MetricsPort), nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
