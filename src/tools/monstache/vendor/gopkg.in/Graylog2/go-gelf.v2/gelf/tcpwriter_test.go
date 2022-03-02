package gelf

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewTCPWriter(t *testing.T) {
	w, err := NewTCPWriter("")
	if err == nil && w != nil {
		t.Error("New didn't fail")
		return
	}
}

func TestNewTCPWriterConfig(t *testing.T) {
	r, _, _, err := newTCPReader("127.0.0.1:0")
	if err != nil {
		t.Error("Could not open TCPReader")
		return
	}
	w, err := NewTCPWriter(r.addr())
	if err != nil {
		t.Errorf("NewTCPWriter: %s", err)
		return
	}

	if w.MaxReconnect != 3 {
		t.Errorf("Default MaxReconnect: expected %d, got %d", 3, w.MaxReconnect)
		return
	}
	w.MaxReconnect = 5
	if w.MaxReconnect != 5 {
		t.Errorf("Custom MaxReconnect: expected %d, got %d", 5, w.MaxReconnect)
		return
	}

	if w.ReconnectDelay != 1 {
		t.Errorf("Default ReconnectDelay: expected %d, got %d", 1, w.ReconnectDelay)
		return
	}
	w.ReconnectDelay = 5
	if w.ReconnectDelay != 5 {
		t.Errorf("Custom ReconnectDelay: expected %d, got %d", 5, w.ReconnectDelay)
		return
	}
}

func assertMessages(msg *Message, short, full string, t *testing.T) {
	if msg.Short != short {
		t.Errorf("msg.Short: expected %s, got %s", short, msg.Short)
		return
	}

	if msg.Full != full {
		t.Errorf("msg.Full: expected %s, got %s", full, msg.Full)
	}

}

func TestWriteSmallMultiLineTCP(t *testing.T) {
	msgData := "awesomesauce\nbananas"

	msg, err := sendAndRecvTCP(msgData)
	if err != nil {
		t.Errorf("sendAndRecvTCP: %s", err)
		return
	}

	assertMessages(msg, "awesomesauce", msgData, t)
}

func TestWriteSmallOneLineTCP(t *testing.T) {
	msgData := "some awesome thing\n"
	msgDataTrunc := msgData[:len(msgData)-1]

	msg, err := sendAndRecvTCP(msgData)
	if err != nil {
		t.Errorf("sendAndRecvTCP: %s", err)
		return
	}

	assertMessages(msg, msgDataTrunc, "", t)

	fileExpected := "/go-gelf/gelf/tcpwriter_test.go"
	if !strings.HasSuffix(msg.Extra["_file"].(string), fileExpected) {
		t.Errorf("msg.File: expected %s, got %s", fileExpected,
			msg.Extra["_file"].(string))
		return
	}

	if len(msg.Extra) != 2 {
		t.Errorf("extra fields in %v (expect only file and line)", msg.Extra)
		return
	}
}

func TestWriteBigMessageTCP(t *testing.T) {
	randData := make([]byte, 4096)
	if _, err := rand.Read(randData); err != nil {
		t.Errorf("cannot get random data: %s", err)
		return
	}
	msgData := "awesomesauce\n" + base64.StdEncoding.EncodeToString(randData)

	msg, err := sendAndRecvTCP(msgData)
	if err != nil {
		t.Errorf("sendAndRecv: %s", err)
		return
	}

	assertMessages(msg, "awesomesauce", msgData, t)
}

func TestWriteMultiPacketMessageTCP(t *testing.T) {
	randData := make([]byte, 150000)
	if _, err := rand.Read(randData); err != nil {
		t.Errorf("cannot get random data: %s", err)
		return
	}
	msgData := "awesomesauce\n" + base64.StdEncoding.EncodeToString(randData)

	msg, err := sendAndRecvTCP(msgData)
	if err != nil {
		t.Errorf("sendAndRecv: %s", err)
		return
	}

	assertMessages(msg, "awesomesauce", msgData, t)
}

func TestExtraDataTCP(t *testing.T) {

	// time.Now().Unix() seems fine, UnixNano() won't roundtrip
	// through string -> float64 -> int64
	extra := map[string]interface{}{
		"_a":    10 * time.Now().Unix(),
		"C":     9,
		"_file": "writer_test.go",
		"_line": 186,
	}

	short := "quick"
	full := short + "\nwith more detail"
	m := Message{
		Version:  "1.0",
		Host:     "fake-host",
		Short:    string(short),
		Full:     string(full),
		TimeUnix: float64(time.Now().UnixNano()) / float64(time.Second),
		Level:    6, // info
		Facility: "writer_test",
		Extra:    extra,
		RawExtra: []byte(`{"woo": "hoo"}`),
	}

	msg, err := sendAndRecvMsgTCP(&m)
	if err != nil {
		t.Errorf("sendAndRecvMsgTCP: %s", err)
		return
	}

	assertMessages(msg, short, full, t)

	if len(msg.Extra) != 3 {
		t.Errorf("extra extra fields in %v", msg.Extra)
		return
	}

	if int64(msg.Extra["_a"].(float64)) != extra["_a"].(int64) {
		t.Errorf("_a didn't roundtrip (%v != %v)", int64(msg.Extra["_a"].(float64)), extra["_a"].(int64))
		return
	}

	if string(msg.Extra["_file"].(string)) != extra["_file"] {
		t.Errorf("_file didn't roundtrip (%v != %v)", msg.Extra["_file"].(string), extra["_file"].(string))
		return
	}

	if int(msg.Extra["_line"].(float64)) != extra["_line"].(int) {
		t.Errorf("_line didn't roundtrip (%v != %v)", int(msg.Extra["_line"].(float64)), extra["_line"].(int))
		return
	}
}

func TestWrite2MessagesWithConnectionDropTCP(t *testing.T) {
	msgData1 := "First message\nThis happens before the connection drops"
	msgData2 := "Second message\nThis happens after the connection drops"

	msg1, msg2, err := sendAndRecv2MessagesWithDropTCP(msgData1, msgData2)
	if err != nil {
		t.Errorf("sendAndRecv2MessagesWithDropTCP: %s", err)
		return
	}

	assertMessages(msg1, "First message", msgData1, t)
	assertMessages(msg2, "Second message", msgData2, t)
}

func TestWrite2MessagesWithServerDropTCP(t *testing.T) {
	msgData1 := "First message\nThis happens before the server drops"
	msgData2 := "Second message\nThis happens after the server drops"

	msg1, err := sendAndRecv2MessagesWithServerDropTCP(msgData1, msgData2)
	if err != nil {
		t.Errorf("sendAndRecv2MessagesWithDropTCP: %s", err)
		return
	}

	assertMessages(msg1, "First message", msgData1, t)
}

func setupConnections() (*TCPReader, chan string, chan string, *TCPWriter, error) {
	r, closeSignal, doneSignal, err := newTCPReader("127.0.0.1:0")

	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("newTCPReader: %s", err)
	}

	w, err := NewTCPWriter(r.addr())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("NewTCPWriter: %s", err)
	}

	return r, closeSignal, doneSignal, w, nil
}

func sendAndRecvTCP(msgData string) (*Message, error) {
	r, closeSignal, doneSignal, w, err := setupConnections()
	if err != nil {
		return nil, err
	}

	if _, err = w.Write([]byte(msgData)); err != nil {
		return nil, fmt.Errorf("w.Write: %s", err)
	}

	closeSignal <- "stop"
	done := <-doneSignal
	if done != "done" {
		return nil, errors.New("Wrong signal received")
	}

	message, err := r.readMessage()
	if err != nil {
		return nil, fmt.Errorf("r.readMessage: %s", err)
	}

	return message, nil
}

func sendAndRecvMsgTCP(msg *Message) (*Message, error) {
	r, closeSignal, doneSignal, w, err := setupConnections()
	if err != nil {
		return nil, err
	}

	if err = w.WriteMessage(msg); err != nil {
		return nil, fmt.Errorf("w.Write: %s", err)
	}

	closeSignal <- "stop"
	done := <-doneSignal
	if done != "done" {
		return nil, errors.New("Wrong signal received")
	}

	w.Close()
	message, err := r.readMessage()
	if err != nil {
		return nil, fmt.Errorf("r.readMessage: %s", err)
	}

	return message, nil
}

func sendAndRecv2MessagesWithDropTCP(msgData1 string, msgData2 string) (*Message, *Message, error) {
	r, closeSignal, doneSignal, w, err := setupConnections()
	if err != nil {
		return nil, nil, err
	}

	if _, err = w.Write([]byte(msgData1)); err != nil {
		return nil, nil, fmt.Errorf("w.Write: %s", err)
	}

	time.Sleep(200 * time.Millisecond)

	closeSignal <- "drop"
	done := <-doneSignal
	if done != "done" {
		return nil, nil, fmt.Errorf("Wrong signal received: %s", done)
	}

	message1, err := r.readMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("readmessage: %s", err)
	}

	// Need to write twice to force the detection of the dropped connection
	if _, err = w.Write([]byte(msgData2)); err != nil {
		return nil, nil, fmt.Errorf("write 1 w.Write: %s", err)
	}
	time.Sleep(200 * time.Millisecond)
	if _, err = w.Write([]byte(msgData2)); err != nil {
		return nil, nil, fmt.Errorf("write 2 w.Write: %s", err)
	}
	time.Sleep(200 * time.Millisecond)

	closeSignal <- "stop"
	done = <-doneSignal
	if done != "done" {
		return nil, nil, errors.New("Wrong signal received")
	}

	message2, err := r.readMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("readmessage: %s", err)
	}

	w.Close()
	return message1, message2, nil
}

func sendAndRecv2MessagesWithServerDropTCP(msgData1 string, msgData2 string) (*Message, error) {
	r, closeSignal, doneSignal, w, err := setupConnections()
	if err != nil {
		return nil, err
	}

	if _, err = w.Write([]byte(msgData1)); err != nil {
		return nil, fmt.Errorf("w.Write: %s", err)
	}

	closeSignal <- "stop"
	done := <-doneSignal
	if done != "done" {
		return nil, fmt.Errorf("Wrong signal received: %s", done)
	}

	message1, err := r.readMessage()
	if err != nil {
		return nil, fmt.Errorf("readmessage: %s", err)
	}

	// Need to write twice to force the detection of the dropped connection
	// The first write will not cause an error, but the subsequent ones will
	for {
		_, err = w.Write([]byte(msgData2))
		if err != nil {
			break
		}
	}

	w.Close()
	return message1, nil
}
