// package contains unit tests and integration tests for the app
package tests

import (
	"database/sql"
	"errors"
	db "go-nauka/location-service/db"
	"go-nauka/location-service/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// initializesz a mock database for testing
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize sqlmock: %v", err)
	}
	db.DB = mockDB

	cleanup := func() {
		mockDB.Close()
	}

	return mock, cleanup
}

// tests the AddLocation function for adding or updating a location in the database
func TestAddLocation(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	tests := []struct {
		name      string
		location  models.Location
		mockError error
		wantErr   bool
	}{
		{
			name: "Successful Insert",
			location: models.Location{
				Name:      "john_doe",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "Failed Insert (DB Error)",
			location: models.Location{
				Name:      "error_user",
				Latitude:  34.0522,
				Longitude: -118.2437,
			},
			mockError: errors.New("insert error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT name FROM location WHERE name = ?").
				WithArgs(tt.location.Name).
				WillReturnError(sql.ErrNoRows)

			mock.ExpectExec("INSERT INTO location").
				WithArgs(tt.location.Name, tt.location.Latitude, tt.location.Longitude).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockError)

			_, err := db.AddLocation(tt.location)

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

// tests the DBGETLocations function for retrieving user locations for the database
func TestDBGetLocations(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude", "updated_at"}).
		AddRow("john_doe", 40.7128, -74.0060, "2024-01-16 10:00:00").
		AddRow("jane_doe", 34.0522, -118.2437, "2024-01-16 11:00:00")

	mock.ExpectQuery("SELECT \\* FROM location").
		WillReturnRows(rows)

	locations, err := db.DBGetLocations()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}

	if locations[0].Name != "john_doe" || locations[1].Name != "jane_doe" {
		t.Errorf("Returned locations do not match expected values")
	}
}

// tests the SearchLocations function for finding users within a specified radius
func TestSearchLocations(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	lat, lon, radius := 40.7128, -74.0060, 10.0
	page, pageSize := 1, 5
	offset := (page - 1) * pageSize

	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude", "updated_at", "distance"}).
		AddRow("john_doe", 40.7128, -74.0060, "2024-01-16 10:00:00", 5.0).
		AddRow("jane_doe", 40.7306, -73.9352, "2024-01-16 11:00:00", 8.0)

	mock.ExpectQuery("SELECT name, latitude, longitude, updated_at,").
		WithArgs(lat, lon, lat, radius, pageSize, offset).
		WillReturnRows(rows)

	locations, err := db.SearchLocations(lat, lon, radius, page, pageSize)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(locations))
	}
}
