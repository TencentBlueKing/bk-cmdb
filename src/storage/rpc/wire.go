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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Wire define a wire
type Wire interface {
	Write(*Message) error
	Read(*Message,) ( error)
	Close() error
}

// BinaryWire implements Wire interface
type BinaryWire struct {
	conn   io.ReadWriteCloser
	writer flushWriter
	reader io.Reader
}

// NewBinaryWire returns a new BinaryWire
func NewBinaryWire(rwc io.ReadWriteCloser, compress string) (*BinaryWire, error) {
	compressor, err := newCompressor(rwc,rwc, compress)
	if err != nil {
		return nil, err
	}
	return &BinaryWire{
		conn:   rwc,
		writer: compressor,
		reader: compressor,
	}, nil
}

func (w *BinaryWire) Write(msg *Message) error {
	if msg == nil {
		return errors.New("wire could not write empty message")
	}
	var err error
	// LittleEndian: x86 cpu 为小端字节序
	// 如 0x01234567，地址范围为0x100~0x103字节,小端字节序则存储为: 0x100: 67, 0x101: 45,..
	if err = binary.Write(w.writer, binary.LittleEndian, msg.magicVersion); err != nil {
		return err
	}
	if err = binary.Write(w.writer, binary.LittleEndian, msg.seq); err != nil {
		return err
	}
	if err = binary.Write(w.writer, binary.LittleEndian, msg.typz); err != nil {
		return err
	}
	if err = writeString(w.writer, msg.cmd); err != nil {
		return err
	}
	if err = writeBytes(w.writer, msg.Data); err != nil {
		return err
	}
	return w.writer.Flush()
}

func (w *BinaryWire) Read(msg *Message) (error) {
	var (
		err error
	)

	if err = binary.Read(w.reader, binary.LittleEndian, &msg.magicVersion); err != nil {
		return err
	}

	if msg.magicVersion != MagicVersion {
		return fmt.Errorf("Wrong API version received: 0x%x", &msg.magicVersion)
	}

	if err = binary.Read(w.reader, binary.LittleEndian, &msg.seq); err != nil {
		return err
	}
	if err = binary.Read(w.reader, binary.LittleEndian, &msg.typz); err != nil {
		return err
	}
	if msg.cmd, err = readString(w.reader); err != nil {
		return err
	}
	if msg.Data, err = readBytes(w.reader); err != nil {
		return err
	}
	return  nil
}

// Close close the wire
func (w *BinaryWire) Close() error {
	return w.conn.Close()
}

func readString(reader io.Reader) (string, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	if length <= 0 {
		return "", nil
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return "", err
	}
	return string(data), nil
}

func writeString(writer io.Writer, s string) error {
	if err := binary.Write(writer, binary.LittleEndian, uint32(len(s))); err != nil {
		return err
	}
	if _, err := io.WriteString(writer, s); err != nil {
		return err
	}
	return nil
}

func readBytes(reader io.Reader) (data []byte, err error) {
	var length uint32
	if err = binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	if length > 0 {
		data = make([]byte, length)
		if _, err = io.ReadFull(reader, data); err != nil {
			return nil, err
		}
	}
	return
}

func writeBytes(writer io.Writer, data []byte) (err error) {
	if err = binary.Write(writer, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err = writer.Write(data); err != nil {
			return err
		}
	}
	return
}
