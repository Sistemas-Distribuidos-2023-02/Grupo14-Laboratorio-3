package main

import (
    "log"
    "net"
	"os"
	"fmt"
	"strconv"
	"context"

    "google.golang.org/grpc"
    pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-3/proto"
)

type FulcrumServer struct {
	pb.UnimplementedFulcrumServer
    id int
    state map[string]map[string]int
    vClocks map[string][]int
}

func NewFulcrumServer(id int) *FulcrumServer {
    return &FulcrumServer{
        id: id,
        state: make(map[string]map[string]int),
        vClocks: make(map[string][]int),
    }
}

func (s *FulcrumServer) AgregarBase(sector string, base string, quantity int) {
    // If the quantity is 0 and the sector does not exist, return an error
    if quantity == 0 {
        if _, ok := s.state[sector]; !ok {
            log.Println("Cannot create a new sector with a quantity of 0")
            return
        }
    }

    // Update the state
    if _, ok := s.state[sector]; !ok {
        s.state[sector] = make(map[string]int)
    }
    s.state[sector][base] = quantity

    // Update the vector clock
    if _, ok := s.vClocks[sector]; !ok {
        s.vClocks[sector] = make([]int, 3)
    }
    s.vClocks[sector][s.id]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("AgregarBase %s %s %d\n", sector, base, quantity))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector)
}

func (s *FulcrumServer) RenombrarBase(sector string, base string, newBase string) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    s.state[sector][newBase] = s.state[sector][base]
    delete(s.state[sector], base)

    // Update the vector clock
    s.vClocks[sector][s.id]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("RenombrarBase %s %s %s\n", sector, base, newBase))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector)
}

func (s *FulcrumServer) ActualizarValor(sector string, base string, newValue int) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    s.state[sector][base] = newValue

    // Update the vector clock
    s.vClocks[sector][s.id]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("ActualizarValor %s %s %d\n", sector, base, newValue))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector)
}

func (s *FulcrumServer) BorrarBase(sector string, base string) {
    // If the sector or the base does not exist, return an error
    if _, ok := s.state[sector]; !ok {
        log.Println("Sector does not exist:", sector)
        return
    }
    if _, ok := s.state[sector][base]; !ok {
        log.Println("Base does not exist:", base)
        return
    }

    // Update the state
    delete(s.state[sector], base)

    // Update the vector clock
    s.vClocks[sector][s.id]++

    // Write to the log file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = f.WriteString(fmt.Sprintf("BorrarBase %s %s\n", sector, base))
    if err != nil {
        log.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        log.Println(err)
        return
    }

    // Update the sector file
    s.updateSectorFile(sector)
}

func (s *FulcrumServer) updateSectorFile(sector string) {
    // Open the sector file
    f, err := os.OpenFile(fmt.Sprintf("Sector%s.txt", sector), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
        return
    }
    defer f.Close()

    // Write to the sector file
    for base, quantity := range s.state[sector] {
        _, err = f.WriteString(fmt.Sprintf("%s %s %d\n", sector, base, quantity))
        if err != nil {
            log.Println(err)
            return
        }
    }

    // Write the vector clock
    _, err = f.WriteString(fmt.Sprintf("[%d,%d,%d]\n", s.vClocks[sector][0], s.vClocks[sector][1], s.vClocks[sector][2]))
    if err != nil {
        log.Println(err)
    }
}

func (s *FulcrumServer) ApplyCommand(ctx context.Context, command *pb.CommandRequest) (*pb.CommandResponse, error) {
    switch command.Action {
    case "AgregarBase":
        s.AgregarBase(command.Sector, command.Base, int(command.Value))
    case "ActualizarValor":
        s.ActualizarValor(command.Sector, command.Base, int(command.Value))
    case "RenombrarBase":
        s.RenombrarBase(command.Sector, command.Base, command.NewBase)
    case "BorrarBase":
        s.BorrarBase(command.Sector, command.Base)
    default:
        return nil, fmt.Errorf("unknown action: %s", command.Action)
    }

    vectorClock := make([]int32, len(s.vClocks[command.Sector]))
	for i, v := range s.vClocks[command.Sector] {
		vectorClock[i] = int32(v)
	}

	// Return the vector clock of the modified sector
	return &pb.CommandResponse{
		VectorClock: vectorClock,
	}, nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run main.go <server_id>")
        os.Exit(1)
    }

    id, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Invalid server ID:", os.Args[1])
        os.Exit(1)
    }

    // Initialize the server
    s := NewFulcrumServer(id)

    // Start a gRPC server
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051+id))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterFulcrumServer(grpcServer, s)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}