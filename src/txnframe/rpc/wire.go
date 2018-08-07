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
