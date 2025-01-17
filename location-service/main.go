// package initializes and starts the location-service
package main

import (
	"context"
	"fmt"
	DB "go-nauka/location-service/db"
	grpc "go-nauka/location-service/grpc"
	"go-nauka/location-service/models"
	"go-nauka/location-service/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// initalizes the database, grpc client and starts the http server
func main() {
	DB.InitDB()

	go grpc.InitGRPCClient()

	locID, err := DB.AddLocation(models.Location{
		Name:      "antek",
		Latitude:  80.112323,
		Longitude: 120.123,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added location: %v\n", locID)

	locations, err := DB.DBGetLocations()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Locations found %v\n", locations)

	pingErr := DB.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	router := routes.SetupRouter()
	go router.Run("localhost:8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	gracefulShutdown(server)

}

// shuts down the http server and closes the database
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	DB.CloseDB()
	fmt.Println("Server exited gracefully")
}
