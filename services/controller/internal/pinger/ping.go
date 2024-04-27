package pinger

import (
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
)

func PingBroadCast() {
	var arrayConn *grpc.ClientConn
	arrayConn = append(arrayConn, Ping("localhost:50051"))
}

func Ping(port string) *grpc.ClientConn {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		portInt, err := strconv.Atoi(port[9:])
		if err != nil {
			log.Fatal("Port is not string")
		}
		portInt++
		port = fmt.Sprintf("localhost:%v", portInt)
	}
	return conn
}
