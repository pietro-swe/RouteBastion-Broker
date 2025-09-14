package customer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/crypto"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

func MakeCreateCustomerHandler(
	encryptionKey []byte,
	tracer trace.Tracer,
	provider db.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		body := shared.CreateCustomerInput{}
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})

			return
		}

		store := NewCustomersStore(provider)
		apiKeyGen := crypto.NewHashGenerator(encryptionKey)
		txManager := db.NewPgTxManager(provider.GetConn())

		traceCtx, span := tracer.Start(ctx, "customer.CreateCustomer")
		customer, err := CreateCustomer(
			traceCtx,
			txManager,
			store,
			apiKeyGen,
			body,
		)
		span.End()

		if err != nil {
			switch e := err.(type) {
			case errors.DomainError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(),
				})
			case errors.ApplicationError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(),
				})
			case errors.InfrastructureError:
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": e.Error(),
				})
			}

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":                 customer.ID,
			"name":               customer.Name,
			"businessIdentifier": customer.BusinessIdentifier,
			"apiKey":             customer.APIKey,
			"createdAt":          customer.CreatedAt,
			"updatedAt":          customer.UpdatedAt,
			"deletedAt":          customer.DeletedAt,
		})
	}
}

func MakeGetOneByAPIKeyHandler(
	tracer trace.Tracer,
	provider db.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		apiKey := c.Param("apiKey")

		traceCtx, span := tracer.Start(ctx, "customer.GetOneCustomerByAPIKey")
		store := NewCustomersStore(provider)
		foundCustomer, err := GetOneCustomerByAPIKey(
			traceCtx,
			store,
			apiKey,
		)
		span.End()

		if err != nil {
			switch e := err.(type) {
			case errors.DomainError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(),
				})
			case errors.ApplicationError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(),
				})
			case errors.InfrastructureError:
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": e.Error(),
				})
			}

			return
		}

		if foundCustomer == nil {
			message := fmt.Sprintf("Customer for API key %s not found", apiKey)

			c.JSON(http.StatusNotFound, gin.H{
				"error": message,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                 foundCustomer.ID,
			"name":               foundCustomer.Name,
			"businessIdentifier": foundCustomer.BusinessIdentifier,
			"apiKey":             foundCustomer.APIKey,
			"createdAt":          foundCustomer.CreatedAt,
			"updatedAt":          foundCustomer.UpdatedAt,
			"deletedAt":          foundCustomer.DeletedAt,
		})
	}
}
