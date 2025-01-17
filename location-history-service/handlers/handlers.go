// handlers provides HTTP handlers for processing location history data
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-nauka/location-history-service/db"
	"go-nauka/location-history-service/utils"

	"github.com/gin-gonic/gin"
)

// handles GET requests for calculating the total distance traveled by the user
func CalculateDistance(c *gin.Context) {
	username := c.Query("username")
	startDate := c.DefaultQuery("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
	endDate := c.DefaultQuery("end", time.Now().Format(time.RFC3339))

	locations, err := db.GetUserLocations(username, startDate, endDate)
	if err != nil {
		log.Printf("Error fetching user locations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch locations"})
		return
	}

	totalDistance := utils.CalculateTotalDistance(locations)
	c.JSON(http.StatusOK, gin.H{
		"username":      username,
		"totalDistance": fmt.Sprintf("%.2f km", totalDistance),
	})
}
