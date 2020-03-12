/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

func (s *service) Get(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

func (s *service) Put(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

func (s *service) Post(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

func (s *service) Delete(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

func (s *service) Do(req *restful.Request, resp *restful.Response) {
	start := time.Now()
	url := fmt.Sprintf("%s://%s%s", req.Request.URL.Scheme, req.Request.URL.Host, req.Request.RequestURI)
	proxyReq, err := http.NewRequest(req.Request.Method, url, req.Request.Body)
	if err != nil {
		blog.Errorf("new proxy request[%s] failed, err: %v", url, err)
		if err := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("proxy request failed, %s", err.Error()),
			ErrCode: common.CCErrProxyRequestFailed,
			Data:    nil,
		}); err != nil {
			blog.Errorf("response request[url: %s] failed, err: %v", req.Request.RequestURI, err)
		}
		return
	}

	for k, v := range req.Request.Header {
		if len(v) > 0 {
			proxyReq.Header.Set(k, v[0])
		}
	}

	response, err := s.client.Do(proxyReq)
	if err != nil {
		blog.Errorf("*failed do request[%s url: %s] , err: %v", req.Request.Method, url, err)

		if err := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("proxy request failed, %s", err.Error()),
			ErrCode: common.CCErrProxyRequestFailed,
			Data:    nil,
		}); err != nil {
			blog.Errorf("response request[%s url: %s] failed, err: %v", req.Request.Method, url, err)
		}
		return
	}
	blog.V(5).Infof("success [%s] do request[%s url: %s]  ", response.Status, req.Request.Method, url)

	defer response.Body.Close()

	for k, v := range response.Header {
		if len(v) > 0 {
			resp.Header().Set(k, v[0])
		}
	}

	resp.ResponseWriter.WriteHeader(response.StatusCode)

	if _, err := io.Copy(resp, response.Body); err != nil {
		blog.Errorf("response request[url: %s] failed, err: %v", req.Request.RequestURI, err)
		return
	}
	blog.V(4).Infof("request id: %s, cost: %dms, action: %s, status code: %d, url: %s, user: %s, app code: %s",
		req.Request.Header.Get(common.BKHTTPCCRequestID),
		time.Since(start).Nanoseconds()/int64(time.Millisecond),
		req.Request.Method, response.StatusCode, url,
		req.Request.Header.Get(common.BKHTTPHeaderUser),
		req.Request.Header.Get(common.BKHTTPRequestAppCode))
	return
}
