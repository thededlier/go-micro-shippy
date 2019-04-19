package main

import (
  "fmt"
  "context"

  // Import generated protobuf code
  pb "github.com/thededlier/go-micro-shippy/consignment-service/proto/consignment"
  vesselProto "github.com/thededlier/go-micro-shippy/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

type repository interface {
  Create(*pb.Consignment) (*pb.Consignment, error)
  GetAll() []*pb.Consignment
}

// Dummy repo to represent some kind of datastore
type Repository struct {
  consignments  []*pb.Consignment
}

// Creates a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
  updated := append(repo.consignments, consignment)
  repo.consignments = updated
  return consignment, nil
}

// Get all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
  return repo.consignments
}

// Service should implement all methods to satisfy the service defined in our protobuf definition
type service struct {
  repo repository
  vesselClient vesselProto.VesselServiceClient
}

// Implemented on our service. Takes context and a request as an argument.
// These are handled by gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

  // Determine vessel
  vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification {
    MaxWeight: req.Weight,
    Capacity: int32(len(req.Containers)),
  })

  fmt.Printf("Vessel found: %v \n", vesselResponse.Vessel.Name)
  if err != nil {
    return err
  }
  // Set vessel id
  req.VesselId = vesselResponse.Vessel.Id

  // Save our consignment
  consignment, err := s.repo.Create(req)
  if err != nil {
    return err
  }
  res.Created = true
  res.Consignment = consignment
  return nil
}

// Get consignments. Implemented on service
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
  consignments := s.repo.GetAll()
  res.Consignments = consignments
  return nil
}

func main() {
  repo := &Repository{}

  // Create a new service
  srv := micro.NewService(
    // Name must match the same as package name in protobuf definition
    micro.Name("go.micro.srv.consignment"),
    micro.Version("latest"),
  )

  srv.Init()

  // Setup Vessel Client
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

  // Register handler
  pb.RegisterShippingServiceHandler(srv.Server(), &service{ repo, vesselClient })

  // Run the server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}
