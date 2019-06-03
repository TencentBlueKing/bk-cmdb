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

package service

import (
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

type RequestType string

const (
	UnknownType     RequestType = "unknown"
	TopoType        RequestType = "topo"
	HostType        RequestType = "host"
	ProcType        RequestType = "proc"
	EventType       RequestType = "event"
	DataCollectType RequestType = "collect"
	OperationType   RequestType = "operation"
)

func (s *service) URLFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	var kind RequestType
	var err error
	kind, err = URLPath(req.Request.RequestURI).FilterChain(req)
	if err != nil {
		blog.Errorf("rewrite request url[%s] failed, err: %v", req.Request.RequestURI, err)
		if err := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
			ErrCode: common.CCErrRewriteRequestUriFailed,
			Data:    nil,
		}); err != nil {
			blog.Errorf("response request[url: %s] failed, err: %v", req.Request.RequestURI, err)
			return
		}
		return
	}

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v", req.Request.RequestURI, err)
			if rerr := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			}); rerr != nil {
				blog.Errorf("proxy request[url: %s] failed, err: %v", req.Request.RequestURI, rerr)
				return
			}
			return
		}
	}()

	servers := make([]string, 0)
	switch kind {
	case TopoType:
		servers, err = s.discovery.TopoServer().GetServers()

	case ProcType:
		servers, err = s.discovery.ProcServer().GetServers()

	case EventType:
		servers, err = s.discovery.EventServer().GetServers()

	case HostType:
		servers, err = s.discovery.HostServer().GetServers()

	case DataCollectType:
		servers, err = s.discovery.DataCollect().GetServers()

	case OperationType:
		servers, err = s.discovery.DataCollect().GetServers()
	}

	if err != nil {
		return
	}

	if strings.HasPrefix(servers[0], "https://") {
		req.Request.URL.Host = servers[0][8:]
		req.Request.URL.Scheme = "https"
	} else {
		req.Request.URL.Host = servers[0][7:]
		req.Request.URL.Scheme = "http"
	}

	chain.ProcessFilter(req, resp)
}
