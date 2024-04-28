package pinger

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "hack/services/controller/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionInfo struct {
	Port   string
	Stream pb.SumService_SumClient
	Conn   *grpc.ClientConn
}

type StreamManager struct {
	Connections []*ConnectionInfo
	Mutex       sync.RWMutex
}

func (m *StreamManager) AddConnection(info *ConnectionInfo) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Connections = append(m.Connections, info)
}

func (m *StreamManager) RemoveConnection(port string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	for i, conn := range m.Connections {
		if conn.Port == port {
			m.Connections = append(m.Connections[:i], m.Connections[i+1:]...)
			break
		}
	}
}

var manager = StreamManager{}

func ConnectToStream(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if err := AttemptStream(port); err != nil {
			fmt.Printf("Stream on port %s failed: %v\n", port, err)
			time.Sleep(5 * time.Second)
			fmt.Printf("Attempting to reconnect to port %s\n", port)
			continue
		}
		break
	}
}

func AttemptStream(port string) error {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewSumServiceClient(conn)
	stream, err := client.Sum(context.Background())
	if err != nil {
		return err
	}

	manager.AddConnection(&ConnectionInfo{
		Port:   port,
		Stream: stream,
		Conn:   conn,
	})

	var num1, num2 int64

	num1, num2 = -10, 2

	go SendData(num1, num2, stream, port)

	response, err := stream.Recv()
	if err != nil {
		manager.RemoveConnection(port)
		return err
	}
	respData := *response.Data
	for {
		response, err := stream.Recv()
		if err != nil {
			manager.RemoveConnection(port)
			return err
		}
		fmt.Printf("INFO Received from %s, success: %s\n", respData, response.Success)
	}
}

func SendData(num1, num2 int64, stream pb.SumService_SumClient, port string) {
	if err := stream.Send(&pb.SumRequest{Num1: num1, Num2: num2}); err != nil {
		fmt.Printf("Failed to send data to port %s: %v\n", port, err)
		return
	}
	time.Sleep(1 * time.Second)
	num1++
	num2++
	SendData(num1, num2, stream, port)
}

func DoAction(oldPort, newPort string) {
	hst := fmt.Sprintf("localhost:%s", oldPort)
	conn, err := grpc.Dial(hst, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewSumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DoAction(ctx, &pb.DoActionRequest{
		Action: "chport",
		Port:   &newPort,
	})
	if err != nil {
		log.Fatalf("could not execute action: %v", err)
	}
	fmt.Printf("DoAction response: %s\n", r.GetData())
	conn.Close()
	hst = fmt.Sprintf("localhost:%s", newPort)
	conn, err = grpc.Dial(hst, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = pb.NewSumServiceClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
}
