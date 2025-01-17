// package initalizesz and starts the GRPCS and REST servers
package main

import (
	"log"
	"net"

	"go-nauka/location-history-service/db"
	"go-nauka/location-history-service/grpc"
	"go-nauka/location-history-service/routes"

	pb "go-nauka/location-history-service/grpc/proto"

	gr "google.golang.org/grpc"
)

// initalizes the database connection and starts GRPC and REST servers
func main() {
	db.InitDB()

	go startGRPCServer()

	startRESTServer()

}

// starts the GRPC server on port 50051
func startGRPCServer() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := gr.NewServer()
	pb.RegisterLocationHistoryServiceServer(grpcServer, &grpc.Server{})

	log.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

// starts the REST server on localhost port 8081 (8080 used by location-service microservice)
func startRESTServer() {
	router := routes.SetupRouter()
	router.Run("localhost:8081")
}
