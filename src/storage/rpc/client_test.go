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
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	gorpc "net/rpc"
	"testing"

	"github.com/stretchr/testify/require"

	"configcenter/src/common/util"
)

type Req struct {
	Name string
}

func (r *Req) MarshalBinary() (data []byte, err error) {
	buf := &bytes.Buffer{}
	writeString(buf, r.Name)
	return buf.Bytes(), nil
}
func (r *Req) UnmarshalBinary(data []byte) error {
	r.Name, _ = readString(bytes.NewBuffer(data))
	return nil
}

type Reply struct {
	OK bool
}

func (r *Reply) MarshalBinary() (data []byte, err error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, r.OK)
	return buf.Bytes(), nil
}
func (r *Reply) UnmarshalBinary(data []byte) error {
	binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &r.OK)
	return nil
}

func OK(msg Request) (interface{}, error) {
	return Reply{OK: true}, nil
}

func BenchmarkGORPC(b *testing.B) {

	rpc := gorpc.NewServer()
	rpc.Register(new(GORPC))

	mux := http.NewServeMux()
	mux.Handle("/rpc", rpc)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	address, err := util.GetDailAddress(ts.URL)
	require.NoError(b, err)
	cli, err := gorpc.DialHTTPPath("tcp", address, "/rpc")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reply := Reply{}
		err := cli.Call("GORPC.OK", &Req{Name: "ok"}, &reply)
		require.NoError(b, err)
		require.True(b, reply.OK)
	}
}
func BenchmarkGORPCParallel(b *testing.B) {

	rpc := gorpc.NewServer()
	rpc.Register(new(GORPC))

	mux := http.NewServeMux()
	mux.Handle("/rpc", rpc)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	address, err := util.GetDailAddress(ts.URL)
	require.NoError(b, err)
	cli, err := gorpc.DialHTTPPath("tcp", address, "/rpc")
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reply := Reply{}
			err := cli.Call("GORPC.OK", &Req{Name: "ok"}, &reply)
			require.NoError(b, err)
			require.True(b, reply.OK)
		}
	})
}

type GORPC struct{}

func (*GORPC) OK(req Req, reply *Reply) error {
	*reply = Reply{OK: true}
	return nil
}

func BenchmarkRPC(b *testing.B) {
	rpc := NewServer()

	mux := http.NewServeMux()
	rpc.Handle("ok", OK)
	mux.Handle("/rpc", rpc)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	address, err := util.GetDailAddress(ts.URL)
	require.NoError(b, err)
	cli, err := DialHTTPPath("tcp", address, "/rpc")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reply := Reply{}
		err := cli.Call("ok", &Req{Name: "ok"}, &reply)
		require.NoError(b, err)
		require.True(b, reply.OK)
	}
}

func BenchmarkRPCParallel(b *testing.B) {
	rpc := NewServer()

	mux := http.NewServeMux()
	rpc.Handle("ok", OK)
	mux.Handle("/rpc", rpc)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	address, err := util.GetDailAddress(ts.URL)
	require.NoError(b, err)
	cli, err := DialHTTPPath("tcp", address, "/rpc")
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reply := Reply{}
			err := cli.Call("ok", &Req{Name: "ok"}, &reply)
			require.NoError(b, err)
			require.True(b, reply.OK)
		}
	})
}
