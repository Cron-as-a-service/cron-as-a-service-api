package routes

import (
	"frderoubaix.me/cron-as-a-service/services"
	"github.com/gin-gonic/gin"
)

func FetchingRoute(router *gin.Engine) {
	router.GET("/api/v1/fetching/attributes", services.AttributesEndpoint)
}
