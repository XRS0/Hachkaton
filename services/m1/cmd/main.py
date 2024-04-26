import grpc
import controller_pb2
import controller_pb2_grpc

def run():
    with grpc.insecure_channel('127.0.0.1:50051') as channel:
        stub = controller_pb2_grpc.StringServiceStub(channel)
        response = stub.SendString(controller_pb2.StringRequest(message="dasfsd"))
        print(f"Response from server: {response.message}")

if __name__ == '__main__':
    run()
