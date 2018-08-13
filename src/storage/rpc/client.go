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
	"configcenter/src/common/blog"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
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
	wire      Wire
	peerAddr  string
	err       error
}

//NewClient replica client
func NewClient(conn net.Conn) *Client {
	c := &Client{
		wire:      NewBinaryWire(conn),
		peerAddr:  conn.RemoteAddr().String(),
		end:       make(chan struct{}, 1024),
		requests:  make(chan *Message, 1024),
		send:      make(chan *Message, 1024),
		responses: make(chan *Message, 1024),
		messages:  map[uint32]*Message{},
	}
	blog.V(3).Infof("connected to %s", c.TargetID())
	go c.loop()
	go c.write()
	go c.read()
	return c
}

// DialHTTPPath connects to an HTTP RPC server
// at the specified network address and path.
func DialHTTPPath(network, address, path string) (*Client, error) {
	var err error
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.Status == connected {
		return NewClient(conn), nil
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

//TargetID operation target ID
func (c *Client) TargetID() string {
	return c.peerAddr
}

//Call replica client
func (c *Client) Call(cmd string, input interface{}, result interface{}) error {
	cmdlength := len(cmd)
	if len(cmd) > commanLimit {
		return ErrCommandOverLength
	}

	ncmd := command{}
	copy(ncmd[:], []byte(cmd)[:cmdlength])

	msg, err := c.operation(TypeRequest, command(ncmd), input)
	if err != nil {
		return err
	}
	return msg.Decode(result)
}

//Ping replica client
func (c *Client) Ping() error {
	_, err := c.operation(TypePing, command{}, nil)
	return err
}

func (c *Client) operation(op uint32, cmd command, data interface{}) (Decoder, error) {
	retry := 0
	for {
		msg := Message{
			complete: make(chan struct{}, 1),
			typz:     op,
			Codec:    CodexJSON,
			cmd:      cmd,
			Data:     nil,
		}

		if op == TypeRequest {
			msg.Encode(data)
			msg.Size = uint32(len(msg.Data))
		}

		timeout := func(op uint32) <-chan time.Time {
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

	req.magicVersion = MagicVersion
	req.seq = c.nextSeq()
	c.messages[req.seq] = req
	c.send <- req
	blog.V(4).Infof("[rpc client]sent message data: %s", req.Data)
}

func (c *Client) handleResponse(resp *Message) {
	blog.V(4).Infof("[rpc client]receive message data: %s", resp.Data)
	if resp.transportErr != nil {
		c.err = resp.transportErr
		// Terminate all in flight
		for _, msg := range c.messages {
			c.replyError(msg)
		}
		// TODO reconnect when transportErr?
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
