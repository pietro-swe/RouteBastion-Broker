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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"

	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/controllers"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/validators"
	"github.com/marechal-dev/RouteBastion-Broker/internal/utils"
)

type Server struct {
	Port int
	EncryptionKey []byte

	Trace trace.Tracer
	DB persistence.DBProvider

	HealthController    controllers.HealthController
	CustomersController controllers.CustomersController
	OptimizationsController controllers.OptimizationsController
}

func NewServer(config utils.AppEnvConfig, trace trace.Tracer) *http.Server {
	port, err := strconv.Atoi(config.ServerPort)
	if err != nil {
		log.Fatalf("could not cast server port: %v", err)
	}

	provider := persistence.NewPgxProvider(
		config.DBDatabase,
		config.DBPassword,
		config.DBUsername,
		config.DBPort,
		config.DBHost,
		config.DBSchema,
	)

	newServer := &Server{
		Port: port,
		EncryptionKey: config.EncryptionKeyBytes,
		Trace: trace,
		DB:   provider,
	}

	newServer.registerCustomValidators()
	newServer.registerControllers()

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
		v.RegisterValidation("cargoKind", validators.IsCargoKindValid)
	}
}

func (s *Server) registerControllers() {
	s.HealthController = controllers.NewHealthController(s.DB)
	s.CustomersController = controllers.NewCustomersController(s.EncryptionKey, s.Trace, s.DB)
	s.OptimizationsController = controllers.NewOptimizationsController(s.Trace)
}

func (s *Server) registerRoutes() http.Handler {
	router := gin.Default()
	router.Use(otelgin.Middleware("Broker-API"))

	// Health-check
	router.GET("/health", s.HealthController.Index)

	customers := router.Group("/customers")
	{
		customers.GET(
			"/:apiKey",
			s.CustomersController.GetOneByAPIKey,
		)
		customers.POST("/", s.CustomersController.Create)
	}

	optimizations := router.Group("/optimizations")
	{
		optimizations.GET("/sync", s.OptimizationsController.Optimize)
	}

	return router
}
