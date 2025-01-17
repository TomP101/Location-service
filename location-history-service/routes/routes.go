// package provides routing configuration for the location-history-service
package routes

import (
	"go-nauka/location-history-service/handlers"

	"github.com/gin-gonic/gin"
)

// initalizes the router and defines HTTP routes for the server
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/history/distance", handlers.CalculateDistance)
	return router
}
