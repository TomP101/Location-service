// package conatins unit an integration tests for the app
package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// sends a GET request to retrieve the total distance traveled bu a user for testing purposes
func fetchDistance(t *testing.T, username string) float64 {
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

// tests if the calculation done by fetchDistance func is correct
func TestDistanceCalculation(t *testing.T) {
	username := "integration"

	time.Sleep(2 * time.Second)

	distance := fetchDistance(t, username)

	assert.InDelta(t, 445.0, distance, 5.0, "Distance should be approximately 445 km")
}
