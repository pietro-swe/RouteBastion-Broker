package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	clients "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/clients/implementations"
	"go.opentelemetry.io/otel/trace"
)

type OptimizationsController struct {
	tracer trace.Tracer
}

func NewOptimizationsController(
	tracer trace.Tracer,
) OptimizationsController {
	return OptimizationsController{
		tracer: tracer,
	}
}

func (oc *OptimizationsController) Optimize(c *gin.Context) {
	requestCtx := c.Request.Context()

	client := clients.GoogleCloudClient{}

	_, span := oc.tracer.Start(requestCtx, "GoogleCloudClient.Optimize")
	client.Optimize(&dtos.OptimizationRequestInput{})
	span.End()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
