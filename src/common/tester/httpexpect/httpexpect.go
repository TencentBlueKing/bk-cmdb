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
 
package httpexpect

import (
	"bytes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type Expect struct {
	t   *testing.T
	err error
}

type Request struct {
	e      *Expect
	methor string
	url    string
	body   []byte
}

type Response struct {
	*httptest.ResponseRecorder
	responseBody []byte
}

func NewExpect(t *testing.T) *Expect {
	return &Expect{t: t}
}

func (e *Expect) clone() *Expect {
	ne := *e
	return &ne
}

func (e *Expect) GET(url string) *Request {
	return &Request{e: e.clone(), methor: "GET", url: url}
}
func (e *Expect) POST(url string) *Request {
	return &Request{e: e.clone(), methor: "POST", url: url}
}
func (e *Expect) PUT(url string) *Request {
	return &Request{e: e.clone(), methor: "PUT", url: url}
}
func (e *Expect) DELETE(url string) *Request {
	return &Request{e: e.clone(), methor: "DELETE", url: url}
}

func (req *Request) WithJSON(data interface{}) *Request {
	body, err := json.Marshal(data)
	assert.NoError(req.e.t, err)
	req.body = body
	return req
}
func (req *Request) WithJSONRaw(data []byte) *Request {
	req.body = data
	return req
}

func (req *Request) Expect(f func(req *restful.Request, resp *restful.Response)) *Response {
	httptest.NewRequest(req.methor, req.url, bytes.NewBuffer(req.body))
	resp := httptest.NewRecorder()
	f(restful.NewRequest(httptest.NewRequest(req.methor, req.url, bytes.NewBuffer(req.body))),
		restful.NewResponse(resp))
	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(req.e.t, err)
	return &Response{resp, respBody}
}

func (resp *Response) BodyData() []byte {
	return resp.responseBody
}
