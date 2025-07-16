package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	// client := clients.GoogleCloudClient{}

	_, span := oc.tracer.Start(requestCtx, "CloudClientSimulator.Optimize")
	// client.Optimize(&dtos.OptimizationRequestInput{})
	time.Sleep(time.Microsecond * 10)
	span.End()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
