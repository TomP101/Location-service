package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

const EarthRadius = 6371.0

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

var locations = []Location{
	{Name: "Blue Train", Latitude: 35.12314, Longitude: 27.64532},
	{Name: "Jeru", Latitude: 36.12314, Longitude: 28.64532},
	{Name: "Sarah Vaughan and Clifford Brown", Latitude: 37.12314, Longitude: 29.64532},
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "users",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	locID, err := addLocation(Location{
		Name:      "antek",
		Latitude:  80.112323,
		Longitude: 120.123,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added location: %v\n", locID)

	locations, err := dbGetLocations()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Locations found %v\n", locations)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/locations", getLocations)
	router.POST("/locations", postLocation)
	router.GET("/search", searchLocationsHandler)

	router.Run("localhost:8080")
}

func getLocations(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, locations)
}

func postLocation(c *gin.Context) {
	var newLocation Location

	if err := c.BindJSON(&newLocation); err != nil {
		return
	}

	locations = append(locations, newLocation)
	c.IndentedJSON(http.StatusCreated, newLocation)
}

func dbGetLocations() ([]Location, error) {
	var locations []Location

	rows, err := db.Query("SELECT * FROM location")
	if err != nil {
		return nil, fmt.Errorf("locations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc Location
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

func addLocation(loc Location) (int64, error) {

	var existingName string
	err := db.QueryRow("SELECT name FROM location WHERE name = ?", loc.Name).Scan(&existingName)

	if err != nil && err != sql.ErrNoRows {

		return 0, fmt.Errorf("addLocation: %v", err)
	}

	if err == nil {

		_, err := db.Exec("UPDATE location SET latitude = ?, longitude = ?, updated_at = CURRENT_TIMESTAMP WHERE name = ?", loc.Latitude, loc.Longitude, loc.Name)
		if err != nil {
			return 0, fmt.Errorf("addLocation (update): %v", err)
		}
		fmt.Printf("Updated location for '%s'\n", loc.Name)
		return 0, nil
	}

	result, err := db.Exec("INSERT INTO location (name, latitude, longitude) VALUES (?, ?, ?)", loc.Name, loc.Latitude, loc.Longitude)
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

func ValidateName(name string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9 ]{4,128}$`)
	if !re.MatchString(name) {
		return errors.New("invalid name")
	}
	return nil
}

func ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return errors.New("invalid coordinates")
	}
	return nil
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

func searchLocations(lat, lon, radius float64, page, pageSize int) ([]Location, error) {
	var locations []Location

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

	rows, err := db.Query(query, lat, lon, lat, radius, pageSize, offset)

	if err != nil {
		return nil, fmt.Errorf("searchLocations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc Location
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

func searchLocationsHandler(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	locations, err := searchLocations(lat, lon, radius, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}
