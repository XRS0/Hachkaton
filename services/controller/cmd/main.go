package main

import (
	. "hack/services/controller/internal/cli"
	cfg "hack/services/controller/pkg/handleconfig"
)

// "context"
// "log"

// pb "hack/services/controller/gen"

// "google.golang.org/grpc"

func main() {
	cfg.ChangeJson(RunCLI())

	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to connection: %v", err)
	// }
	// defer conn.Close()
	// client := pb.NewSumServiceClient(conn)
	// stream, err := client.Sum(context.Background())
	// if err != nil {
	// 	log.Fatalf("%v.Sum(_) = _, %v", client, err)
	// }

}

// TODO:
// 1. Сделать мониторинг работы микросервисов-моков по эндпоинту /status, мы можем добавлять в обработку, мониторим карту запущенных микросервисов ( id:host )
// сервиса контроллера новые микросервисы ( выстроить работу с конфигом JSON ) тобешь добавляем в json файл конфигурации новый порт или хост. ( микросервисы должны быть заранее созданы и запущены)
// 2. Сделать возможность запуска микросервисов через cli
// 3. Добавить поддержку стримов
