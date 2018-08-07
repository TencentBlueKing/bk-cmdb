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
type HandlerStreamFunc func(chan<- *Message) error

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
	id           uint32 // Seq and ID can apparently be collapsed into one (ID)
	seq          uint32
	typz         uint32
	cmd          command

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
