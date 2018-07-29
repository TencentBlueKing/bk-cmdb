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
	"context"
	"fmt"
	"net"
	"net/rpc"
)

// Server rpc server
type Server interface {
	Run(ctx context.Context) error
}

// NewServer create a new rpc server instance
func NewServer(cfg Config, actions []interface{}) Server {
	return &server{
		actions: actions,
		conf:    cfg,
	}
}

type server struct {
	actions []interface{}
	conf    Config
}

func (s *server) Run(ctx context.Context) error {

	for _, action := range s.actions {
		if err := rpc.Register(action); nil != err {
			return err
		}
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.conf.IPAddr, s.conf.Port))
	if nil != err {
		return err
	}

	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := tcpListen.Accept()
			if nil != err {
				fmt.Println("failed to accept, error info is ", err.Error())
				continue
			}
			go rpc.ServeConn(conn)
		}
	}

}
