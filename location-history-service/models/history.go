// package defines the data structures used in the location-history-service
package models

// LocationHisotry provides a structure that records a users location at a specific point in time
type LocationHistory struct {
	ID         int     `json:"id"`
	Username   string  `json:"username"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	RecordedAt string  `json:"recorded_at"`
}
