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
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	tenantset "configcenter/pkg/types/tenant-set"
	"configcenter/src/ac/meta"
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
		if err := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
			ErrCode: common.CCErrRewriteRequestUriFailed,
			Data:    nil,
		}); err != nil {
			blog.Errorf("response request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			return
		}
		return
	}

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			if rerr := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			}); rerr != nil {
				blog.Errorf("proxy request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, rerr, rid)
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
		servers, err = s.discovery.OperationServer().GetServers()

	case TaskType:
		servers, err = s.discovery.TaskServer().GetServers()

	case AdminType:
		servers, err = s.discovery.MigrateServer().GetServers()

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
		resp.WriteAsJson(&metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommAuthNotHavePermission,
			ErrMsg: "not allowed to call transaction api",
		})
		return
	}

	s.urlFilterChan(req, resp, chain, s.discovery.CoreService(), rootPath, "/api/v3")
}

// WebCoreFilterChan core-service api filter chan for web-server
func (s *service) WebCoreFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// right now, only allow calling from web-server
	if !httpheader.IsReqFromWeb(req.Request.Header) {
		resp.WriteAsJson(&metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommAuthNotHavePermission,
			ErrMsg: "not allowed to call this api",
		})
		return
	}

	s.urlFilterChan(req, resp, chain, s.discovery.CoreService(), rootPath, "/api/v3")
}

// TenantSetFilterChan is the filter chan for tenant set api
func (s *service) TenantSetFilterChan(discovery discovery.Interface, targetURL string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// check if the request tenant is system tenant
		tenantID := httpheader.GetTenantID(req.Request.Header)
		if tenantID != common.BKDefaultTenantID {
			resp.WriteAsJson(&metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommAuthNotHavePermission,
				ErrMsg: "only system tenant can access this api",
			})
			return
		}

		// check if tenant set id is the same with the default tenant set id
		if req.PathParameter("tenant_set_id") != strconv.FormatInt(tenantset.DefaultTenantSetID, 10) {
			resp.WriteAsJson(&metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommParamsInvalid,
				ErrMsg: "tenant set id is invalid",
			})
			return
		}

		// authorize access tenant set permission
		if auth.EnableAuthorize() {
			user := meta.UserInfo{
				UserName: httpheader.GetUser(req.Request.Header),
				TenantID: tenantID,
			}
			resource := meta.ResourceAttribute{Basic: meta.Basic{Type: meta.TenantSet, Action: meta.AccessTenantSet,
				InstanceID: tenantset.DefaultTenantSetID}}
			if authResp, authorized := s.authorizeReq(req, user, resource); !authorized {
				resp.WriteAsJson(authResp)
				return
			}
		}

		// replace url with the target url
		for key, value := range req.PathParameters() {
			targetURL = strings.Replace(targetURL, fmt.Sprintf("{%s}", key), value, 1)
		}
		req.Request.URL.Path = targetURL
		req.Request.RequestURI = targetURL

		// replace tenant id with the specified tenant id
		httpheader.SetTenantID(req.Request.Header, req.PathParameter("tenant_id"))

		s.urlFilterChan(req, resp, chain, discovery, "", "")
	}
}

// urlFilterChan url filter chan, modify the request to dispatch it to specific sever
func (s *service) urlFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain,
	discovery discovery.Interface, prevRoot, root string) {

	rid := httpheader.GetRid(req.Request.Header)

	var err error

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			respErr := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			})
			if respErr != nil {
				blog.Errorf("proxy request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, respErr, rid)
				return
			}
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

	return s.authorizeReq(req, attribute.User, attribute.Resources...)
}

func (s *service) authorizeReq(req *restful.Request, user meta.UserInfo, resources ...meta.ResourceAttribute) (
	*metadata.BaseResp, bool) {

	rid := httpheader.GetRid(req.Request.Header)
	path := req.Request.URL.Path
	errf := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(req.Request.Header))

	decisions, err := s.authorizer.AuthorizeBatch(req.Request.Context(), req.Request.Header, user, resources...)
	if err != nil {
		blog.Errorf("authorized request failed, url: %s, err: %v, rid: %s", path, err, rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommCheckAuthorizeFailed,
			ErrMsg: errf.Error(common.CCErrCommCheckAuthorizeFailed).Error(),
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

	if authorized {
		return nil, true
	}

	s.noPermissionRequestTotal.With(
		prometheus.Labels{
			metrics.LabelHandler: path,
			metrics.LabelAppCode: httpheader.GetAppCode(req.Request.Header),
		},
	).Inc()

	permission, err := s.authorizer.GetPermissionToApply(req.Request.Context(), req.Request.Header, resources)
	if err != nil {
		blog.Errorf("get permission to apply failed, err: %v, rid: %s", err, rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommCheckAuthorizeFailed,
			ErrMsg: errf.Error(common.CCErrCommCheckAuthorizeFailed).Error(),
			Result: false,
		}, false
	}

	blog.Warnf("request has no permission, url: %s, user: %+v, resources: %+v, permission: %+v, rid: %s", path,
		user, resources, permission, rid)

	return &metadata.BaseResp{
		Code:        common.CCNoPermission,
		ErrMsg:      errf.Error(common.CCErrCommAuthNotHavePermission).Error(),
		Result:      false,
		Permissions: permission,
	}, false
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
			rsp := metadata.BaseResp{
				Code:   common.CCErrTooManyRequestErr,
				ErrMsg: "too many requests",
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		key := common.ApiCacheLimiterRulePrefix + rule.RuleName
		result, err := s.cache.Eval(context.Background(), setRequestCntTTLScript, []string{key}, rule.TTL).Result()
		if err != nil {
			blog.Errorf("redis Eval failed, key:%s, rule:%#v, err: %v, rid: %s", key, *rule, err, rid)
			fchain.ProcessFilter(req, resp)
			return
		}
		cnt, ok := result.(int64)
		if !ok {
			blog.Errorf("execute setRequestCntTTLScript failed, key:%s, rule:%#v, err: %v, rid: %s",
				key, *rule, result, rid)
			fchain.ProcessFilter(req, resp)
			return
		}

		if cnt > rule.Limit {
			blog.Errorf("too many requests, matched rule is %#v, rid: %s", *rule, rid)
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

// JwtFilter the filter that handles the source of the jwt request
func (s *service) JwtFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		header, err := jwt.GetHandler().Parse(req.Request.Header)
		if err != nil {
			rsp := metadata.BaseResp{
				Code:   common.CCErrAPINoPassSourceCertification,
				ErrMsg: err.Error(),
				Result: false,
			}
			_ = resp.WriteAsJson(rsp)
			return
		}
		req.Request.Header = header

		fchain.ProcessFilter(req, resp)
		return
	}
}

// TenantVerify the filter that handles tenant verification
func (s *service) TenantVerify() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {

		if !s.config.EnableMultiTenantMode {
			if tenant := req.Request.Header.Get(httpheader.TenantHeader); tenant == "" ||
				tenant == common.BKUnconfiguredTenantID {
				req.Request.Header.Set(httpheader.TenantHeader, common.BKUnconfiguredTenantID)
				fchain.ProcessFilter(req, resp)
				return
			}
			blog.Errorf("tenant mode is not enabled but tenant is set, rid: %s", httpheader.GetRid(req.Request.Header))
			rsp := metadata.BaseResp{
				Code:   common.CCErrAPICheckTenantInvalid,
				ErrMsg: "tenant mode is not enabled but tenant is set",
				Result: false,
			}
			_ = resp.WriteAsJson(rsp)
			return
		}

		tenantID := httpheader.GetTenantID(req.Request.Header)
		tenantData, exist := tenant.GetTenant(tenantID)
		if !exist || tenantData.Status != types.EnabledStatus {
			blog.Errorf("invalid tenant: %s, rid: %s", tenantID, httpheader.GetRid(req.Request.Header))
			rsp := metadata.BaseResp{
				Code:   common.CCErrAPICheckTenantInvalid,
				ErrMsg: "invalid tenant",
				Result: false,
			}
			_ = resp.WriteAsJson(rsp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}
