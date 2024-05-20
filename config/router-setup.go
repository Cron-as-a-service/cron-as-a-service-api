package config

import (
	"fmt"
	"frderoubaix.me/cron-as-a-service/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	ginConfig "github.com/tavsec/gin-healthcheck/config"
	"go.uber.org/zap"
	"os"
)

func InitRouter() {
	router := gin.Default()

	// Configurer le middleware CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	routes.GlobalRoute(router)
	routes.TaskRoute(router)
	routes.FetchingRoute(router)

	healthcheckError := healthcheck.New(router, ginConfig.DefaultConfig(), []checks.Check{})
	if healthcheckError != nil {
		zap.L().Error(fmt.Sprintf("Healthcheck initialization failed: %v", healthcheckError))
	}

	port := os.Getenv("SERVER_PORT")
	routerError := router.Run(":" + port)
	if routerError != nil {
		zap.L().Error(fmt.Sprintf("GinRouter initialization failed: %v", routerError))
	}
}
