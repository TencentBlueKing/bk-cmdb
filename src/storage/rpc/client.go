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
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"configcenter/src/common/blog"
)

var (
	opRetries     = 0
	opReadTimeout = 15 * time.Second // client read
	opPingTimeout = 20 * time.Second
)

const commanLimit = 40

//Client replica client
type Client struct {
	end       chan struct{}
	requests  chan *Message
	send      chan *Message
	responses chan *Message
	seq       uint32
	messages  map[uint32]*Message
	stream    *streamstore
	wire      Wire
	peerAddr  string
	err       error
	codec     Codec
}

//NewClient replica client
func NewClient(conn net.Conn, compress string) (*Client, error) {
	wire, err := NewBinaryWire(conn, compress)
	if err != nil {
		return nil, err
	}
	c := &Client{
		wire:      wire,
		peerAddr:  conn.RemoteAddr().String(),
		end:       make(chan struct{}, 1024),
		requests:  make(chan *Message, 1024),
		send:      make(chan *Message, 1024),
		responses: make(chan *Message, 1024),
		messages:  map[uint32]*Message{},
		codec:     JSONCodec,
		stream:    newStreamStore(),
	}
	blog.V(3).Infof("connected to %s", c.TargetID())
	go c.loop()
	go c.write()
	go c.read()
	return c, nil
}

func Dial(connect string) (*Client, error) {
	uri, err := url.Parse(connect)
	if err != nil {
		return nil, err
	}
	var port = uri.Port()
	if uri.Port() == "" {
		port = "80"
	}
	return DialHTTPPath("tcp", uri.Hostname()+":"+port, uri.Path)
}

// DialHTTPPath connects to an HTTP RPC server
// at the specified network address and path.
func DialHTTPPath(network, address, path string) (*Client, error) {
	var err error
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.Status == connected {
		return NewClient(conn, "deflate")
	}
	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}
	conn.Close()
	return nil, &net.OpError{
		Op:   "dial-http",
		Net:  network + " " + address,
		Addr: nil,
		Err:  err,
	}
}

//Close replica client
func (c *Client) Close() error {
	c.wire.Close()
	c.end <- struct{}{}
	return nil
}

// TargetID operation target ID
func (c *Client) TargetID() string {
	return c.peerAddr
}

// Call replica client
func (c *Client) Call(cmd string, input interface{}, result interface{}) error {
	msg, err := c.operation(TypeRequest, cmd, input)
	if err != nil {
		return err
	}
	return msg.Decode(result)
}

// CallStream replica client
func (c *Client) CallStream(cmd string, input interface{}) (*StreamMessage, error) {
	msg, err := c.operation(TypeRequest, cmd, input)
	if err != nil {
		return nil, err
	}

	sm := NewStreamMessage(msg)
	c.stream.store(msg.seq, sm)
	go func() {
	clientstreamloop:
		for {
			select {
			case streammsg := <-sm.output:
				c.send <- streammsg
				if msg.typz == TypeStreamClose {
					break clientstreamloop
				}
			case <-sm.done:
				break clientstreamloop
			case <-c.end:
				break clientstreamloop
			}
		}
		sm.err = ErrStreamStoped
		close(sm.input)
		close(sm.output)
		close(sm.done)
	}()

	return sm, nil
}

//Ping replica client
func (c *Client) Ping() error {
	_, err := c.operation(TypePing, "", nil)
	return err
}

func (c *Client) operation(op MessageType, cmd string, data interface{}) (*Message, error) {
	retry := 0
	for {
		msg := Message{
			magicVersion: MagicVersion,
			codec:        c.codec,
			seq:          c.nextSeq(),
			complete:     make(chan struct{}, 1),
			typz:         op,
			cmd:          cmd,
			Data:         nil,
		}

		if op == TypeRequest {
			err := msg.Encode(data)
			if err != nil {
				return nil, err
			}
		}

		timeout := func(op MessageType) <-chan time.Time {
			switch op {
			case TypeRequest:
				return time.After(opReadTimeout)
			}
			return time.After(opPingTimeout)
		}(msg.typz)

		c.requests <- &msg

		select {
		case <-msg.complete:
			if msg.typz == TypeError {
				return nil, errors.New(string(msg.Data))
			}
			return &msg, nil
		case <-timeout:
			switch msg.typz {
			case TypeRequest:
				blog.Errorf("request timeout on replcia %s, seq= %d", c.TargetID(), msg.seq)
			case TypePing:
				blog.Errorf("Ping timeout on replica %s, seq= %d", c.TargetID(), msg.seq)
			}
			if retry < opRetries {
				retry++
				blog.Errorf("Retry %d on replica, seq= %d", retry, msg.seq)
			} else {
				err := ErrRWTimeout
				if msg.typz == TypePing {
					err = ErrPingTimeout
				}
				// not transportErr when timeout?
				// c.responses <- &Message{
				// 	transportErr: err,
				// }
				// TODO print journal
				return nil, err
			}
		}
	}
}

func (c *Client) loop() {
	defer close(c.send)

	for {
		select {
		case <-c.end:
			return
		case req := <-c.requests:
			c.handleRequest(req)
		case resp := <-c.responses:
			c.handleResponse(resp)
		}
	}
}

func (c *Client) nextSeq() uint32 {
	c.seq++
	return c.seq
}

func (c *Client) replyError(req *Message) {
	delete(c.messages, req.seq) // no need lock cause the loop is serial
	req.typz = TypeError
	req.Data = []byte(c.err.Error())
	req.complete <- struct{}{}
}

func (c *Client) handleRequest(req *Message) {
	// TODO req.ID = journal.InsertPendingOp(time.Now(), c.TargetID(), journal.OpPing, 0)
	if c.err != nil {
		c.replyError(req)
		return
	}

	c.messages[req.seq] = req
	c.send <- req
	blog.V(5).Infof("[rpc client]sent message data: %s", req.Data)
}

func (c *Client) handleResponse(resp *Message) {
	blog.V(5).Infof("[rpc client]receive message data: %s", resp.Data)
	resp.codec = c.codec
	if resp.transportErr != nil {
		c.err = resp.transportErr
		// Terminate all in flight
		for _, msg := range c.messages {
			c.replyError(msg)
		}
		// TODO reconnect when transportErr?
		return
	}
	if resp.typz == TypeStream || resp.typz == TypeStreamClose {
		stream, ok := c.stream.get(resp.seq)
		if ok {
			stream.input <- resp
		} else {
			blog.Warnf("[rpc client] stream not found, resp is %s", resp.Data)
		}
		return
	}
	if req, ok := c.messages[resp.seq]; ok {
		if c.err != nil {
			c.replyError(req)
			return
		}

		delete(c.messages, resp.seq) // no need lock cause the loop is serial
		req.typz = resp.typz
		req.Data = resp.Data
		req.complete <- struct{}{}
	}
}

func (c *Client) write() {
	for msg := range c.send {
		if err := c.wire.Write(msg); err != nil {
			c.responses <- &Message{
				transportErr: err,
			}
		}
	}
}

func (c *Client) read() {
	for {
		msg, err := c.wire.Read()
		if err != nil {
			blog.Errorf("Error reading from wire: %v", err)
			c.responses <- &Message{
				transportErr: err,
			}
			break
		}
		c.responses <- msg
	}
}
