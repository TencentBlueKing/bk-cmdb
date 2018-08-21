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

package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/blog"
)

// http request verb type
type VerbType string

const (
	PUT    VerbType = http.MethodPut
	POST   VerbType = http.MethodPost
	GET    VerbType = http.MethodGet
	DELETE VerbType = http.MethodDelete
	PATCH  VerbType = http.MethodPatch
)

type Request struct {
	capability *util.Capability

	verb    VerbType
	params  url.Values
	headers http.Header
	body    io.Reader
	ctx     context.Context

	// prefixed url
	baseURL string
	// sub path of the url, will be append to baseURL
	subPath string

	// request timeout value
	timeout time.Duration

	err error
}

func (r *Request) WithParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (r *Request) WithHeaders(header http.Header) *Request {
	if r.headers == nil {
		r.headers = header
		return r
	}

	for key, values := range header {
		for _, v := range values {
			r.headers.Add(key, v)
		}
	}
	return r
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) WithTimeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

func (r *Request) SubResource(subPath string) *Request {
	subPath = strings.TrimLeft(subPath, "/")
	r.subPath = subPath
	return r
}

func (r *Request) Body(body interface{}) *Request {
	if nil == body {
		r.body = bytes.NewReader([]byte(""))
		return r
	}

	valueOf := reflect.ValueOf(body)
	switch valueOf.Kind() {
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		if valueOf.IsNil() {
			r.body = bytes.NewReader([]byte(""))
			return r
		}
		break

	case reflect.Struct:
		break

	default:
		r.err = errors.New("body should be one of interface, map, pointer or slice value")
		r.body = bytes.NewReader([]byte(""))
		return r
	}

	data, err := json.Marshal(body)
	if nil != err {
		r.err = err
		r.body = bytes.NewReader([]byte(""))
		return r
	}

	r.body = bytes.NewReader(data)
	return r
}

func (r *Request) WrapURL() *url.URL {
	finalUrl := &url.URL{}
	if len(r.baseURL) != 0 {
		u, err := url.Parse(r.baseURL)
		if err != nil {
			r.err = err
			return new(url.URL)
		}
		*finalUrl = *u
	}

	finalUrl.Path = finalUrl.Path + r.subPath

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	if r.timeout != 0 {
		query.Set("timeout", r.timeout.String())
	}

	finalUrl.RawQuery = query.Encode()
	return finalUrl
}

func (r *Request) Do() *Result {
	result := new(Result)
	if r.err != nil {
		result.Err = r.err
		return result
	}

	client := r.capability.Client
	if client == nil {
		client = http.DefaultClient
	}

	maxRetryCycle := 3
	retries := 0

	hosts, err := r.capability.Discover.GetServers()
	if err != nil {
		result.Err = err
		return result
	}

	for try := 0; try < maxRetryCycle; try++ {
		for index, host := range hosts {
			retries = try + index
			url := host + r.WrapURL().String()
			req, err := http.NewRequest(string(r.verb), url, r.body)
			if err != nil {
				result.Err = err
				return result
			}

			if r.ctx != nil {
				req.WithContext(r.ctx)
			}

			req.Header = r.headers
			if len(req.Header) == 0 {
				req.Header = make(http.Header)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			if retries > 0 {
				r.tryThrottle(url)
			}

			resp, err := client.Do(req)
			if err != nil {
				// "Connection reset by peer" is a special err which in most scenario is a a transient error.
				// Which means that we can retry it. And so does the GET operation.
				// While the other "write" operation can not simply retry it again, because they are not idempotent.

				if !isConnectionReset(err) || r.verb != GET {
					result.Err = err
					return result
				}

				// retry now
				time.Sleep(20 * time.Millisecond)
				continue

			}

			var body []byte
			if resp.Body != nil {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					if err == io.ErrUnexpectedEOF {
						// retry now
						time.Sleep(20 * time.Millisecond)
						continue
					}
					result.Err = err
					return result
				}
				body = data
			}
			result.Body = body
			result.StatusCode = resp.StatusCode
			return result
		}

	}

	result.Err = errors.New("unexpected error")
	return result
}

const maxLatency = 100 * time.Millisecond

func (r *Request) tryThrottle(url string) {
	now := time.Now()
	if r.capability.Throttle != nil {
		r.capability.Throttle.Accept()
	}

	if latency := time.Since(now); latency > maxLatency {
		blog.V(3).Infof("Throttling request took %d ms, request: %s", latency, r.verb, url)
	}
}

type Result struct {
	Body       []byte
	Err        error
	StatusCode int
}

func (r *Result) Into(obj interface{}) error {
	if nil != r.Err {
		return r.Err
	}

	if 0 != len(r.Body) {
		err := json.Unmarshal(r.Body, obj)
		if nil != err {
			if http.StatusOK != r.StatusCode {
				return fmt.Errorf("error info %s", string(r.Body))
			}
			blog.Errorf("http reply not json, reply:%s, error:%s", string(r.Body), err.Error())
			return err
		}
	}
	return nil
}

// Returns if the given err is "connection reset by peer" error.
func isConnectionReset(err error) bool {
	if urlErr, ok := err.(*url.Error); ok {
		err = urlErr.Err
	}
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	if osErr, ok := err.(*os.SyscallError); ok {
		err = osErr.Err
	}
	if errno, ok := err.(syscall.Errno); ok && errno == syscall.ECONNRESET {
		return true
	}
	return false
}
