package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lyyubava/solidgate-software-engineering-school.git/controllers"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	routerApi := r.Group("/api")
	routerApi.GET("/rate", controllers.Rate)
	routerApi.POST("/subscribe", controllers.Subscribe)
	return r
}
