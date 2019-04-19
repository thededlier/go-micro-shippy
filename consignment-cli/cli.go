package main

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "log"
  "context"

  pb "github.com/thededlier/go-micro-shippy/consignment-service/proto/consignment"
  "google.golang.org/grpc"
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
  // Set up connection to grpc
  // TODO: Will add proper auth later
  conn, err := grpc.Dial(address, grpc.WithInsecure())

  if err != nil {
    log.Fatalf("Failed to connect : %v", err)
  }
  defer conn.Close()
  client := pb.NewShippingServiceClient(conn)
  // Setup file as the default sample file. If cli args are given for another file, use that
  file := sampleFile
  if len(os.Args) > 1 {
    file = os.Args[1]
  }

  consignment, err := parseFile(file)
  if err != nil {
    log.Fatalf("Could not parse file: %v", err)
  }

  response, err := client.CreateConsignment(context.Background(), consignment)
  if err != nil {
    log.Fatalf("Failed to create: %v", err)
  }
  log.Printf("Created: %t", response.Created)
}
