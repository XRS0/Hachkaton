gen:
	cd api && protoc --go_out=../services/goserv/gen --go-grpc_out=../services/goserv/gen controller.proto && python3 -m grpc_tools.protoc -I. --python_out=../services/pyserv/cmd/ --grpc_python_out=../services/pyserv/cmd/ controller.proto

gorun:
	go run services/goserv/cmd/main.go

pyrun:
	python3 services/pyserv/cmd/main.py
