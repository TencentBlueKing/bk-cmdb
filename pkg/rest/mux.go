/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package rest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type ctxKey string

var (
	patternCtxKey = ctxKey("mux.pattern")
)

// Router consisting of the core routing methods using only the standard net/http handler
type Router interface {
	http.Handler

	// Use appends one or more middlewares onto the Router stack.
	// Note: This should be used only before any routes are added to the mux.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a sub-Router along a `pattern` string.
	Route(pattern string, fn func(r Router)) Router

	// Mount attaches another http.Handler along `pattern/` string
	Mount(pattern string, h http.Handler)

	// Handle and HandleFunc adds routes for `pattern{$}` that matches all HTTP methods.
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)

	// HTTP-method routing along `pattern{$}`
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)

	// NotFound defines a handler to respond whenever a route could not be found.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is not allowed.
	MethodNotAllowed(h http.HandlerFunc)
}

// router is a HTTP route multiplexer that parses a request path, base the standard net/http mux
type router struct {
	// The underlying mux to register the routes to
	mux *http.ServeMux
	// The Subrouter base path
	basePath string
	// The middleware stack
	middlewares []func(http.Handler) http.Handler
	// Custom route not found handler
	notFoundHandler http.HandlerFunc
	// Custom method not allowed handler
	methodNotAllowedHandler http.HandlerFunc
}

// NewRouter returns a new mux object that implements the Router interface.
func NewRouter() Router {
	r := &router{
		mux:         http.NewServeMux(),
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
	return r
}

// Use appends a middleware handler to the mux middleware stack.
func (r *router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// With adds inline middlewares for an endpoint handler.
func (r *router) With(middlewares ...func(http.Handler) http.Handler) Router {
	newRouter := r.clone()
	newRouter.middlewares = append(newRouter.middlewares, middlewares...)
	return newRouter
}

// Group creates a new inline-mux with a copy of middleware stack. It's useful
// for a group of handlers along the same routing path that use an additional
// set of middlewares.
func (r *router) Group(fn func(r Router)) Router {
	newRouter := r.clone()
	if fn != nil {
		fn(newRouter)
	}

	return newRouter
}

// Route creates a new mux and mounts it along the `pattern` as a subrouter.
func (r *router) Route(pattern string, fn func(r Router)) Router {
	newRouter := r.clone()
	newRouter.basePath = path.Join(r.basePath, pattern)
	if fn != nil {
		fn(newRouter)
	}

	return newRouter
}

// Mount attaches another http.Handler or mux Router as a subrouter along a routing
// path. It's very useful to split up a large API as many independent routers and
// compose them as a single service using Mount
func (r *router) Mount(pattern string, handler http.Handler) {
	if pattern == "" || pattern[0] != '/' {
		panic(fmt.Errorf("pattern must begin with /"))
	}

	basePattern := strings.TrimRight(pattern, "/")
	r.mux.Handle(basePattern+"/", r.stripMountPrefix(basePattern, r.chain(handler)))
}

// stripMountPrefix trims the mount prefix or {var} from the url if present.
func (r *router) stripMountPrefix(pattern string, handler http.Handler) http.Handler {
	if pattern == "" {
		return handler
	}

	// add a next / prefix
	count := strings.Count(pattern, "/") + 1

	f := func(w http.ResponseWriter, req *http.Request) {
		// 忽略路径中的前缀
		skipCount := count
		offset := strings.IndexFunc(req.URL.Path, func(r rune) bool {
			if r == '/' {
				skipCount--
			}
			if skipCount == 0 {
				return true
			}
			return false
		})
		req.URL.Path = req.URL.Path[offset:]

		// the mount route pattern
		rp := strings.TrimSuffix(RoutePattern(req), "/")
		ctx := context.WithValue(req.Context(), patternCtxKey, rp)
		req = req.WithContext(ctx)

		handler.ServeHTTP(w, req)
	}

	return http.HandlerFunc(f)
}

// Handle adds the route `pattern` that matches any http method to
// execute the `handler` http.Handler.
func (r *router) Handle(pattern string, handler http.Handler) {
	r.register("", pattern, handler)
}

// HandleFunc adds the route `pattern` that matches any http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.register("", pattern, handler)
}

// Get adds the route `pattern` that matches a POST http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) Get(pattern string, handler http.HandlerFunc) {
	r.register(http.MethodGet, pattern, handler)
}

// Post adds the route `pattern` that matches a POST http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) Post(pattern string, handler http.HandlerFunc) {
	r.register(http.MethodPost, pattern, handler)
}

// Put adds the route `pattern` that matches a Put http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) Put(pattern string, handler http.HandlerFunc) {
	r.register(http.MethodPut, pattern, handler)
}

// Patch adds the route `pattern` that matches a Patch http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) Patch(pattern string, handler http.HandlerFunc) {
	r.register(http.MethodPatch, pattern, handler)
}

// Delete adds the route `pattern` that matches a Delete http method to
// execute the `handlerFn` http.HandlerFunc.
func (r *router) Delete(pattern string, handler http.HandlerFunc) {
	r.register(http.MethodDelete, pattern, handler)
}

// NotFound sets a custom http.HandlerFunc for routing paths that could
// not be found. The default 404 handler is `http.NotFound`.
func (r *router) NotFound(h http.HandlerFunc) {
	r.notFoundHandler = h
}

// MethodNotAllowed sets a custom http.HandlerFunc for routing paths where the
// method is unresolved. The default handler returns a 405 with an empty body.
func (r *router) MethodNotAllowed(h http.HandlerFunc) {
	r.methodNotAllowedHandler = h
}

// ServeHTTP is the single method of the http.Handler interface that makes
// Mux interoperable with the standard library
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 没有自定义404和405处理函数，直接调用mux的ServeHTTP, 减少一次路由匹配
	if r.notFoundHandler == nil && r.methodNotAllowedHandler == nil {
		r.mux.ServeHTTP(w, req)
		return
	}

	h, pattern := r.mux.Handler(req)

	// handle custom 404 and 405 route
	if pattern == "" {
		h = r.serveNotFound(h)
	} else {
		h = r.mux
	}

	h.ServeHTTP(w, req)
}

func (r *router) serveNotFound(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, req *http.Request) {
		wrapper := &notFoundWrapper{
			header: http.Header{},
			buf:    bytes.NewBuffer(nil),
		}
		handler.ServeHTTP(wrapper, req)

		// 404
		if r.notFoundHandler != nil && wrapper.code == http.StatusNotFound {
			r.notFoundHandler(w, req)
			return
		}

		// 405
		if r.methodNotAllowedHandler != nil && wrapper.code == http.StatusMethodNotAllowed {
			w.Header().Set("Allow", wrapper.header.Get("Allow"))
			r.methodNotAllowedHandler(w, req)
			return
		}

		wrapper.ServeHTTP(w, req)
	}

	return http.HandlerFunc(f)
}

func (r *router) register(method string, pattern string, handler http.Handler) {
	if r.basePath != "" {
		pattern = path.Join(r.basePath, pattern)
	}

	if strings.HasSuffix(pattern, "/") {
		pattern = pattern + "{$}" // 除开mount, 其他路径都精准匹配
	}

	if method != "" {
		r.mux.Handle(method+" "+pattern, r.chain(handler))
	} else {
		r.mux.Handle(pattern, r.chain(handler))
	}
}

// chain builds a http.Handler composed of an inline middleware stack and endpoint
// handler in the order they are passed.
func (r *router) chain(endpoint http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		endpoint = r.middlewares[i](endpoint)
	}

	return endpoint
}

func (r *router) clone() *router {
	newMiddlewares := make([]func(http.Handler) http.Handler, len(r.middlewares))
	copy(newMiddlewares, r.middlewares)

	newRouter := &router{
		mux:         r.mux,
		basePath:    r.basePath,
		middlewares: newMiddlewares,
	}
	return newRouter
}

// notFoundWrapper that implements the minimal http.ResponseWriter interface.
type notFoundWrapper struct {
	code   int
	header http.Header
	buf    *bytes.Buffer
}

// Header implement the http.ResponseWriter interface Header method
func (nf *notFoundWrapper) Header() http.Header {
	return nf.header
}

// WriteHeader implement the http.ResponseWriter interface WriteHeader method
func (nf *notFoundWrapper) WriteHeader(code int) {
	nf.code = code
}

// Write implement the http.ResponseWriter interface Write method
func (nf *notFoundWrapper) Write(buf []byte) (int, error) {
	return nf.buf.Write(buf)
}

// ServeHTTP implement the http.Handler interface ServeHTTP method
func (nf *notFoundWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for k, v := range nf.header {
		for i := range v {
			w.Header().Add(k, v[i])
		}
	}

	w.WriteHeader(nf.code)
	w.Write(nf.buf.Bytes())
}

// RoutePattern returns the matched route pattern include mount prefix, but ignore method and host
func RoutePattern(req *http.Request) string {
	pattern, _ := req.Context().Value(patternCtxKey).(string)

	// ignore method and host
	i := strings.IndexByte(req.Pattern, '/')
	if i < 0 {
		return pattern
	}

	return pattern + req.Pattern[i:]
}
