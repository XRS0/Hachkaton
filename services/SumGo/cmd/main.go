package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	pb "hack/services/SumGo/gen"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSumServiceServer
	lis    net.Listener
	server *grpc.Server
	port   string
	mu     sync.Mutex
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{
		lis:    lis,
		server: s,
		port:   port,
	}
	pb.RegisterSumServiceServer(s, srv)
	log.Printf("Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Sum(stream pb.SumService_SumServer) error {
	success := fmt.Sprintf("SumService %s", s.port)
	response := &pb.SumResponse{Data: &success}
	if err := stream.Send(response); err != nil {
		return err
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		num1 := in.Num1
		num2 := in.Num2
		fmt.Println(num1 + num2)
		res := strconv.Itoa(int(num1) + int(num2))
		success := fmt.Sprintf("HTTP Server is successfully started result is %s", res)
		response := &pb.SumResponse{Success: success}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func (s *server) DoAction(ctx context.Context, in *pb.DoActionRequest) (*pb.DoActionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.Action == "chport" {
		s.server.Stop()

		newPort := ":" + *in.Port
		lis, err := net.Listen("tcp", newPort)
		if err != nil {
			log.Printf("Failed to listen on new port: %v", err)
			return nil, err
		}

		s.server = grpc.NewServer()
		pb.RegisterSumServiceServer(s.server, s)
		s.lis = lis
		s.port = newPort

		go func() {
			if err := s.server.Serve(lis); err != nil {
				log.Fatalf("Failed to serve on new port: %v", err)
			}
		}()
		log.Printf("Server is restarted and listening at %v", lis.Addr())
		return &pb.DoActionResponse{Success: "Server restarted on new port " + *in.Port}, nil
	}

	return &pb.DoActionResponse{Success: "OK"}, nil
}
