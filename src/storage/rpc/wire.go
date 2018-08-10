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
	"encoding/binary"
	"fmt"
	"io"
)

type Wire interface {
	Write(*Message) error
	Read() (*Message, error)
	Close() error
}

type BinaryWire struct {
	conn   io.ReadWriteCloser
	writer *bufio.Writer
	reader io.Reader
}

func NewBinaryWire(rwc io.ReadWriteCloser) *BinaryWire {
	return &BinaryWire{
		conn:   rwc,
		writer: bufio.NewWriterSize(rwc, writeBufferSize),
		reader: bufio.NewReaderSize(rwc, readBufferSize),
	}
}

func (w *BinaryWire) Write(msg *Message) error {
	// LittleEndian: x86 cpu 为小端字节序
	// 如 0x01234567，地址范围为0x100~0x103字节,小端字节序则存储为: 0x100: 67, 0x101: 45,..
	if err := binary.Write(w.writer, binary.LittleEndian, msg.magicVersion); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, msg.seq); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, msg.typz); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, msg.cmd); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, msg.Codec); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, msg.Size); err != nil {
		return err
	}
	if err := binary.Write(w.writer, binary.LittleEndian, uint32(len(msg.Data))); err != nil {
		return err
	}
	if len(msg.Data) > 0 {
		if _, err := w.writer.Write(msg.Data); err != nil {
			return err
		}
	}
	return w.writer.Flush()
}

func (w *BinaryWire) Read() (*Message, error) {
	var (
		msg    Message
		length uint32
	)

	if err := binary.Read(w.reader, binary.LittleEndian, &msg.magicVersion); err != nil {
		return nil, err
	}

	if msg.magicVersion != MagicVersion {
		return nil, fmt.Errorf("Wrong API version received: 0x%x", &msg.magicVersion)
	}

	if err := binary.Read(w.reader, binary.LittleEndian, &msg.seq); err != nil {
		return nil, err
	}
	if err := binary.Read(w.reader, binary.LittleEndian, &msg.typz); err != nil {
		return nil, err
	}
	if err := binary.Read(w.reader, binary.LittleEndian, &msg.cmd); err != nil {
		return nil, err
	}
	if err := binary.Read(w.reader, binary.LittleEndian, &msg.Codec); err != nil {
		return nil, err
	}
	if err := binary.Read(w.reader, binary.LittleEndian, &msg.Size); err != nil {
		return nil, err
	}
	if err := binary.Read(w.reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	if length > 0 {
		msg.Data = make([]byte, length)
		if _, err := io.ReadFull(w.reader, msg.Data); err != nil {
			return nil, err
		}
	}
	return &msg, nil
}

func (w *BinaryWire) Close() error {
	return w.conn.Close()
}
