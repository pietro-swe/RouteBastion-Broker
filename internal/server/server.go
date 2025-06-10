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

	infraDB "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/database"
	customers "github.com/marechal-dev/RouteBastion-Broker/internal/modules/customers/infrastructure/http/controllers"
	health "github.com/marechal-dev/RouteBastion-Broker/internal/modules/health/infrastructure/http/controllers"
	"github.com/marechal-dev/RouteBastion-Broker/internal/modules/shared/application/validators"
	"github.com/marechal-dev/RouteBastion-Broker/internal/utils"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	port int

	tracer *trace.TracerProvider
	db infraDB.DBProvider

	healthController    health.HealthController
	customersController customers.CustomersController
}

func NewServer(config utils.AppEnvConfig, tracer *trace.TracerProvider) *http.Server {
	port, err := strconv.Atoi(config.ServerPort)
	if err != nil {
		log.Fatalf("could not cast server port: %v", err)
	}

	provider := infraDB.NewPgxProvider(
		config.DBDatabase,
		config.DBPassword,
		config.DBUsername,
		config.DBPort,
		config.DBHost,
		config.DBSchema,
	)

	newServer := &Server{
		port: port,
		tracer: tracer,
		db:   provider,
	}

	newServer.RegisterCustomValidators()
	newServer.RegisterControllers()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("cargoKind", validators.IsValidCargoKind)
	}
}

func (s *Server) RegisterControllers() {
	s.healthController = health.NewHealthController(s.db)
	s.customersController = customers.NewCustomersController(s.db)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Health-check
	r.GET("/health", s.healthController.Index)

	// Customers
	customers := r.Group("/customers")
	{
		customers.GET(
			"/:apiKey",
			s.customersController.GetOneByApiKey,
		)
		customers.POST("/", s.customersController.Create)
	}

	return r
}
