package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
)

func MakeHealthCheckHandler(provider db.DBProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"database": provider.Health(),
		})
	}
}
