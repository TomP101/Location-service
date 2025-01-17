// package contains unit tests and integration tests for the app
package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	GRPC "go-nauka/location-service/grpc"
	"go-nauka/location-service/handlers"
	"go-nauka/location-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// simulates grpcs client behavior for testing
type MockGRPCClient struct {
	ShouldFail bool
}

// mocks the behavior of sending a location update over grpc
func (m *MockGRPCClient) SendLocationUpdate(username string, latitude, longitude float64) error {
	if m.ShouldFail {
		return errors.New("mocked gRPC failure")
	}
	return nil
}

// tests the POST /locations endpoint for adding new locations
func TestPostLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	router := gin.Default()
	router.POST("/locations", handlers.PostLocation)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		grpcShouldFail bool
		mockDB         bool
	}{
		{
			name:           "Successful gRPC Call",
			payload:        `{"name":"tomek_prus","latitude":40.7128,"longitude":-74.0060}`,
			expectedStatus: http.StatusCreated,
			grpcShouldFail: false,
			mockDB:         true,
		},
		{
			name:           "gRPC Failure",
			payload:        `{"name":"tomek_prus","latitude":40.7128,"longitude":-74.0060}`,
			expectedStatus: http.StatusInternalServerError,
			grpcShouldFail: true,
			mockDB:         true,
		},
		{
			name:           "Invalid Payload",
			payload:        `{"name":"","latitude":999,"longitude":999}`,
			expectedStatus: http.StatusBadRequest,
			grpcShouldFail: false,
			mockDB:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GRPC.Client = &MockGRPCClient{ShouldFail: tt.grpcShouldFail}

			var location models.Location
			_ = json.Unmarshal([]byte(tt.payload), &location)

			if tt.mockDB {
				mock.ExpectQuery("SELECT name FROM location WHERE name = ?").
					WithArgs("tomek_prus").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectExec("INSERT INTO location").
					WithArgs("tomek_prus", 40.7128, -74.0060).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			req, _ := http.NewRequest("POST", "/locations", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.mockDB {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("Unfulfilled DB expectations: %v", err)
				}
			}
		})
	}
}

// tests the GET /locations endpoint for retrieving all stored locations
func TestGetLocations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	router := gin.Default()
	router.GET("/locations", handlers.GetLocations)

	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude", "updated_at"}).
		AddRow("tomek_prus", 40.7128, -74.0060, "2024-01-16 10:00:00").
		AddRow("jane_doe", 34.0522, -118.2437, "2024-01-16 11:00:00")

	mock.ExpectQuery("SELECT \\* FROM location").
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/locations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled DB expectations: %v", err)
	}
}

// tests the GET /search endpoint for finding locations within a specified radius
func TestSearchLocationsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	router := gin.Default()
	router.GET("/search", handlers.SearchLocationsHandler)

	lat, lon, radius, page, pageSize := 40.7128, -74.0060, 10.0, 1, 5
	offset := (page - 1) * pageSize

	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude", "updated_at", "distance"}).
		AddRow("tomek_prus", 40.7128, -74.0060, "2024-01-16 10:00:00", 5.0).
		AddRow("jane_doe", 40.7306, -73.9352, "2024-01-16 11:00:00", 8.0)

	mock.ExpectQuery("SELECT name, latitude, longitude, updated_at,").
		WithArgs(lat, lon, lat, radius, pageSize, offset).
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/search?latitude=40.7128&longitude=-74.0060&radius=10&page=1&page_size=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled DB expectations: %v", err)
	}
}
