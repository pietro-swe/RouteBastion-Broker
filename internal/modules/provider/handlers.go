package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/customerrors"
	uuid "github.com/satori/go.uuid"
)

func CreateProviderHandler(
	provider db.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		body := shared.CreateProviderInput{}
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		store := NewProvidersStore(provider)

		providerEntity, err := CreateProvider(
			ctx,
			store,
			body,
		)
		if err != nil {
			switch e := err.(type) {
			case *customerrors.AppError:
				c.JSON(customerrors.ToHttpStatus(e), gin.H{
					"code":    e.Code,
					"message": e.Msg,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    customerrors.ErrUnknown,
					"message": e.Error(),
				})
			}

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         providerEntity.ID,
			"name":       providerEntity.Name,
			"createdAt":  providerEntity.CreatedAt,
			"modifiedAt": providerEntity.ModifiedAt,
			"deletedAt":  providerEntity.DeletedAt,
		})
	}
}

func UpdateProviderHandler(
	provider db.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id := c.Param("id")

		parsedID, err := uuid.FromString(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  customerrors.ErrCodeInvalidInput,
				"error": "invalid provider ID",
			})

			return
		}

		body := shared.UpdateProviderInput{
			ID: parsedID.String(),
		}

		bindErr := c.BindJSON(&body)
		if bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    customerrors.ErrCodeInvalidInput,
				"message": bindErr.Error(),
			})

			return
		}

		store := NewProvidersStore(provider)

		providerEntity, err := UpdateProvider(
			ctx,
			store,
			body,
		)
		if err != nil {
			switch e := err.(type) {
			case *customerrors.AppError:
				c.JSON(customerrors.ToHttpStatus(e), gin.H{
					"code":    e.Code,
					"message": e.Msg,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    customerrors.ErrUnknown,
					"message": e.Error(),
				})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         providerEntity.ID,
			"name":       providerEntity.Name,
			"createdAt":  providerEntity.CreatedAt,
			"modifiedAt": providerEntity.ModifiedAt,
			"deletedAt":  providerEntity.DeletedAt,
		})
	}
}

func DeleteProviderHandler(
	provider db.DBProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id := c.Param("id")

		parsedID, err := uuid.FromString(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    customerrors.ErrCodeInvalidInput,
				"message": "invalid provider ID",
			})

			return
		}

		store := NewProvidersStore(provider)

		deleteErr := DeleteProvider(
			ctx,
			store,
			parsedID.String(),
		)
		if deleteErr != nil {
			switch e := deleteErr.(type) {
			case *customerrors.AppError:
				c.JSON(customerrors.ToHttpStatus(e), gin.H{
					"code":    e.Code,
					"message": e.Msg,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    customerrors.ErrUnknown,
					"message": e.Error(),
				})
			}

			return
		}

		c.Status(http.StatusNoContent)
	}
}
