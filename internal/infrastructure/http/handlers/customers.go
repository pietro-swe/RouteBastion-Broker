package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	usecases "github.com/marechal-dev/RouteBastion-Broker/internal/application/use_cases"
	cryptoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/cryptography/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/http/presenters"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	repoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/shared"
)

func MakeCreateCustomerHandler(
	encryptionKey []byte,
	tracer trace.Tracer,
	db persistence.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		body := dtos.CreateCustomerInput{}
		err := c.BindJSON(&body)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})

			return
		}

		repository := repoImpl.NewPgCustomersRepository(db)
		apiKeyGen := cryptoImpl.NewAPIKeyGenerator()
		txManager := persistence.NewPgTxManager(db.GetConn())

		useCase := usecases.NewCreateCustomerUseCase(
			txManager,
			repository,
			encryptionKey,
			apiKeyGen,
		)

		traceCtx, span := tracer.Start(ctx, "CreateCustomerUseCase.Execute")
		customer, err := useCase.Execute(traceCtx, body)
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
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": e.Error(),
				})
			}

			return
		}

		payload := presenters.CustomerToHTTP(customer)

		c.JSON(http.StatusCreated, payload)
	}
}

func MakeGetOneByAPIKeyHandler(
	db persistence.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Param("apiKey")

		repository := repoImpl .NewPgCustomersRepository(db)
		useCase := usecases.NewGetOneCustomerUseCaseImpl(repository)

		foundCustomer := useCase.Execute(apiKey)

		if foundCustomer == nil {
			message := fmt.Sprintf("Customer for API key %s not found", apiKey)

			c.JSON(http.StatusNotFound, gin.H{
				"error": message,
			})

			return
		}

		payload := presenters.CustomerToHTTP(foundCustomer)

		c.JSON(http.StatusOK, payload)
	}
}
