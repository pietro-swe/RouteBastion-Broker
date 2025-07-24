/*
Package handlers provides HTTP handlers for different routes
*/
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecases "github.com/marechal-dev/RouteBastion-Broker/internal/application/use_cases"
	"go.opentelemetry.io/otel/trace"
)

func MakeOptimizeSyncHandler(
	tracer trace.Tracer,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		useCase := usecases.NewOptimizeSyncUseCaseImpl()

		_, span := tracer.Start(ctx, "OptimizeSyncUseCase.Execute")
		useCase.Execute(ctx)
		span.End()

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
