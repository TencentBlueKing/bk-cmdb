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
	"encoding/json"
	"errors"
	"strings"
)

// Errors define
var (
	//ErrRWTimeout r/w operation timeout
	ErrRWTimeout         = errors.New("r/w timeout")
	ErrPingTimeout       = errors.New("Ping timeout")
	ErrCommandOverLength = errors.New("Command overlength")
	ErrCommandNotFount   = errors.New("Command not found")
)

type HandlerFunc func(*Message) (interface{}, error)
type HandlerStreamFunc func(input <-chan *Message, output <-chan *Message) error

const (
	TypeRequest = iota
	TypeResponse
	TypeError
	TypeClose
	TypePing

	readBufferSize  = 8096
	writeBufferSize = 8096
)

type Codec uint32

const (
	CodexJSON Codec = iota
	CodexGob
)

type Decoder interface {
	Decode(interface{}) error
}

var (
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

const (
	MagicVersion = uint16(0x1b01) // cmdb01
)

type command [40]byte

func NewCommand(cmd string) (command, error) {
	ncmd := command{}
	cmdlength := len(cmd)
	if len(cmd) > commanLimit {
		return ncmd, ErrCommandOverLength
	}
	copy(ncmd[:], []byte(cmd)[:cmdlength])
	return ncmd, nil
}
func (c command) String() string {
	return strings.TrimSpace(string(c[:]))
}

type Message struct {
	complete     chan struct{}
	transportErr error

	magicVersion uint16
	seq          uint32
	typz         uint32
	cmd          command // maybe should use uint32

	Codec Codec
	Size  uint32
	Data  []byte
}

func (msg *Message) Decode(value interface{}) error {
	if msg.Codec == CodexJSON {
		return json.Unmarshal(msg.Data, value)
	}
	return ErrUnsupportedCodec
}
func (msg *Message) Encode(value interface{}) error {
	var err error
	if msg.Codec == CodexJSON {
		msg.Data, err = json.Marshal(value)
		return err
	}

	return ErrUnsupportedCodec
}
