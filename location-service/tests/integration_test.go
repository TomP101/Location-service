// package contains unit tests and integration tests for the app
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// sends a POST request to the location service to update a users location
func postLocation(t *testing.T, username string, latitude, longitude float64) {
	url := "http://localhost:8080/locations"
	payload := map[string]interface{}{
		"name":      username,
		"latitude":  latitude,
		"longitude": longitude,
	}

	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// gets the total distance a user has traveled from the location history service
func getDistance(t *testing.T, username string) float64 {
	url := fmt.Sprintf("http://localhost:8081/history/distance?username=%s", username)
	resp, err := http.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	distanceStr := result["totalDistance"]
	var distance float64
	fmt.Sscanf(distanceStr, "%f km", &distance)
	return distance
}

// tests the entire application workflow
func TestFullIntegration(t *testing.T) {
	username := "integration"

	postLocation(t, username, 35.12314, 27.64532)
	time.Sleep(1 * time.Second)

	postLocation(t, username, 39.12355, 27.64538)
	time.Sleep(1 * time.Second)

	distance := getDistance(t, username)

	assert.InDelta(t, 445.0, distance, 5.0, "Distance should be approximately 445 km")
}
