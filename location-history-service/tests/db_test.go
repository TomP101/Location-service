// package conatins unit an integration tests for the app
package tests

import (
	"errors"

	DB "go-nauka/location-history-service/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// initializes a mock database for testing
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize sqlmock: %v", err)
	}
	DB.DB = mockDB

	cleanup := func() {
		mockDB.Close()
	}

	return mock, cleanup
}

// tests the SaveLocation func
func TestSaveLocation(t *testing.T) {
	mock, _ := setupMockDB(t)

	tests := []struct {
		name       string
		username   string
		latitude   float64
		longitude  float64
		recordedAt string
		mockError  error
		wantErr    bool
	}{
		{
			name:       "Successful Insert",
			username:   "john_doe",
			latitude:   40.7128,
			longitude:  -74.0060,
			recordedAt: "2024-01-16T10:00:00Z",
			mockError:  nil,
			wantErr:    false,
		},
		{
			name:       "Failed Insert (DB Error)",
			username:   "error_user",
			latitude:   34.0522,
			longitude:  -118.2437,
			recordedAt: "2024-01-16T11:00:00Z",
			mockError:  errors.New("insert error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO location_history").
				WithArgs(tt.username, tt.latitude, tt.longitude, tt.recordedAt).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockError)

			err := DB.SaveLocation(tt.username, tt.latitude, tt.longitude, tt.recordedAt)

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled DB expectations: %v", err)
			}
		})
	}
}

// Tests the behavior of GetUserLocations func
func TestGetUserLocations(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	username := "john_doe"
	startDate := "2024-01-16T00:00:00Z"
	endDate := "2024-01-17T00:00:00Z"

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		expectedLen int
		wantErr     bool
	}{
		{
			name: "Successful Fetch",
			mockRows: sqlmock.NewRows([]string{"id", "username", "latitude", "longitude", "recorded_at"}).
				AddRow(1, username, 40.7128, -74.0060, "2024-01-16T10:00:00Z").
				AddRow(2, username, 40.7138, -74.0070, "2024-01-16T11:00:00Z"),
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "No Data Found",
			mockRows:    sqlmock.NewRows([]string{"id", "username", "latitude", "longitude", "recorded_at"}),
			expectedLen: 0,
			wantErr:     false,
		},
		{
			name:        "Query Error",
			mockRows:    nil,
			expectedLen: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockRows != nil {
				mock.ExpectQuery("SELECT id, username, latitude, longitude, recorded_at FROM location_history").
					WithArgs(username, startDate, endDate).
					WillReturnRows(tt.mockRows)
			} else {
				mock.ExpectQuery("SELECT id, username, latitude, longitude, recorded_at FROM location_history").
					WithArgs(username, startDate, endDate).
					WillReturnError(errors.New("query error"))
			}

			history, err := DB.GetUserLocations(username, startDate, endDate)

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			if len(history) != tt.expectedLen {
				t.Errorf("Expected %d records, got %d", tt.expectedLen, len(history))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled DB expectations: %v", err)
			}
		})
	}
}
