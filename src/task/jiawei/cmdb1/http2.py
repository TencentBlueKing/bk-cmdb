# coding:utf-8

import socket
import threading
import json


def handle_client(client_socket):
    """
    处理客户端请求
    """
    request_data = client_socket.recv(2048)
    body = request_data.split('\r\n\r\n')[-1]
    print request_data
    print "request data:", json.loads(body)
    # 构造响应数据
    response_start_line = "HTTP/1.1 200 OK\r\n"
    response_headers = "Server: My server\r\n"
    response_body = "<h1>Python HTTP Test</h1>"
    response = response_start_line + response_headers + "\r\n" + response_body
    # 向客户端返回响应数据
    client_socket.send(bytes(response))
    # 关闭客户端连接
    client_socket.close()


if __name__ == "__main__":
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(("", 8045))
    server_socket.listen(128)

    while True:
        client_socket, client_address = server_socket.accept()
        print "[%s, %s]user connected" % client_address
        t = threading.Thread(target=handle_client, args=(client_socket,))
        t.start()
    client_socket.close()
