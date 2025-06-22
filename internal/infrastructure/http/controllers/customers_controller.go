/*
Package controllers provide HTTP Request Controllers for the Broker REST API
*/
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	usecases "github.com/marechal-dev/RouteBastion-Broker/internal/application/use_cases"
	cryptoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/cryptography/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/presenters"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	repoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/shared"
	"go.opentelemetry.io/otel/trace"
)

type CustomersController struct {
	encrytionKey []byte
	tracer trace.Tracer
	db persistence.DBProvider
}

func NewCustomersController(encrytionKey []byte, tracer trace.Tracer, db persistence.DBProvider) CustomersController {
	return CustomersController{
		encrytionKey: encrytionKey,
		tracer: tracer,
		db: db,
	}
}

func (cc *CustomersController) Create(c *gin.Context) {
	requestCtx := c.Request.Context()

	dto := &dtos.CreateCustomerInput{}
	err := c.BindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload",
		})

		return
	}

	repository := repoImpl.NewPgCustomersRepository(cc.db)
	apiKeyGen := cryptoImpl.NewAPIKeyGenerator()
	txManager := persistence.NewPgTxManager(cc.db.GetConn())

	useCase := usecases.NewCreateCustomerUseCase(
		txManager,
		repository,
		cc.encrytionKey,
		apiKeyGen,
	)


	traceCtx, span := cc.tracer.Start(requestCtx, "CreateCustomerUseCaseImpl.Execute")
	customer, err := useCase.Execute(traceCtx, dto)
	span.End()
	if err != nil {
		switch e := err.(type) {
		case shared.DomainError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": e.Error(),
			})
		case shared.ApplicationError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": e.Error(),
			})
		case shared.InfrastructureError:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": e.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": e.Error(),
			})
		}

		return
	}

	payload := presenters.CustomerFromDomain(customer)

	c.JSON(http.StatusCreated, payload)
}

func (cc *CustomersController) GetOneByAPIKey(c *gin.Context) {
	apiKey := c.Param("apiKey")

	repository := repoImpl.NewPgCustomersRepository(cc.db)
	useCase := usecases.NewGetOneCustomerUseCaseImpl(repository)

	foundCustomer := useCase.Execute(apiKey)

	if foundCustomer == nil {
		message := fmt.Sprintf("customer for API key %s not found", apiKey)

		c.JSON(http.StatusNotFound, map[string]string{
			"error": message,
		})

		return
	}

	payload := presenters.CustomerFromDomain(foundCustomer)

	c.JSON(http.StatusOK, payload)
}
