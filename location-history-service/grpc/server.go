// package implements the GRPC server for handling location updates
package grpc

import (
	"context"

	"go-nauka/location-history-service/db"
	pb "go-nauka/location-history-service/grpc/proto"
)

// implements the LocationHistoryServiceServer interface for handling GRPC requests
type Server struct {
	pb.UnimplementedLocationHistoryServiceServer
}

// handles incoming GRPC requests to store user location data using SaveLocation func from db package
func (s *Server) RecordLocation(ctx context.Context, req *pb.LocationRequest) (*pb.LocationResponse, error) {
	err := db.SaveLocation(req.Username, req.Latitude, req.Longitude, req.RecordedAt)
	if err != nil {
		return &pb.LocationResponse{Status: "Failed"}, err
	}
	return &pb.LocationResponse{Status: "Success"}, nil
}
