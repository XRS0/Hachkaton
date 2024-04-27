gen:
	cd api && protoc --go_out=../services/controller/gen --go-grpc_out=../services/controller/gen m1.proto && python3 -m grpc_tools.protoc -I. --python_out=../services/m1/cmd/ --grpc_python_out=../services/m1/cmd/ m1.proto

goruncli:
	cd services/controller/cmd/ && go build -o main.go && ./cmd start m1 50051

pyrun:
	python3 services/sumService/cmd/main.py

gobuild:
	cd services/controller/cmd && go build -o controller . && sudo mv controller /usr/local/bin/

gentest:
	cd api && protoc --go_out=../services/controller/gen --go-grpc_out=../services/controller/gen sumgo.proto && protoc --go_out=../services/SumGo/gen --go-grpc_out=../services/SumGo/gen sumgo.proto
