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

package rpc_test

import (
	"configcenter/src/common/util"
	"configcenter/src/storage/rpc"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStream(t *testing.T) {
	fmt.Println("test started")
	type Param struct {
		Args string
	}

	type Req struct {
		Args string
	}

	type Resp struct {
		Args string
	}

	var streamfunc = func(param *rpc.Message, stream rpc.ServerStream) error {
		fmt.Fprintln(os.Stdout, "server started")
		var p = Param{}
		err := param.Decode(&p)
		require.NoError(t, err)
		require.Equal(t, "param", p.Args)
		var req = Req{}
		err = stream.Recv(&req)
		require.NoError(t, err)
		require.Equal(t, "req", req.Args)

		err = stream.Send(&Resp{Args: "resp"})
		require.NoError(t, err)

		fmt.Fprintln(os.Stdout, "server stoped")
		return nil
	}

	rpcsrv := rpc.NewServer()
	rpcsrv.HandleStream("streamtest", streamfunc)
	mux := http.NewServeMux()
	mux.Handle("/rpc", rpcsrv)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	address, err := util.GetDailAddress(ts.URL)
	require.NoError(t, err)
	cli, err := rpc.DialHTTPPath("tcp", address, "/rpc")
	require.NoError(t, err)

	stream, err := cli.CallStream("streamtest", &Param{Args: "param"})
	require.NoError(t, err)

	err = stream.Send(&Req{Args: "req"})
	require.NoError(t, err)

	resp := Resp{}
	err = stream.Recv(&resp)
	require.NoError(t, err)
	require.Equal(t, "resp", resp.Args)

	err = stream.Recv(&resp)
	require.EqualError(t, err, rpc.ErrStreamStoped.Error())

	err = stream.Close()
	require.NoError(t, err)

}
