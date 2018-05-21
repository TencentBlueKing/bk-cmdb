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

package httpserver

import (
	"configcenter/src/common/http/httpserver"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/option"
	"github.com/emicklei/go-restful"
	"net"
	"strconv"
)

type HttpServer struct {
	server *httpserver.HttpServer
	rootWS *restful.WebService

	addr string
	port uint
}

var _ Server = &HttpServer{}

func NewServer(opt *option.Options) (*HttpServer, error) {
	addr, port, err := net.SplitHostPort(opt.Addrport)
	if err != nil {
		log.Errorf("get address error %v", *opt)
		return nil, err
	}

	p, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return nil, err
	}

	server := httpserver.NewHttpServer(uint(p), addr, "")
	ws := server.NewWebService("/", nil)

	return &HttpServer{
		server: server,
		rootWS: ws,
	}, nil
}

func (s *HttpServer) GetPort() uint {
	return s.port
}
func (s *HttpServer) GetAddr() string {
	return s.addr
}

func (s *HttpServer) RegisterActions(as ...Action) {
	var httpactions []*httpserver.Action
	for _, a := range as {
		httpactions = append(httpactions, &httpserver.Action{Verb: a.Method, Path: a.Path, Handler: a.Handler})
	}
	s.server.RegisterActions(s.rootWS, httpactions)
}

func (s *HttpServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
