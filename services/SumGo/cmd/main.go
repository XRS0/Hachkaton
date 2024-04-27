package main

import (
	"fmt"
	pb "hack/services/SumGo/gen"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSumServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSumServiceServer(s, &server{})
	log.Printf("Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func (s *server) Sum(stream pb.SumService_SumServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		port := *in.Port
		fmt.Println(port)
		success := "SumServer: HTTP Server is successfully started"
		response := &pb.SumResponse{Success: &success}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}
