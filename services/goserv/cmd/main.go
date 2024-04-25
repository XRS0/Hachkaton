package main

import (
	"context"
	"log"
	"net"

	pb "hack/services/goserv/gen"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStringServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStringServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}

func (s *server) SendString(ctx context.Context, in *pb.StringRequest) (*pb.StringReply, error) {
	log.Printf("Received message: %s", in.GetMessage())
	return &pb.StringReply{Message: "Received: " + in.GetMessage()}, nil
}