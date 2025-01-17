// package handles all database actions for managing user location data in microservice 1
package DB

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go-nauka/location-service/models"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// safely closes the database connection
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("Error closing the database:", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}

// initializes the database connection
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

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to the database!")
}

// Gets all user location records from the database
func DBGetLocations() ([]models.Location, error) {
	var locations []models.Location

	rows, err := DB.Query("SELECT * FROM location")
	if err != nil {
		return nil, fmt.Errorf("locations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.Name, &loc.Latitude, &loc.Longitude, &loc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("locations: %v", err)
		}

		locations = append(locations, loc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("locations: %v", err)
	}
	return locations, nil
}

// inserts a new location or updates one if it exists( name )
func AddLocation(loc models.Location) (int64, error) {

	var existingName string
	err := DB.QueryRow("SELECT name FROM location WHERE name = ?", loc.Name).Scan(&existingName)

	if err != nil && err != sql.ErrNoRows {

		return 0, fmt.Errorf("addLocation: %v", err)
	}

	if err == nil {

		_, err := DB.Exec("UPDATE location SET latitude = ?, longitude = ?, updated_at = CURRENT_TIMESTAMP WHERE name = ?", loc.Latitude, loc.Longitude, loc.Name)
		if err != nil {
			return 0, fmt.Errorf("addLocation (update): %v", err)
		}
		fmt.Printf("Updated location for '%s'\n", loc.Name)
		return 0, nil
	}

	result, err := DB.Exec("INSERT INTO location (name, latitude, longitude) VALUES (?, ?, ?)", loc.Name, loc.Latitude, loc.Longitude)
	if err != nil {
		return 0, fmt.Errorf("addLocation (insert): %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("addLocation (insert): %v", err)
	}
	fmt.Printf("Inserted new location for '%s'\n", loc.Name)
	return rowsAffected, nil
}

// retrives locations within a specified radius of given coordinates(supports pagination)
func SearchLocations(lat, lon, radius float64, page, pageSize int) ([]models.Location, error) {
	var locations []models.Location

	offset := (page - 1) * pageSize

	query := `
	SELECT name, latitude, longitude, updated_at,
		(6371 * ACOS(
			COS(RADIANS(?)) * COS(RADIANS(latitude)) *
			COS(RADIANS(longitude) - RADIANS(?)) +
			SIN(RADIANS(?)) * SIN(RADIANS(latitude))
		)) AS distance
	FROM location
	HAVING distance <= ?
	ORDER BY distance ASC
	LIMIT ? OFFSET ?
`

	rows, err := DB.Query(query, lat, lon, lat, radius, pageSize, offset)

	if err != nil {
		return nil, fmt.Errorf("searchLocations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc models.Location
		var distance float64
		if err := rows.Scan(&loc.Name, &loc.Latitude, &loc.Longitude, &loc.UpdatedAt, &distance); err != nil {
			return nil, fmt.Errorf("searchLocations: %v", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("searchLocations: %v", err)
	}

	return locations, nil
}
