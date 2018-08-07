/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rpc

import (
	"configcenter/src/common/blog"
	"io"
	"net/http"
)

type Server struct {
	handlers       map[string]HandlerFunc
	streamHandlers map[string]HandlerStreamFunc
}

func NewServer() *Server {
	return &Server{
		handlers:       map[string]HandlerFunc{},
		streamHandlers: map[string]HandlerStreamFunc{},
	}
}

var connected = "200 Connected to CC RPC"

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		resp.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(resp, "405 must CONNECT\n")
		return
	}
	conn, _, err := resp.(http.Hijacker).Hijack()
	if err != nil {
		blog.Errorf("rpc hijacking %s: %s", req.RemoteAddr, err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")

	session := NewServerSession(s, conn)
	if err = session.Run(); err != nil {
		blog.Errorf("Run ServerSession error:  %s: %s ", req.RemoteAddr, err.Error())
		return
	}
}

func (s *Server) Handle(name string, f HandlerFunc) {
	s.handlers[name] = f
}
func (s *Server) HandleStream(name string, f HandlerStreamFunc) {
	s.streamHandlers[name] = f
}

type ServerSession struct {
	srv       *Server
	wire      Wire
	responses chan *Message
	done      chan struct{}
}

func NewServerSession(srv *Server, conn io.ReadWriteCloser) *ServerSession {
	return &ServerSession{
		srv:       srv,
		wire:      NewBinaryWire(conn),
		responses: make(chan *Message, 1024),
		done:      make(chan struct{}, 5),
	}
}

func (s *ServerSession) Run() error {
	go s.writeloop()
	defer func() {
		s.done <- struct{}{}
	}()
	return s.readloop()
}

func (s *ServerSession) Stop() {
	s.done <- struct{}{}
}

func (s *ServerSession) readFromWire(ret chan<- error) {
	msg, err := s.wire.Read()
	if err == io.EOF {
		ret <- err
		return
	} else if err != nil {
		blog.Errorf("Failed to read: %v", err)
		ret <- err
		return
	}

	switch msg.typz {
	case TypeRequest:
		blog.Infof("calling %s", msg.cmd.String())
		if handlerFunc, ok := s.srv.handlers[msg.cmd.String()]; ok {
			go s.handle(handlerFunc, msg)
		} else if handlerFunc, ok := s.srv.streamHandlers[msg.cmd.String()]; ok {
			go s.handleStream(handlerFunc, msg)
		} else {
			s.pushResponse(0, msg, ErrCommandNotFount)
		}

	case TypePing:
		go s.handlePing(msg)
	}
	ret <- nil
}

func (s *ServerSession) handle(f HandlerFunc, msg *Message) {
	result, err := f(msg)
	if encodeErr := msg.Encode(result); encodeErr != nil {
		blog.Errorf("EncodeData error: %s", encodeErr.Error())
	}
	s.pushResponse(0, msg, err)
}
func (s *ServerSession) handleStream(f HandlerStreamFunc, msg *Message) {
	// TODO handle stream
	ch := make(chan *Message)
	err := f(ch)
	s.pushResponse(0, msg, err)
}

func (s *ServerSession) handlePing(msg *Message) {
	s.pushResponse(0, msg, nil)
}

func (s *ServerSession) pushResponse(count int, msg *Message, err error) {
	msg.magicVersion = MagicVersion
	msg.Size = uint32(len(msg.Data))

	msg.typz = TypeResponse
	if err != nil {
		msg.typz = TypeError
		msg.Data = []byte(err.Error())
		msg.Size = uint32(len(msg.Data))
	}
	s.responses <- msg
}

func (s *ServerSession) readloop() error {
	ret := make(chan error)
	for {
		go s.readFromWire(ret)

		select {
		case err := <-ret:
			if err != nil {
				return err
			}
			continue
		case <-s.done:
			blog.Infof("RPC server stopped")
			return nil
		}
	}
}

func (s *ServerSession) writeloop() {
	for {
		select {
		case msg := <-s.responses:
			if err := s.wire.Write(msg); err != nil {
				blog.Errorf("Failed to write: %v", err)
			}
		case <-s.done:
			msg := &Message{
				typz: TypeClose,
			}
			//Best effort to notify client to close connection
			s.wire.Write(msg)
			break
		}
	}
}
