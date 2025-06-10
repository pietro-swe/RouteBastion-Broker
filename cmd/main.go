package main

import (
	"context"
	"fmt"
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
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown with error: %v", err)
	}

	log.Println("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	exporter, err := instrumentation.InitExporter(config)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	defer exporter.Shutdown(context.Background())

	tracer := instrumentation.InitTracer(exporter)
	defer tracer.Shutdown(context.Background())

	server := server.NewServer(config, tracer)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	serverErr := server.ListenAndServe()
	if serverErr != nil && serverErr != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", serverErr))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
