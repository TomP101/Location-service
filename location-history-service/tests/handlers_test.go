// package conatins unit an integration tests for the app
package tests

import (
	"encoding/json"
	"errors"

	"go-nauka/location-history-service/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// tests the CalculateDistance handler
func TestCalculateDistance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	router := gin.Default()
	router.GET("/history/distance", handlers.CalculateDistance)

	tests := []struct {
		name           string
		username       string
		startDate      string
		endDate        string
		mockRows       *sqlmock.Rows
		mockError      error
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:      "Successful Calculation",
			username:  "john_doe",
			startDate: "2024-01-16T00:00:00Z",
			endDate:   "2024-01-17T00:00:00Z",
			mockRows: sqlmock.NewRows([]string{"id", "username", "latitude", "longitude", "recorded_at"}).
				AddRow(1, "john_doe", 35.12314, 27.64532, "2024-01-16T10:00:00Z").
				AddRow(2, "john_doe", 39.12355, 27.64538, "2024-01-16T12:00:00Z"),
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"username":      "john_doe",
				"totalDistance": "444.83 km",
			},
		},
		{
			name:           "No Data for User",
			username:       "unknown_user",
			startDate:      "2024-01-16T00:00:00Z",
			endDate:        "2024-01-17T00:00:00Z",
			mockRows:       sqlmock.NewRows([]string{"id", "username", "latitude", "longitude", "recorded_at"}),
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"username":      "unknown_user",
				"totalDistance": "0.00 km",
			},
		},
		{
			name:           "Database Error",
			username:       "error_user",
			startDate:      "2024-01-16T00:00:00Z",
			endDate:        "2024-01-17T00:00:00Z",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Could not fetch locations",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError != nil {
				mock.ExpectQuery("SELECT id, username, latitude, longitude, recorded_at FROM location_history").
					WithArgs(tt.username, tt.startDate, tt.endDate).
					WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery("SELECT id, username, latitude, longitude, recorded_at FROM location_history").
					WithArgs(tt.username, tt.startDate, tt.endDate).
					WillReturnRows(tt.mockRows)
			}

			req, _ := http.NewRequest("GET", "/history/distance?username="+tt.username+"&start="+tt.startDate+"&end="+tt.endDate, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled DB expectations: %v", err)
			}
		})
	}
}
