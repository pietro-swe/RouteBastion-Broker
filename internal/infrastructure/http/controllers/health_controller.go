package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
)

type HealthController interface {
	Index(c *gin.Context)
}

type healthController struct {
	db persistence.DBProvider
}

func NewHealthController(db persistence.DBProvider) HealthController {
	return &healthController{
		db: db,
	}
}

func (h *healthController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"db": h.db.Health(),
	})
}
