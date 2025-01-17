// defines routes and initializes the api routes for the location service
package routes

import (
	"go-nauka/location-service/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main Gin router with defined routes
// GET  /locations - Retrieves all stored user locations
// POST /locations - Adds or updates a users location and notifies the history service
// GET  /search    - Searches for users within a specified radius with pagination support
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/locations", handlers.GetLocations)
	router.POST("/locations", handlers.PostLocation)
	router.GET("/search", handlers.SearchLocationsHandler)

	return router
}
