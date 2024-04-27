package main

import (
	"context"
	"fmt"
	pb "hack/services/SumGo/gen"

	// . "hack/services/controller/internal/cli"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// cfg.ChangeJson(RunCLI())
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewSumServiceClient(conn)
	stream, err := client.Sum(context.Background())
	if err != nil {
		log.Fatalf("Ошибка при создании стрима: %v", err)
	}

	port := ":50052"

	go func() {
		request := &pb.SumRequest{Port: &port}
		if err := stream.Send(request); err != nil {
			log.Fatalf("Ошибка отправки: %v", err)
			return
		}
	}()

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Ошибка получения: %v", err)
			return
		}
		fmt.Println(*response.Success)
	}

}

// TODO:
// 1. Сделать мониторинг работы микросервисов-моков по эндпоинту /status, мы можем добавлять в обработку, мониторим карту запущенных микросервисов
// сервиса контроллера новые микросервисы ( выстроить работу с конфигом JSON ) тобешь добавляем в json файл конфигурации новый хост. ( микросервисы должны быть заранее созданы и запущены)
// 2. Сделать возможность запуска микросервисов через cli
// 3. Добавить поддержку стримов

// мы обрабатываем данные из стрима микросервиса и записываем их в json состояние сервиса
// мы делаем перенаправление команд из cli в нужный микросервис по его имени
// делаем uuid для каждого микросервиса(для обращения внутри программы)
