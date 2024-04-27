import grpc
import controller_pb2
import controller_pb2_grpc
from concurrent import futures
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler

http_server_running = False


class StringService(controller_pb2_grpc.StringServiceServicer):
    def StreamStrings(self, request_iterator, context):
        global http_server_running
        # Send a stream message to indicate that the server is running
        yield controller_pb2.StringReply(message="Сервер работает")
        for request in request_iterator:
            message = request.message
            message_parts = message.split(", ")

            # Initialize variables for port and action
            port = None
            action = None

            # Iterate over each part of the message
            for part in message_parts:
                key, value = part.split(":")
                if key == "port":
                    port = int(value)
                    print(port)
                    if now_port != port:
                        server.stop(grace=None)
                        serve(port)
                elif key == "action":
                    action = value
            print(f"принято с клиента: {request.message}")
            yield controller_pb2.StringReply(message=f"стрим от сервера")
            if action == "start":
                http_server_thread = threading.Thread(target=start_http_server)
                http_server_thread.start()
                http_server_running = True
                return controller_pb2.StringReply(message="HTTP server started")
            elif action == "stop":
                http_server_running = False
                return controller_pb2.StringReply(message="HTTP server stopped")


def start_http_server():
    print('HTTP was started')

    class MyHTTPRequestHandler(BaseHTTPRequestHandler):
        def do_GET(self):
            if self.path == "/m1":
                print('/m1')
                self.send_response(200)
                self.send_header("Content-type", "text/plain")
                self.end_headers()
                self.wfile.write("is working".encode())
            else:
                self.send_response(404)
                self.end_headers()

    httpd = HTTPServer(("localhost", 8080), MyHTTPRequestHandler)
    print("HTTP server started on port 8080")
    while http_server_running:
        httpd.handle_request()


def serve(grpc_port):
    global server, now_port
    now_port = grpc_port
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    controller_pb2_grpc.add_StringServiceServicer_to_server(StringService(), server)
    server.add_insecure_port(f"[::]:{grpc_port}")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve(50051)
