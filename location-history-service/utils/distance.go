// package provides utility functions for distance calculation
package utils

import (
	"go-nauka/location-history-service/models"
	"math"
)

const EarthRadius = 6371.0

// converts degrees to radians(used for calculations using trigonometry)
func DegreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// calculates the distance between two points
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Asin(math.Sqrt(a))

	return EarthRadius * c
}

// calculates total distance traveled by a user
func CalculateTotalDistance(locations []models.LocationHistory) float64 {
	var total float64
	for i := 1; i < len(locations); i++ {
		total += HaversineDistance(
			locations[i-1].Latitude, locations[i-1].Longitude,
			locations[i].Latitude, locations[i].Longitude,
		)
	}
	return total
}
