package main

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "log"
  "context"

  pb "github.com/thededlier/go-micro-shippy/consignment-service/proto/consignment"
  micro "github.com/micro/go-micro"

)

const (
  address         = "localhost:50051"
  sampleFile      = "sample-consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
  var consignment *pb.Consignment
  data, err := ioutil.ReadFile(file)

  if err != nil {
    return nil, err
  }
  json.Unmarshal(data, &consignment)
  return consignment, err
}

func main() {
  // Set up connection using go micro service
  service := micro.NewService(
    micro.Name("shippy.consignment.cli"),
  )
  service.Init()

  client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())

  // Setup file as the default sample file. If cli args are given for another file, use that
  file := sampleFile
  if len(os.Args) > 1 {
    file = os.Args[1]
  }

  consignment, err := parseFile(file)
  if err != nil {
    log.Fatalf("Could not parse file: %v", err)
  }

  // Create new consignment
  createResponse, err := client.CreateConsignment(context.Background(), consignment)
  if err != nil {
    log.Fatalf("Failed to create: %v", err)
  }
  log.Printf("Created: %t", createResponse.Created)

  getAllResponse, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
  if err != nil {
    log.Fatalf("Failed to list all consignments: %v", err)
  }
  for _, consignment := range getAllResponse.Consignments {
    log.Println(consignment)
  }
}
