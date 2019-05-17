package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadArray(t *testing.T) {
	type cas struct {
		payload string
		err     error
		res     []string
	}
	for i, c := range []cas{
		{
			payload: "*1\r\n$4\r\nPING\r\n",
			res:     []string{"PING"},
		},
		{
			payload: "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n",
			res:     []string{"LLEN", "mylist"},
		},
		{
			payload: "*2\r\n$4\r\nLLEN\r\n$6\r\nmyl",
			err:     io.EOF,
		},
		{
			payload: "PING",
			err:     io.EOF,
		},
		{
			payload: "*0\r\n",
		},
		{
			payload: "*-1\r\n", // not sure this is legal in a request
		},
	} {
		res, err := readArray(bufio.NewReader(bytes.NewBufferString(c.payload)))
		if have, want := err, c.err; have != want {
			t.Errorf("err %d: have %v, want %v", i, have, want)
			continue
		}
		if have, want := res, c.res; !reflect.DeepEqual(have, want) {
			t.Errorf("case %d: have %v, want %v", i, have, want)
		}
	}
}

func TestReadString(t *testing.T) {
	type cas struct {
		payload string
		err     error
		res     string
	}
	bigPayload := strings.Repeat("X", 1<<24)
	for i, c := range []cas{
		{
			payload: "+hello world\r\n",
			res:     "hello world",
		},
		{
			payload: "-some error\r\n",
			res:     "some error",
		},
		{
			payload: ":42\r\n",
			res:     "42",
		},
		{
			payload: ":\r\n",
			res:     "",
		},
		{
			payload: "$4\r\nabcd\r\n",
			res:     "abcd",
		},
		{
			payload: fmt.Sprintf("$%d\r\n%s\r\n", len(bigPayload), bigPayload),
			res:     bigPayload,
		},

		{
			payload: "",
			err:     io.EOF,
		},
		{
			payload: ":42",
			err:     io.EOF,
		},
		{
			payload: "XXX",
			err:     io.EOF,
		},
		{
			payload: "XXXXXX",
			err:     io.EOF,
		},
		{
			payload: "\r\n",
			err:     ErrProtocol,
		},
		{
			payload: "XXXX\r\n",
			err:     ErrProtocol,
		},
	} {
		res, err := readString(bufio.NewReader(bytes.NewBufferString(c.payload)))
		if have, want := err, c.err; have != want {
			t.Errorf("err %d: have %v, want %v", i, have, want)
			continue
		}
		if have, want := res, c.res; !reflect.DeepEqual(have, want) {
			t.Errorf("case %d: have %#v, want %#v", i, have, want)
		}
	}
}
