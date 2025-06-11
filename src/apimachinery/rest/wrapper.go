/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

import "net/http"

// clientWrapper is a wrapper for restful client
type clientWrapper struct {
	client   ClientInterface
	wrappers RequestWrapperChain
}

// NewClientWrapper new restful client wrapper
func NewClientWrapper(client ClientInterface, wrappers ...RequestWrapper) ClientInterface {
	return &clientWrapper{
		client:   client,
		wrappers: wrappers,
	}
}

// Verb generate restful request by http method
func (r *clientWrapper) Verb(verb VerbType) *Request {
	return ProcessRequestWrapperChain(r.client.Verb(verb), r.wrappers)
}

// Post generate post restful request
func (r *clientWrapper) Post() *Request {
	return ProcessRequestWrapperChain(r.client.Post(), r.wrappers)
}

// Put generate put restful request
func (r *clientWrapper) Put() *Request {
	return ProcessRequestWrapperChain(r.client.Put(), r.wrappers)
}

// Get generate get restful request
func (r *clientWrapper) Get() *Request {
	return ProcessRequestWrapperChain(r.client.Get(), r.wrappers)
}

// Delete generate delete restful request
func (r *clientWrapper) Delete() *Request {
	return ProcessRequestWrapperChain(r.client.Delete(), r.wrappers)
}

// Patch generate patch restful request
func (r *clientWrapper) Patch() *Request {
	return ProcessRequestWrapperChain(r.client.Patch(), r.wrappers)
}

// RequestWrapper is the restful request wrapper
type RequestWrapper func(*Request) *Request

// RequestWrapperChain is the restful request wrapper chain
type RequestWrapperChain []RequestWrapper

// ProcessRequestWrapperChain process restful request wrapper chain
func ProcessRequestWrapperChain(req *Request, wrappers RequestWrapperChain) *Request {
	req.wrappers = wrappers
	return req
}

// BaseUrlWrapper returns a restful request wrapper that changes request's base url
func BaseUrlWrapper(baseUrl string) RequestWrapper {
	return func(request *Request) *Request {
		request.baseURL = baseUrl
		return request
	}
}

// HeaderWrapper returns a restful request wrapper that changes request's header
func HeaderWrapper(handler func(http.Header) http.Header) RequestWrapper {
	return func(request *Request) *Request {
		if request.headers == nil {
			request.headers = make(http.Header)
		}
		request.headers = handler(request.headers)
		return request
	}
}
