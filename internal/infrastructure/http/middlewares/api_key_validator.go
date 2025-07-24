/*
Package middlewares provide HTTP server middlewares to run before requests enter the main handler
*/
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	repoImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories/implementations"
)

func WithValidAPIKey(db persistence.DBProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const apiKeyHeader string = "RouteBastion-API-Key"

		apiKey := ctx.GetHeader(apiKeyHeader)

		if apiKey == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "API Key header is missing",
				},
			)

			return
		}

		repository := repoImpl.NewPgCustomersRepository(db)
		customer := repository.GetOneByAPIKey(ctx, apiKey)

		if customer == nil {
			ctx.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "Invalid API key",
				},
			)

			return
		}

		ctx.Next()
	}
}
