package httpserver

import (
	"github.com/emicklei/go-restful"
)

type Server interface {
	ListenAndServe() error
	RegisterActions(as ...Action)
}

type Action struct {
	Method  string
	Path    string
	Handler restful.RouteFunction // actually this should be http.HandlerFunc
}
