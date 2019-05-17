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
	"strings"

	"configcenter/src/apimachinery/util"
)

type ClientInterface interface {
	Verb(verb VerbType) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
	Patch() *Request
}

func NewRESTClient(c *util.Capability, baseUrl string) ClientInterface {
	baseUrl = strings.Trim(baseUrl, "/")
	baseUrl = "/" + baseUrl + "/"
	return &RESTClient{
		baseUrl:    baseUrl,
		capability: c,
	}
}

type RESTClient struct {
	baseUrl    string
	capability *util.Capability
}

func (r *RESTClient) Verb(verb VerbType) *Request {
	return &Request{
		verb:       verb,
		baseURL:    r.baseUrl,
		capability: r.capability,
	}
}

func (r *RESTClient) Post() *Request {
	return r.Verb(POST)
}

func (r *RESTClient) Put() *Request {
	return r.Verb(PUT)
}

func (r *RESTClient) Get() *Request {
	return r.Verb(GET)
}

func (r *RESTClient) Delete() *Request {
	return r.Verb(DELETE)
}

func (r *RESTClient) Patch() *Request {
	return r.Verb(PATCH)
}
