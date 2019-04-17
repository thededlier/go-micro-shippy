package main

import (
  "context",
  "log",
  "net",
  "sync",

  // Import generated protobuf code
  pb "github.com/thededlier/go-micro-shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
  port = ":50051"
)

type repository interface {
  Create(*pb.Consignment) (*pb.Consignment, error)
}

// Dummy repo to represent some kind of datastore
type Repository struct {
  mu            sync.RWMutex
  consignments  []*pb.Consignment
}

// Creates a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment error) {
  repo.mu.Lock()
  updated := append(repo.consignments, consignment)
  repo.consignments = updated
  repo.mu.Unlock()
  return consignment, nil
}

// Service should implement all methods to satisfy the service defined in our protobuf definition
type service struct {
  repo repository
}

// Implemented on our service. Takes context and a request as an argument.
// These are handled by gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
  // Save our consignment
  consignment, err := s.repo.Create(req)
  if err != nil {
    return nil, err
  }

  return &pb.Response{ Created: true, Consignment: consignment }, nil
}

func main() {
  repo := &Repository{}

  // Set up gRPC server
  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  grpcServer := grpc.NewServer()

  // Register our consignment service with the gRPC server
  // This will tie our implementation into the auto-generated interface code for our
	// protobuf definition.
  pb.RegisterShippingServiceServer(grpcServer, &service{ repo })

  // Register reflection service on gRPC server.
  reflection.Register(grpcServer)

  log.Println("Running on port: ", port)
  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("Failed to serve: %v", err)
  }
}
