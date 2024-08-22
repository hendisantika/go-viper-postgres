package routers

import (
	"github.com/gin-gonic/gin"
	"go-viper-postgres/src/api/handlers"
)

func Health(r *gin.RouterGroup) {
	handler := handlers.NewHealthHandler()

	r.GET("/", handler.Health)
}
