import grpc
import controller_pb2
import controller_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = controller_pb2_grpc.StringServiceStub(channel)
        response = stub.SendString(controller_pb2.StringRequest(message="hello Жопа"))
        print(f"Greeting: {response.message}")

if __name__ == '__main__':
    run()
