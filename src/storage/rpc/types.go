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
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"sync"

	"configcenter/src/common/util"
)

// Errors define
var (
	//ErrRWTimeout r/w operation timeout
	ErrRWTimeout       = errors.New("r/w timeout")
	ErrPingTimeout     = errors.New("Ping timeout")
	ErrCommandNotFount = errors.New("Command not found")
	ErrStreamStoped    = errors.New("Stream stoped")
)

// HandlerFunc define a HandlerFunc
type HandlerFunc func(Request) (interface{}, error)

// HandlerStreamFunc define a HandlerStreamFunc
type HandlerStreamFunc func(Request, ServerStream) error

type streamstore struct {
	sync.RWMutex
	stream map[uint32]*StreamMessage
}

func newStreamStore() *streamstore {
	return &streamstore{
		stream: map[uint32]*StreamMessage{},
	}
}

func (s *streamstore) store(seq uint32, stream *StreamMessage) {
	s.Lock()
	s.stream[seq] = stream
	s.Unlock()
}

func (s *streamstore) get(seq uint32) (*StreamMessage, bool) {
	stream, ok := s.stream[seq]
	return stream, ok
}
func (s *streamstore) remove(seq uint32) {
	s.Lock()
	delete(s.stream, seq)
	s.Unlock()
}

// StreamMessage define
type StreamMessage struct {
	root   *Message
	input  chan *Message
	output chan *Message
	done   *util.AtomicBool
	err    error
}

// ServerStream define interface
type ServerStream interface {
	Recv(result interface{}) error
	Send(data interface{}) error
}

// NewStreamMessage returns a new StreamMessage
func NewStreamMessage(root *Message) *StreamMessage {
	return &StreamMessage{
		root:   root,
		input:  make(chan *Message, 10),
		output: make(chan *Message, 10),
		done:   util.NewBool(false),
	}
}

// Recv receive message
func (m StreamMessage) Recv(result interface{}) error {
	if m.err != nil {
		return m.err
	}
	msg := <-m.input
	if msg.typz == TypeStreamClose {
		m.err = ErrStreamStoped
		if len(msg.Data) > 0 {
			return errors.New(string(msg.Data))
		}
		return m.err
	}
	return msg.Decode(result)
}

// Send send message
func (m StreamMessage) Send(data interface{}) error {
	if m.err != nil {
		return m.err
	}
	msg := m.root.copy()
	msg.typz = TypeStream
	if err := msg.Encode(data); err != nil {
		return err
	}
	m.output <- msg
	return nil
}

// Close should only call by client
func (m StreamMessage) Close() error {
	msg := m.root.copy()
	msg.typz = TypeStreamClose
	m.output <- msg
	return nil
}

// MessageType define
type MessageType uint32

// MessageType enumeration
const (
	TypeRequest MessageType = iota
	TypeResponse
	TypeStream
	TypeError
	TypeClose
	TypePing
	TypeStreamClose
)

const (
	readBufferSize  = 8096
	writeBufferSize = 8096
)

// Codec define a codec
type Codec interface {
	Decode(data []byte, v interface{}) error
	Encode(v interface{}) ([]byte, error)
}

type jsonCodec struct{}

// JSONCodec implements Codec interface
var JSONCodec Codec = new(jsonCodec)

func (jsonCodec) Decode(data []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}
func (jsonCodec) Encode(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

const (
	// MagicVersion is the cc rpc protocol version
	MagicVersion = uint16(0x1b01) // cmdb01
)

// Request define a request interface
type Request interface {
	Decode(value interface{}) error
}

// Message define a rpc message
type Message struct {
	complete     chan struct{}
	transportErr error
	codec        Codec

	magicVersion uint16
	seq          uint32
	typz         MessageType
	cmd          string // maybe should use uint32
	Data         []byte
}

func (msg Message) copy() *Message {
	return &Message{
		magicVersion: msg.magicVersion,
		seq:          msg.seq,
		typz:         msg.typz,
		cmd:          msg.cmd,
		codec:        msg.codec,
	}
}

// Decode decode the message data
func (msg *Message) Decode(value interface{}) error {
	if decoder, ok := value.(encoding.BinaryUnmarshaler); ok {
		return decoder.UnmarshalBinary(msg.Data)
	}
	return msg.codec.Decode(msg.Data, value)
}

// Encode encode the value to message data
func (msg *Message) Encode(value interface{}) error {
	if value == nil {
		msg.Data = msg.Data[:0]
		return nil
	}
	var err error
	if encoder, ok := value.(encoding.BinaryMarshaler); ok {
		msg.Data, err = encoder.MarshalBinary()
	} else {
		msg.Data, err = msg.codec.Encode(value)
	}
	return err
}

type ClientConfig struct {
	Address string
}
