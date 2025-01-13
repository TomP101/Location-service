package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// represents data about a user
type User struct {
	ID        string  `json:"id"`
	name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var users = []User{
	{ID: "1", name: "Blue Train", Latitude: 35.12314, Longitude: 27.64532},
	{ID: "2", name: "Jeru", Latitude: 36.12314, Longitude: 28.64532},
	{ID: "3", name: "Sarah Vaughan and Clifford Brown", Latitude: 37.12314, Longitude: 29.64532},
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
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	usID, err := addUser(User{
		name:      "tomek",
		Latitude:  7.112323,
		Longitude: 123.123,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added user: %v\n", usID)

	users, err := db_get_users()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Users found %v\n", users)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	fmt.Println("works")
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", postUsers)

	router.Run("localhost:8080")

}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func postUsers(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func db_get_users() ([]User, error) {

	var users []User

	rows, err := db.Query("SELECT * FROM user ")
	if err != nil {
		return nil, fmt.Errorf("users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var us User
		if err := rows.Scan(&us.ID, &us.name, &us.Latitude, &us.Longitude); err != nil {
			return nil, fmt.Errorf("users: %v", err)
		}
		users = append(users, us)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("users: %v", err)
	}
	return users, nil
}

func addUser(us User) (int64, error) {
	result, err := db.Exec("INSERT INTO user (name, latitude, longitude) VALUES (?, ?, ?)", us.name, us.Latitude, us.Longitude)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}
