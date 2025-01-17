// package conatins unit an integration tests for the app
package tests

import (
	"go-nauka/location-history-service/utils"
	"math"
	"testing"
)

// checks whether a and b are almost equal within the provided tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// tests the HaversineDistance func
func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{
			name:     "London to New York",
			lat1:     51.5007,
			lon1:     0.1246,
			lat2:     40.6892,
			lon2:     -74.0445,
			expected: 5591,
		},
		{
			name:     "Same Location",
			lat1:     34.0522,
			lon1:     -118.2437,
			lat2:     34.0522,
			lon2:     -118.2437,
			expected: 0,
		},
		{
			name:     "Task Example",
			lat1:     35.12314,
			lon1:     27.64532,
			lat2:     39.12355,
			lon2:     27.64538,
			expected: 445,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if !almostEqual(result, tt.expected, 1.0) {
				t.Errorf("Expected %.2f km, but got %.2f km", tt.expected, result)
			}
		})
	}
}
