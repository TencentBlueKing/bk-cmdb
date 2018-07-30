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
	"configcenter/src/common/blog"
	"configcenter/src/common/ssl"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

// HttpServer is data struct of http server
type HttpServer struct {
	addr         string
	port         uint
	sock         string
	isSSL        bool
	caFile       string
	certFile     string
	keyFile      string
	certPasswd   string
	webContainer *restful.Container
}

func NewHttpServer(port uint, addr, sock string) *HttpServer {

	wsContainer := restful.NewContainer()

	// AddUserConfig container filter to enable CORS
	//	cors := restful.CrossOriginResourceSharing{
	//		AllowedHeaders: []string{"Content-Type", "Accept"},
	//		AllowedDomains: []string{},
	//		CookiesAllowed: true,
	//		Container:      wsContainer}
	//	wsContainer.Filter(cors.Filter)
	//	wsContainer.Filter(wsContainer.OPTIONSFilter)
	return &HttpServer{
		addr:         addr,
		port:         port,
		sock:         sock,
		webContainer: wsContainer,
		isSSL:        false,
	}
}

func (s *HttpServer) SetSsl(cafile, certfile, keyfile, certPasswd string) {
	s.caFile = cafile
	s.certFile = certfile
	s.keyFile = keyfile
	s.certPasswd = certPasswd
	s.isSSL = true
}

func (s *HttpServer) RegisterWebServer(rootPath string, filter restful.FilterFunction, actions []*Action) error {
	//new a web service
	ws := s.NewWebService(rootPath, filter)

	//register action
	s.RegisterActions(ws, actions)

	return nil
}

func (s *HttpServer) NewWebService(rootPath string, filter restful.FilterFunction) *restful.WebService {
	ws := new(restful.WebService)
	if "" != rootPath {
		ws.Path(rootPath)
	}

	ws.Produces(restful.MIME_JSON)

	if nil != filter {
		ws.Filter(filter)
	}

	s.webContainer.Add(ws)

	return ws
}
func (s *HttpServer) GetWebContainer() *restful.Container {
	return s.webContainer
}
func (s *HttpServer) RegisterActions(ws *restful.WebService, actions []*Action) {
	blog.Debug("RegisterActions")
	for _, action := range actions {
		switch action.Verb {
		case "POST":
			route := ws.POST(action.Path).To(action.Handler)
			s.registerActionsFilter(route, action.FilterHandler)
			ws.Route(route)
			blog.Debug("register post api, url(%s)", action.Path)
		case "GET":
			route := ws.GET(action.Path).To(action.Handler)
			s.registerActionsFilter(route, action.FilterHandler)
			ws.Route(route)
			blog.Debug("register get api, url(%s)", action.Path)
		case "PUT":
			route := ws.PUT(action.Path).To(action.Handler)
			s.registerActionsFilter(route, action.FilterHandler)
			ws.Route(route)
			blog.Debug("register put api, url(%s)", action.Path)
		case "DELETE":
			route := ws.DELETE(action.Path).To(action.Handler)
			s.registerActionsFilter(route, action.FilterHandler)
			ws.Route(route)
			blog.Debug("register delete api, url(%s)", action.Path)
		default:
			blog.Error("unrecognized action verb: %s", action.Verb)
		}
	}
}

func (s *HttpServer) registerActionsFilter(r *restful.RouteBuilder, filters []restful.FilterFunction) {
	for _, f := range filters {
		r.Filter(f)
	}
}

func (s *HttpServer) ListenAndServe() error {

	var chError = make(chan error)
	//list and serve by addrport
	go func() {
		addrport := net.JoinHostPort(s.addr, strconv.FormatUint(uint64(s.port), 10))
		httpserver := &http.Server{Addr: addrport, Handler: s.webContainer}
		if s.isSSL {
			tlsConf, err := ssl.ServerTslConf(s.caFile, s.certFile, s.keyFile, s.certPasswd)
			if err != nil {
				blog.Error("fail to load certfile, err:%s", err.Error())
				chError <- fmt.Errorf("fail to load certfile")
				return
			}
			httpserver.TLSConfig = tlsConf
			blog.Info("Start https service on(%s)", addrport)
			chError <- httpserver.ListenAndServeTLS("", "")
		} else {
			blog.Info("Start http service on(%s)", addrport)
			chError <- httpserver.ListenAndServe()
		}
	}()

	return <-chError
}
