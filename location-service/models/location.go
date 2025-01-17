// package defines data structures used in the app
package models

// Location represents the structure for storing user location data
// Name - unique user identifier
// Latitude and Longitude are used for defininf a users location
// UpdatedAt is used to record the time of the last update
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}
