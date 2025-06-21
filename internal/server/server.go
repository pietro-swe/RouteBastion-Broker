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

	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"

	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/controllers"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/validators"
	"github.com/marechal-dev/RouteBastion-Broker/internal/utils"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	Port int

	Tracer *trace.TracerProvider

	DB persistence.DBProvider

	EncryptionKey []byte

	HealthController    controllers.HealthController
	CustomersController controllers.CustomersController
}

func NewServer(config utils.AppEnvConfig, tracer *trace.TracerProvider) *http.Server {
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
		Tracer: tracer,
		DB:   provider,
		EncryptionKey: config.EncryptionKeyBytes,
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
	s.CustomersController = controllers.NewCustomersController(s.EncryptionKey, s.DB)
}

func (s *Server) registerRoutes() http.Handler {
	router := gin.Default()

	// Health-check
	router.GET("/health", s.HealthController.Index)

	customers := router.Group("/customers")
	{
		customers.GET(
			"/:apiKey",
			s.CustomersController.GetOneByApiKey,
		)
		customers.POST("/", s.CustomersController.Create)
	}

	return router
}
