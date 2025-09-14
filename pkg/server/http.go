/*
Package server provides an unified interface that represents the application Server
*/
package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
	"github.com/pietro-swe/RouteBastion-Broker/internal/modules/customer"
	"github.com/pietro-swe/RouteBastion-Broker/internal/modules/health"
	"github.com/pietro-swe/RouteBastion-Broker/internal/modules/optimization"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/env"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/httputils"

	// "github.com/pietro-swe/RouteBastion-Broker/pkg/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	Port int

	EncryptionKey []byte
	Trace         trace.Tracer
	DB            db.DBProvider
}

func NewHTTPServer(config env.AppEnvConfig, trace trace.Tracer) *http.Server {
	port, err := strconv.Atoi(config.ServerPort)
	if err != nil {
		log.Fatalf("could not cast server port: %v", err)
	}

	provider := db.NewPgxProvider(
		config.DBDatabase,
		config.DBPassword,
		config.DBUsername,
		config.DBPort,
		config.DBHost,
		config.DBSchema,
	)

	newServer := &Server{
		Port:          port,
		EncryptionKey: config.EncryptionKeyBytes,
		Trace:         trace,
		DB:            provider,
	}

	newServer.registerCustomValidators()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.Port),
		Handler:      newServer.registerRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("cargoKind", httputils.IsCargoKindValid)
	}
}

func (s *Server) registerRoutes() http.Handler {
	router := gin.Default()

	// Middlewares
	router.Use(otelgin.Middleware("Broker-REST-API"))

	// Health-check
	router.GET("/health", health.MakeHealthCheckHandler(s.DB))

	v1 := router.Group("/v1")
	{
		customers := v1.Group("/customers")
		{
			customers.GET(
				"/:apiKey",
				customer.MakeGetOneByAPIKeyHandler(s.Trace, s.DB),
			)

			customers.POST("/", customer.MakeCreateCustomerHandler(
				s.EncryptionKey,
				s.Trace,
				s.DB,
			))
		}

		optimizations := v1.Group("/optimizations")
		// optimizations.Use(middlewares.WithValidAPIKey(s.DB))
		{
			optimizations.GET("/sync", optimization.MakeOptimizeSyncHandler(s.Trace))
		}
	}

	return router
}
