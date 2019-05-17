# coding = utf-8

import json
from BaseHTTPServer import HTTPServer, BaseHTTPRequestHandler


class RequestHandler(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()

    def do_GET(self):
        response = {
            'status': 'SUCCESS',
            'data': 'hello from server'
        }

        self._set_headers()
        self.wfile.write(json.dumps(response))

    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        print 'post data from client:'
        print post_data

        response = {
            'status': 'SUCCESS',
            'data': 'server got your post data'
        }
        self._set_headers()
        self.wfile.write(json.dumps(response))


def run():
    port = 8000
    print('Listening on localhost:%s' % port)
    server = HTTPServer(('', port), RequestHandler)
    server.serve_forever()


run()
