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
}

var _ Server = &HttpServer{}

func NewServer(opt *option.Options) (Server, error) {
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
