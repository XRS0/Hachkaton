package main

import (
	"fmt"
	pb "hack/services/SumGo/gen"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSumServiceServer
}

var port string = ":50051"

func main() {
	lis, err := net.Listen("tcp", port)
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
	success := fmt.Sprintf("SumService %s", port)
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

		num1 := *in.Num1
		num2 := *in.Num2
		fmt.Println(num1 + num2)
		res := strconv.Itoa(int(num1) + int(num2))
		success := fmt.Sprintf("HTTP Server is successfully started result is %s", res)
		response := &pb.SumResponse{Success: &success}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}
