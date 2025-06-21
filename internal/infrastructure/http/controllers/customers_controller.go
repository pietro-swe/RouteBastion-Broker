package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	usecases "github.com/marechal-dev/RouteBastion-Broker/internal/application/use_cases"
	cryptoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/cryptography/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/presenters"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories/implementations"
	repoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/shared"
)

type CustomersController struct {
	encrytionKey []byte
	db persistence.DBProvider
}

func NewCustomersController(encrytionKey []byte, db persistence.DBProvider) CustomersController {
	return CustomersController{
		encrytionKey: encrytionKey,
		db: db,
	}
}

func (cc *CustomersController) Create(c *gin.Context) {
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

	ctx := context.Background()
	customer, err := useCase.Execute(ctx, dto)
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

func (cc *CustomersController) GetOneByApiKey(c *gin.Context) {
	apiKey := c.Param("apiKey")

	repository := implementations.NewPgCustomersRepository(cc.db)
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
