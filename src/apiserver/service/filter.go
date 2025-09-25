/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/ac/parser"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metrics"
	"configcenter/src/common/resource/jwt"

	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus"
)

// RequestType TODO
type RequestType string

const (
	// UnknownType TODO
	UnknownType RequestType = "unknown"
	// TopoType TODO
	TopoType RequestType = "topo"
	// HostType TODO
	HostType RequestType = "host"
	// ProcType TODO
	ProcType RequestType = "proc"
	// EventType TODO
	EventType RequestType = "event"
	// DataCollectType TODO
	DataCollectType RequestType = "collect"
	// OperationType TODO
	OperationType RequestType = "operation"
	// TaskType TODO
	TaskType RequestType = "task"
	// AdminType TODO
	AdminType RequestType = "admin"
	// CloudType TODO
	CloudType RequestType = "cloud"
	// CacheType TODO
	CacheType RequestType = "cache"
)

// URLFilterChan url filter chan
func (s *service) URLFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	rid := httpheader.GetRid(req.Request.Header)

	var kind RequestType
	var err error
	kind, err = URLPath(req.Request.RequestURI).FilterChain(req)
	if err != nil {
		blog.Errorf("rewrite request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)

		s.RespError(req, resp, http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
			ErrCode: common.CCErrRewriteRequestUriFailed,
			Data:    nil,
		})
		return
	}

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)

			s.RespError(req, resp, http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			})
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
		servers, err = s.discovery.OperationServer().GetServers()

	case TaskType:
		servers, err = s.discovery.TaskServer().GetServers()

	case AdminType:
		servers, err = s.discovery.MigrateServer().GetServers()

	case CloudType:
		servers, err = s.discovery.CloudServer().GetServers()

	case CacheType:
		servers, err = s.discovery.CacheService().GetServers()

	default:
		name := string(kind)
		if name != "" {
			servers, err = s.discovery.Server(name).GetServers()
		}
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

// BizFilterChan biz api filter chan
func (s *service) BizFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	s.urlFilterChan(req, resp, chain, s.discovery.TopoServer(), rootPath+"/biz", "/topo/v3/app")
}

// HostFilterChan host server api filter chan
func (s *service) HostFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	s.urlFilterChan(req, resp, chain, s.discovery.HostServer(), rootPath, "/host/v3")
}

// TopoFilterChan topo server api filter chan
func (s *service) TopoFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	s.urlFilterChan(req, resp, chain, s.discovery.TopoServer(), rootPath, "/topo/v3")
}

// CacheFilterChan cache service api filter chan
func (s *service) CacheFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	s.urlFilterChan(req, resp, chain, s.discovery.CacheService(), rootPath+"/cache", "/cache/v3")
}

// TxnFilterChan transaction api filter chan
func (s *service) TxnFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// right now, only allow calling from web-server
	if !httpheader.IsReqFromWeb(req.Request.Header) {
		s.RespError(req, resp, http.StatusOK, &metadata.RespError{
			ErrCode: common.CCErrCommAuthNotHavePermission,
			Msg:     fmt.Errorf("not allowed to call transaction api"),
		})
		return
	}

	s.urlFilterChan(req, resp, chain, s.discovery.CoreService(), rootPath, "/api/v3")
}

// WebCoreFilterChan core-service api filter chan for web-server
func (s *service) WebCoreFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// right now, only allow calling from web-server
	if !httpheader.IsReqFromWeb(req.Request.Header) {
		s.RespError(req, resp, http.StatusOK, &metadata.RespError{
			ErrCode: common.CCErrCommAuthNotHavePermission,
			Msg:     fmt.Errorf("not allowed to call this api"),
		})
		return
	}

	s.urlFilterChan(req, resp, chain, s.discovery.CoreService(), rootPath, "/api/v3")
}

// urlFilterChan url filter chan, modify the request to dispatch it to specific sever
func (s *service) urlFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain,
	discovery discovery.Interface, prevRoot, root string) {

	rid := httpheader.GetRid(req.Request.Header)

	var err error

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			s.RespError(req, resp, http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			})
			return
		}
	}()

	servers, err := discovery.GetServers()
	if err != nil {
		return
	}

	// set the request server address through discovery
	if strings.HasPrefix(servers[0], "https://") {
		req.Request.URL.Host = servers[0][8:]
		req.Request.URL.Scheme = "https"
	} else {
		req.Request.URL.Host = servers[0][7:]
		req.Request.URL.Scheme = "http"
	}

	// set the request server root path instead of previous api serverS root path
	req.Request.RequestURI = root + strings.TrimPrefix(req.Request.RequestURI, prevRoot)
	req.Request.URL.Path = root + strings.TrimPrefix(req.Request.URL.Path, prevRoot)

	chain.ProcessFilter(req, resp)
}

func (s *service) authFilter(errFunc func() errors.CCErrorIf) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		if !auth.EnableAuthorize() {
			fchain.ProcessFilter(req, resp)
			return
		}

		rsp, success := s.verifyAuthorizeStatus(req, errFunc)
		if !success {
			_ = resp.WriteAsJson(rsp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

func (s *service) verifyAuthorizeStatus(req *restful.Request, errFunc func() errors.CCErrorIf) (*metadata.BaseResp,
	bool) {

	rid := httpheader.GetRid(req.Request.Header)
	path := req.Request.URL.Path
	language := httpheader.GetLanguage(req.Request.Header)
	attribute, err := parser.ParseAttribute(req, s.engine)
	if err != nil {
		blog.Errorf("authFilter failed, caller: %s, parse auth attribute for %s %s failed, err: %v, rid: %s",
			req.Request.RemoteAddr, req.Request.Method, req.Request.URL.Path, err, rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommParseAuthAttributeFailed,
			ErrMsg: err.Error(),
			Result: false,
		}, false
	}

	if blog.V(5) {
		blog.InfoJSON("auth filter parsed attribute: %s, rid: %s", attribute, rid)
	}

	decisions, err := s.authorizer.AuthorizeBatch(req.Request.Context(), req.Request.Header, attribute.User,
		attribute.Resources...)
	if err != nil {
		blog.Errorf("authFilter failed, authorized request failed, url: %s, err: %v, rid: %s",
			path, err, rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommCheckAuthorizeFailed,
			ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
			Result: false,
		}, false
	}

	authorized := true
	for _, decision := range decisions {
		if !decision.Authorized {
			authorized = false
			break
		}
	}

	if !authorized {
		s.noPermissionRequestTotal.With(
			prometheus.Labels{
				metrics.LabelHandler: path,
				metrics.LabelAppCode: httpheader.GetAppCode(req.Request.Header),
			},
		).Inc()

		permission, err := s.authorizer.GetPermissionToApply(req.Request.Context(), req.Request.Header,
			attribute.Resources)
		if err != nil {
			blog.Errorf("get permission to apply failed, err: %v, rid: %s", err, rid)
			return &metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}, false
		}

		blog.WarnJSON("authFilter failed, url: %s, attribute: %s, permission: %s, rid: %s", path, attribute,
			permission, rid)

		return &metadata.BaseResp{
			Code:        common.CCNoPermission,
			ErrMsg:      errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommAuthNotHavePermission).Error(),
			Result:      false,
			Permissions: permission,
		}, false
	}

	return nil, true
}

// KEYS[1] is the redis key to incr and expire
// ARGV[1] is the ttl
const setRequestCntTTLScript = `
local cnt = redis.pcall('INCR', KEYS[1]);
if type(cnt) ~= "number"
then
	return cnt
end

local rs = redis.pcall('TTL', KEYS[1]);
if type(rs) ~= "number"
then
	return rs
end

if rs == -1
then
	rs = redis.pcall('EXPIRE', KEYS[1], ARGV[1]);
	if type(rs) ~= "number"
	then
		return rs
	end
end

return cnt
`

// LimiterFilter limit on a api request according to limiter rules
func (s *service) LimiterFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		rid := httpheader.GetRid(req.Request.Header)
		if s.limiter.LenOfRules() == 0 {
			fchain.ProcessFilter(req, resp)
			return
		}

		rule := s.limiter.GetMatchedRule(req)
		if rule == nil {
			fchain.ProcessFilter(req, resp)
			return
		}

		if rule.DenyAll {
			blog.Errorf("too many requests, matched rule is %#v, rid: %s", *rule, rid)

			s.errorLimiterTotal.With(prometheus.Labels{
				metrics.LabelAppCode: httpheader.GetAppCode(req.Request.Header),
				metrics.LabelHandler: req.Request.RequestURI,
			}).Inc()

			rsp := metadata.BaseResp{
				Code:   common.CCErrTooManyRequestErr,
				ErrMsg: "too many requests",
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		exceedLimit, err := s.checkIfExceedLimit(rule, rid)
		if err != nil {
			fchain.ProcessFilter(req, resp)
			return
		}

		if exceedLimit {
			blog.Errorf("too many requests, matched rule is %#v, rid: %s", *rule, rid)
			s.errorLimiterTotal.With(prometheus.Labels{
				metrics.LabelAppCode: httpheader.GetAppCode(req.Request.Header),
				metrics.LabelHandler: req.Request.RequestURI,
			}).Inc()

			rsp := metadata.BaseResp{
				Code:   common.CCErrTooManyRequestErr,
				ErrMsg: "too many requests",
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

// checkIfExceedLimit checks if the request exceeds the rate limit, return true if request should be rejected
func (s *service) checkIfExceedLimit(rule *metadata.LimiterRule, rid string) (bool, error) {

	key := common.ApiCacheLimiterRulePrefix + rule.RuleName
	result, err := s.cache.Eval(context.Background(), setRequestCntTTLScript, []string{key}, rule.TTL).Result()
	if err != nil {
		blog.Errorf("redis Eval failed, key:%s, rule:%#v, err: %v, rid: %s", key, *rule, err, rid)
		return false, err
	}

	cnt, ok := result.(int64)
	if !ok {
		blog.Errorf("execute setRequestCntTTLScript failed, key:%s, rule:%#v, result: %v, rid: %s",
			key, *rule, result, rid)
		return false, fmt.Errorf("execute setRequestCntTTLScript failed, result: %v", result)
	}

	if cnt > rule.Limit {
		return true, nil
	}
	return false, nil
}

// JwtFilter the filter that handles the source of the jwt request
func (s *service) JwtFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		header, err := jwt.GetHandler().Parse(req.Request.Header)
		if err != nil {
			s.RespError(req, resp, http.StatusOK, &metadata.RespError{
				Msg:     err,
				ErrCode: common.CCErrAPINoPassSourceCertification,
			})
			return
		}
		req.Request.Header = header
		fchain.ProcessFilter(req, resp)
		return
	}
}
