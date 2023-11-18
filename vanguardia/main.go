package main

import (
    "context"
    "log"
    "net"
    "fmt"

    pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-3/proto"

    "google.golang.org/grpc"
)

type LogEntry struct {
    SectorInfo   string
    VectorClock  []int32
    FulcrumServer string
}

var logEntries []LogEntry

type server struct {
    pb.UnimplementedVanguardServer
    brokerClient pb.BrokerClient
    clientClocks map[string][]int32
}

func (s *server) GetSoldados(ctx context.Context, in *pb.Command) (*pb.Response, error) {
    // Get the client's latest vector clock
    clientClock, ok := s.clientClocks[in.ClientId]
    if !ok {
        clientClock = make([]int32, len(s.clientClocks))
    }

    // Forward the command to the Broker server
    message := &pb.Message{
        Sector: in.GetSector(),
        Base: in.GetBase(),
        VectorClock: clientClock,
    }
    ack, err := s.brokerClient.Mediate(ctx, message)
    if err != nil {
        return nil, err
    }

    // Update the client's vector clock
    s.clientClocks[in.ClientId] = ack.GetVectorClock()

    // Log the command and response
    logEntry := LogEntry{
        SectorInfo:   fmt.Sprintf("GetSoldados %s %s", in.GetSector(), in.GetBase()),
        VectorClock:  ack.GetVectorClock(),
        FulcrumServer: ack.GetFulcrumServer(),
    }
    logEntries = append(logEntries, logEntry)

    // Return the response from the Broker server to the user
    return &pb.Response{
        Acknowledgement: ack.GetAcknowledgement(),
        FulcrumServer: ack.GetFulcrumServer(),
        VectorClock: ack.GetVectorClock(),
    }, nil
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
    fmt.Println("Vanguard server running...")

    // Attach the Vanguard service to the gRPC server
    pb.RegisterVanguardServer(grpcServer, vanguardServer)

    go func() { // Consola input
        for {
            fmt.Print("Enter sector (ingresar solo el nombre, primera letra mayuscula): ")
            var sector string
            fmt.Scanln(&sector)
    
            fmt.Print("Enter base (ingresar solo el nombre, primera letra mayuscula): ")
            var base string
            fmt.Scanln(&base)
    
            // Create a Command message
            cmd := &pb.Command{Sector: sector, Base: base}
    
            // Call the GetSoldados method
            res, err := vanguardServer.GetSoldados(context.Background(), cmd)
            if err != nil {
                log.Fatalf("Failed to execute command: %v", err)
            }
    
            // Print the response
            fmt.Println("Response:", res.GetAcknowledgement())
            fmt.Println("Fulcrum Server:", res.GetFulcrumServer())
            fmt.Println("Vector Clock:", res.GetVectorClock())
        }
    }()

    // Start the gRPC server (blocking)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}