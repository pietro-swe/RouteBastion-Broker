package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
)

func MakeHealthCheckHandler(db persistence.DBProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"database": db.Health(),
		})
	}
}
