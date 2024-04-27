import grpc
import controller_pb2
import controller_pb2_grpc
from concurrent import futures
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler

# Define a global flag to control the HTTP server
http_server_running = False


class StringService(controller_pb2_grpc.StringServiceServicer):
    def SendString(self, request, context):
        global http_server_running
        print(f"Received request: {request.message}")
        if request.message == "start":
            # Start the HTTP server in a separate thread
            http_server_thread = threading.Thread(target=start_http_server)
            http_server_thread.start()
            http_server_running = True
            return controller_pb2.StringReply(message="HTTP server started")
        elif request.message == "stop":
            # Stop the HTTP server
            http_server_running = False
            return controller_pb2.StringReply(message="HTTP server stopped")
        else:
            return controller_pb2.StringReply(message="Unknown command")


def start_http_server():
    class MyHTTPRequestHandler(BaseHTTPRequestHandler):
        def do_GET(self):
            if self.path == "/m1":
                self.send_response(200)
                self.send_header("Content-type", "text/plain")
                self.end_headers()
                self.wfile.write("yeba".encode())
            else:
                self.send_response(404)
                self.end_headers()

    httpd = HTTPServer(("localhost", 8080), MyHTTPRequestHandler)
    print("HTTP server started on port 8080")
    while http_server_running:
        httpd.handle_request()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    controller_pb2_grpc.add_StringServiceServicer_to_server(StringService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
