// handles all database opeartions for the location-history-service
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go-nauka/location-history-service/models"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// initalizes the database
func InitDB() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "users",
	}

	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the location history database!")
}

// inserts a new location record into the location_history table
func SaveLocation(username string, lat, lon float64, recordedAt string) error {
	_, err := DB.Exec("INSERT INTO location_history (username, latitude, longitude, recorded_at) VALUES (?, ?, ?, ?)",
		username, lat, lon, recordedAt)
	return err
}

// retrieves a users location history between two dates
func GetUserLocations(username, startDate, endDate string) ([]models.LocationHistory, error) {
	query := `
		SELECT id, username, latitude, longitude, recorded_at 
		FROM location_history 
		WHERE username = ? AND recorded_at BETWEEN ? AND ?
		ORDER BY recorded_at ASC
	`

	rows, err := DB.Query(query, username, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("GetUserLocations: %v", err)
	}
	defer rows.Close()

	var history []models.LocationHistory

	for rows.Next() {
		var loc models.LocationHistory
		if err := rows.Scan(&loc.ID, &loc.Username, &loc.Latitude, &loc.Longitude, &loc.RecordedAt); err != nil {
			return nil, fmt.Errorf("GetUserLocations: %v", err)
		}
		history = append(history, loc)
	}

	return history, nil
}
