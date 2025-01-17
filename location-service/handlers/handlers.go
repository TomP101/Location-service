// package contains HTTP request handlers for managing user locations
package handlers

import (
	DB "go-nauka/location-service/db"
	GRPC "go-nauka/location-service/grpc"
	"go-nauka/location-service/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var locations []models.Location

// handles GET request for retrievies all stored user locations
// Responds with 500 erro if fetching for the datbase fails
func GetLocations(c *gin.Context) {

	locations, err := DB.DBGetLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
		return
	}
	c.IndentedJSON(http.StatusOK, locations)
}

// handles POST requests for adding or updating user's current location
// validates the input, updates the database and notifies the location-history-service over grpc
func PostLocation(c *gin.Context) {
	var newLocation models.Location

	if err := c.BindJSON(&newLocation); err != nil {
		return
	}

	if newLocation.Name == "" || newLocation.Latitude < -90 || newLocation.Latitude > 90 || newLocation.Longitude < -180 || newLocation.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	_, err := DB.AddLocation(newLocation)
	if err != nil {
		log.Println("Failed to update location in DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	err = GRPC.Client.SendLocationUpdate(newLocation.Name, newLocation.Latitude, newLocation.Longitude)
	if err != nil {
		log.Println("Failed to send gRPC location update:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to notify history service"})
		return
	}

	locations = append(locations, newLocation)
	c.IndentedJSON(http.StatusCreated, newLocation)
}

// handles GET request for searching users within a given radius
// validates query and supports pagination
func SearchLocationsHandler(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	locations, err := DB.SearchLocations(lat, lon, radius, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}
