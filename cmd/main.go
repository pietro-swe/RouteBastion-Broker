package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/marechal-dev/RouteBastion-Broker/internal/server"
	"github.com/marechal-dev/RouteBastion-Broker/internal/server/instrumentation"
	"github.com/marechal-dev/RouteBastion-Broker/internal/utils"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("[Broker] shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("[Broker] server forced to shutdown with error: %v", err)
	}

	log.Println("[Broker] server exiting")

	done <- true
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("[Broker] could not load config: %v", err)
	}

	encryptionKeyBytes, err := base64.StdEncoding.DecodeString(config.EncryptionKey)
	if err != nil {
		log.Fatalf("[Broker] failed to decode encryption key: %v", err)
	}

	config.EncryptionKeyBytes = encryptionKeyBytes

	exporter, err := instrumentation.InitExporter(config)
	if err != nil {
		log.Fatalf("[Broker] failed to initialize exporter: %v", err)
	}

	tracer := instrumentation.InitTracer(exporter)

	defer func() {
		exporter.Shutdown(context.Background())
		tracer.Shutdown(context.Background())
	}()

	server := server.NewServer(config, tracer)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	serverErr := server.ListenAndServe()
	if serverErr != nil && serverErr != http.ErrServerClosed {
		log.Fatalf("[Broker] http server error: %s", serverErr)
	}

	<-done
	log.Println("[Broker] graceful shutdown complete.")
}
