package main

import (
    "context"
    "log"
    "net"

    pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-3/proto"

    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedVanguardServer
    brokerClient pb.BrokerClient
}

func (s *server) ExecuteCommand(ctx context.Context, in *pb.Command) (*pb.Response, error) {
    // Forward the command to the Broker server
    message := &pb.Message{Command: in.GetCommand()}
    ack, err := s.brokerClient.Mediate(ctx, message)
    if err != nil {
        return nil, err
    }

    // Return the response from the Broker server to the user
    return &pb.Response{Acknowledgement: ack.GetAcknowledgement()}, nil
}

func main() {
    // Create a listener on TCP port 50051 (or any port you want)
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Create a gRPC server object
    grpcServer := grpc.NewServer()

    // Connect to the Broker server
    conn, err := grpc.Dial("broker_address", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to Broker server: %v", err)
    }
    defer conn.Close()

    // Create a Broker client
    brokerClient := pb.NewBrokerClient(conn)

    // Create a new Vanguard server
    vanguardServer := &server{brokerClient: brokerClient}

    // Attach the Vanguard service to the gRPC server
    pb.RegisterVanguardServer(grpcServer, vanguardServer)

    // Start the gRPC server (blocking)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}