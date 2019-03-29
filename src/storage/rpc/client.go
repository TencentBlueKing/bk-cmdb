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
	"sync"
	"sync/atomic"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

var (
	opRetries     = 0
	opReadTimeout = 15 * time.Second // client read
	opPingTimeout = 20 * time.Second
)

//Client replica client
type Client struct {
	send         chan *Message
	seq          uint32
	messageMutex sync.RWMutex
	messages     map[uint32]*Message
	stream       *streamstore
	wire         Wire
	peerAddr     string
	err          error
	codec        Codec

	request  Message
	response Message
	done     *util.AtomicBool
	wg       sync.WaitGroup
}

//NewClient replica client
func NewClient(conn net.Conn, compress string) (*Client, error) {
	wire, err := NewBinaryWire(conn, compress)
	if err != nil {
		return nil, err
	}
	c := &Client{
		wire:     wire,
		peerAddr: conn.RemoteAddr().String(),
		done:     util.NewBool(false),
		send:     make(chan *Message, 1024),
		messages: map[uint32]*Message{},
		codec:    JSONCodec,
		stream:   newStreamStore(),
	}
	blog.V(3).Infof("connected to %s", c.TargetID())
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
	c.done.Set()
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
		for streammsg := range sm.output {
			c.send <- streammsg
			if msg.typz == TypeStreamClose {
				break
			}
			if sm.done.IsSet() || c.done.IsSet() {
				break
			}
		}
		sm.err = ErrStreamStoped
		close(sm.input)
		close(sm.output)
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
		msg := &Message{
			magicVersion: MagicVersion,
			codec:        c.codec,
			seq:          c.nextSeq(),
			complete:     make(chan struct{}),
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

		c.handleRequest(msg)

		select {
		case <-msg.complete:
			if msg.typz == TypeError {
				return nil, errors.New(string(msg.Data))
			}
			return msg, nil
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

func (c *Client) nextSeq() uint32 {
	return atomic.AddUint32(&c.seq, 1)
}

func (c *Client) replyError(req *Message) {
	c.messageMutex.Lock()
	delete(c.messages, req.seq)
	c.messageMutex.Unlock()

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

	c.messageMutex.Lock()
	c.messages[req.seq] = req
	c.messageMutex.Unlock()

	c.send <- req
	blog.V(5).Infof("[rpc client]sent message data: %s", req.Data)
}

func (c *Client) handleResponse(resp *Message) {
	blog.V(5).Infof("[rpc client]receive message data: %s", resp.Data)
	resp.codec = c.codec
	if resp.transportErr != nil {
		c.err = resp.transportErr
		// Terminate all in flight
		c.messageMutex.RLock()
		for _, msg := range c.messages {
			c.replyError(msg)
		}
		c.messageMutex.RUnlock()
		// TODO reconnect when transportErr?
		return
	}
	if resp.typz == TypeStream || resp.typz == TypeStreamClose {
		c.stream.RLock()
		stream, ok := c.stream.get(resp.seq)
		c.stream.RUnlock()
		if ok {
			stream.input <- resp
		} else {
			blog.Warnf("[rpc client] stream not found, resp is %s", resp.Data)
		}
		return
	}
	c.messageMutex.Lock()
	if req, ok := c.messages[resp.seq]; ok {
		if c.err != nil {
			c.messageMutex.Unlock()
			c.replyError(req)
			return
		}

		delete(c.messages, resp.seq)
		c.messageMutex.Unlock()

		req.typz = resp.typz
		req.Data = resp.Data
		close(req.complete)
	} else {
		c.messageMutex.Unlock()
	}
}

func (c *Client) write() {
	for msg := range c.send {
		if err := c.wire.Write(msg); err != nil {
			c.handleResponse(msg)
		}
	}
}

func (c *Client) read() {
	for {
		err := c.wire.Read(&c.response)
		if err != nil {
			blog.Errorf("Error reading from wire: %v", err)
			c.handleResponse(&c.response)
			break
		}
		c.handleResponse(&c.response)
	}
}
