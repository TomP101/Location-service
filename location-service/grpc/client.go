// package manages the GRPC communcation with the location-history-service
package GRPC

import (
	"context"
	"log"
	"time"

	pb "go-nauka/location-service/grpc/proto"

	"google.golang.org/grpc"
)

// defines the interface for sending location updates over gRPC
type GRPCClient interface {
	SendLocationUpdate(username string, latitude, longitude float64) error
}

// implments the GRPCCLIENT interface using the grpc generated client
type DefaultGRPCClient struct {
	client pb.LocationHistoryServiceClient
}

// initializesz the GRPC client and connects to the
func InitGRPCClient() *DefaultGRPCClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to Location History Service: %v", err)
	}
	client := pb.NewLocationHistoryServiceClient(conn)
	log.Println("Connected to Location History Service (gRPC)")

	return &DefaultGRPCClient{client: client}
}

// send a user location update to location-history-service over grpc
func (d *DefaultGRPCClient) SendLocationUpdate(username string, latitude, longitude float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LocationRequest{
		Username:   username,
		Latitude:   latitude,
		Longitude:  longitude,
		RecordedAt: time.Now().Format(time.RFC3339),
	}

	_, err := d.client.RecordLocation(ctx, req)
	if err != nil {
		log.Printf("Failed to send location update: %v", err)
		return err
	}

	log.Printf("Sent location update for user '%s'", username)
	return nil
}

// global grpc client instance
var Client GRPCClient = InitGRPCClient()
