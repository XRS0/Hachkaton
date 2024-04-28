import grpc
import service_pb2
import service_pb2_grpc
from concurrent import futures
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler

http_server_running = False
now_port = None
is_port_change = False


class SumService(service_pb2_grpc.SumServiceServicer):
    def Sum(self, request_iterator, context):
        global http_server_running
        for request in request_iterator:
            if request.action == 'sum':
                print(request.num1, request.num2, request.port)
                result = request.num1 + request.num2
                yield service_pb2.SumResponse(data=str(result), success='true')
            elif request.action == 'start':
                http_server_thread = threading.Thread(target=start_http_server)
                http_server_thread.start()
                http_server_running = True
                yield service_pb2.SumResponse(data='HTTP сервер запущен', success='true')
            elif request.action == "stop":
                http_server_running = False
                yield service_pb2.SumResponse(data='HTTP сервер остановлен', success='true')
            elif request.action == "chport":
                change_port(request.port)


def change_port(port):
    global is_port_change
    if port != now_port:
        server.stop(grace=None)
        is_port_change = True
        serve(port)
    return


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
