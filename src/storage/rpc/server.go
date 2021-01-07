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
	"context"
	"io"
	"net/http"
	"runtime/debug"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

// Server define
type Server struct {
	ctx            context.Context
	codec          Codec
	handlers       map[string]HandlerFunc
	streamHandlers map[string]HandlerStreamFunc
}

// NewServer returns new server
func NewServer() *Server {
	return &Server{
		ctx:            context.Background(),
		codec:          JSONCodec,
		handlers:       map[string]HandlerFunc{},
		streamHandlers: map[string]HandlerStreamFunc{},
	}
}

var connected = "200 Connected to CC RPC"
var connectfaile = "400 Connect failed to CC RPC"

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

	session, err := NewServerSession(s, conn, "deflate")
	if err != nil {
		blog.Errorf("rpc new server session faile %s: %s", req.RemoteAddr, err.Error())
		io.WriteString(conn, "HTTP/1.0 "+connectfaile+"\n\n")
		return
	}

	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	blog.V(3).Infof("connect from rpc client %s", req.RemoteAddr)

	if err = session.Run(); err != nil {
		blog.Errorf("dissconnect from rpc client %s: %s ", req.RemoteAddr, err.Error())
		return
	}
}

// Handle regist new handler
func (s *Server) Handle(name string, f HandlerFunc) {
	s.handlers[name] = f
}

// HandleStream regist new stream handler
func (s *Server) HandleStream(name string, f HandlerStreamFunc) {
	s.streamHandlers[name] = f
}

// SetCodec set server codec
func (s *Server) SetCodec(codec Codec) {
	s.codec = codec
}

// ServerSession define a server session
type ServerSession struct {
	srv       *Server
	wire      Wire
	request   Message
	responses chan *Message
	done      *util.AtomicBool
	stream    *streamstore
}

// NewServerSession returns a new ServerSession
func NewServerSession(srv *Server, conn io.ReadWriteCloser, compress string) (*ServerSession, error) {
	wire, err := NewBinaryWire(conn, compress)
	if err != nil {
		return nil, err
	}
	return &ServerSession{
		srv:       srv,
		wire:      wire,
		responses: make(chan *Message, 1024),
		done:      util.NewBool(false),
		stream:    newStreamStore(),
		request:   Message{codec: srv.codec},
	}, nil
}

// Run run the Serssion
func (s *ServerSession) Run() error {
	go s.writeloop()
	return s.readloop()
}

// Stop stop the server session
func (s *ServerSession) Stop() {
	s.done.Set()
}

func (s *ServerSession) readFromWire() error {
	msg := Message{codec: s.srv.codec}
	err := s.wire.Read(&msg)
	if err == io.EOF {
		return err
	} else if err != nil {
		blog.Errorf("Failed to read: %v", err)
		return err
	}

	switch msg.typz {
	case TypeRequest:
		blog.V(3).Infof("[rpc server] calling [%s]", msg.cmd)
		if handlerFunc, ok := s.srv.handlers[msg.cmd]; ok {
			go s.handle(handlerFunc, &msg)
		} else if handlerFunc, ok := s.srv.streamHandlers[msg.cmd]; ok {
			go s.handleStream(handlerFunc, &msg)
		} else {
			cmds := []string{}
			for cmd := range s.srv.handlers {
				cmds = append(cmds, cmd)
			}
			blog.V(3).Infof("[rpc server] command [%s] not found, existing command are: %#v", msg.cmd, s.srv.handlers)
			s.pushResponse(&msg, ErrCommandNotFount)
		}
	case TypeStream, TypeStreamClose:
		s.stream.RLock()
		stream, ok := s.stream.get(msg.seq)
		if ok {
			stream.input <- &msg
		}
		s.stream.RUnlock()
	case TypePing:
		go s.handlePing(&msg)
	default:
		blog.Warnf("[rpc server] unknow message type: %v", msg.typz)
	}
	return nil
}

func (s *ServerSession) handle(f HandlerFunc, msg *Message) {
	defer func() {
		runtimeErr := recover()
		if runtimeErr != nil {
			stack := debug.Stack()
			blog.Errorf("command [%s] failed: %v\n%s", msg.cmd, runtimeErr, stack)
		}
	}()
	result, err := f(msg)
	if encodeErr := msg.Encode(result); encodeErr != nil {
		blog.Errorf("EncodeData error: %s", encodeErr.Error())
	}
	s.pushResponse(msg, err)
}
func (s *ServerSession) handleStream(f HandlerStreamFunc, msg *Message) {
	stream := NewStreamMessage(msg)
	s.stream.store(msg.seq, stream)
	s.pushResponse(msg, nil)

	go func() {
		defer func() {
			runtimeErr := recover()
			if runtimeErr != nil {
				stack := debug.Stack()
				blog.Errorf("stream command [%s] failed: %v\n%s", msg.cmd, runtimeErr, stack)
			}
		}()
		err := f(msg, stream)
		nmsg := msg.copy()
		nmsg.typz = TypeStreamClose
		if err != nil {
			nmsg.Data = []byte(err.Error())
		}
		stream.output <- nmsg
	}()
	go func() {
	serverstreamloop:
		for smsg := range stream.output {
			s.responses <- smsg
			if smsg.typz == TypeStreamClose {
				break serverstreamloop
			}
			if s.done.IsSet() || stream.done.IsSet() {
				stream.err = ErrStreamStoped
				break
			}
		}
		s.stream.remove(msg.seq)
		close(stream.input)
		close(stream.output)
	}()
}

func (s *ServerSession) handlePing(msg *Message) {
	s.pushResponse(msg, nil)
}

func (s *ServerSession) pushResponse(msg *Message, err error) {
	msg.magicVersion = MagicVersion

	msg.typz = TypeResponse
	if err != nil {
		msg.typz = TypeError
		msg.Data = []byte(err.Error())
	}
	s.responses <- msg
}

func (s *ServerSession) readloop() error {
	var err error
	for {
		if s.done.IsSet() {
			blog.Infof("[rpc server] RPC server stopped")
			return nil
		}
		if err = s.readFromWire(); err != nil {
			s.Stop()
			return nil
		}
	}
}

func (s *ServerSession) writeloop() {
	for msg := range s.responses {
		if err := s.wire.Write(msg); err != nil {
			blog.Errorf("Failed to write: %v", err)
		}
		if s.done.IsSet() {
			if queuelen := len(s.responses); queuelen > 0 {
				for queuelen > 0 {
					msg := <-s.responses
					if err := s.wire.Write(msg); err != nil {
						blog.Errorf("Failed to write: %v", err)
						break
					}
				}
			}
			msg := &Message{
				typz: TypeClose,
			}
			//Best effort to notify client to close connection
			s.wire.Write(msg)
			break
		}
	}
}
