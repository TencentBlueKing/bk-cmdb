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
	"strings"

	"configcenter/src/ac/parser"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"

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
	TaskType        RequestType = "task"
	AdminType       RequestType = "admin"
	CloudType       RequestType = "cloud"
	CacheType       RequestType = "cache"
)

const publicKey = "public_key"

// URLFilterChan url filter chan
func (s *service) URLFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	rid := util.GetHTTPCCRequestID(req.Request.Header)

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

func (s *service) authFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response,
	fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		rid := util.GetHTTPCCRequestID(req.Request.Header)
		path := req.Request.URL.Path

		blog.V(7).Infof("authFilter on url: %s, rid: %s", path, rid)
		if !auth.EnableAuthorize() {
			blog.V(7).Infof("auth disabled, skip auth filter, rid: %s", rid)
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/verify" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/business_list" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/skip_url" {
			fchain.ProcessFilter(req, resp)
			return
		}

		language := util.GetLanguage(req.Request.Header)
		attribute, err := parser.ParseAttribute(req, s.engine)
		if err != nil {
			blog.Errorf("authFilter failed, caller: %s, parse auth attribute for %s %s failed, err: %v, rid: %s",
				req.Request.RemoteAddr, req.Request.Method, req.Request.URL.Path, err, rid)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommParseAuthAttributeFailed,
				ErrMsg: err.Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		if blog.V(5) {
			blog.InfoJSON("auth filter parsed attribute: %s, rid: %s", attribute, rid)
		}

		decisions, err := s.authorizer.AuthorizeBatch(req.Request.Context(), req.Request.Header, attribute.User,
			attribute.Resources...)
		if err != nil {
			blog.Errorf("authFilter failed, authorized request failed, url: %s, err: %v, rid: %s", path, err, rid)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		authorized := true
		for _, decision := range decisions {
			if !decision.Authorized {
				authorized = false
				break
			}
		}

		if !authorized {
			permission, err := s.authorizer.GetPermissionToApply(req.Request.Context(), req.Request.Header,
				attribute.Resources)
			if err != nil {
				blog.Errorf("get permission to apply failed, err: %v, rid: %s", err, rid)
				rsp := metadata.BaseResp{
					Code: common.CCErrCommCheckAuthorizeFailed,
					ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).
						Error(),
					Result: false,
				}
				resp.WriteAsJson(rsp)
				return
			}
			blog.WarnJSON("authFilter failed, url: %s, attribute: %s, permission: %s, rid: %s",
				path, attribute, permission, rid)
			rsp := metadata.BaseResp{
				Code: common.CCNoPermission,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommAuthNotHavePermission).
					Error(),
				Result:      false,
				Permissions: permission,
			}
			resp.WriteAsJson(rsp)
			return
		}

		blog.V(7).Infof("authFilter authorize on url:%s success, rid: %s", path, rid)
		fchain.ProcessFilter(req, resp)
		return
	}
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
		rid := util.GetHTTPCCRequestID(req.Request.Header)
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

// SourceFilter determine the source of the request and perform related operations
func (s *service) SourceFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		header := req.Request.Header
		rid := util.GetHTTPCCRequestID(header)

		loginVersion, _ := cc.String("webServer.login.version")
		if loginVersion == common.BKOpenSourceLoginPluginVersion || loginVersion == common.BKSkipLoginPluginVersion {
			fchain.ProcessFilter(req, resp)
			return
		}

		// it means request by apiGateway
		if util.GetGatewayName(header) != "" {
			if err := s.doRequestFromAPIGW(req); err != nil {
				rsp := metadata.BaseResp{
					Code:   common.CCErrAPINoPassSourceCertification,
					ErrMsg: err.Error(),
					Result: false,
				}
				resp.WriteAsJson(rsp)
				return
			}
			fchain.ProcessFilter(req, resp)
			return
		}

		webToken := util.GetWebToken(header)
		// it means request by webserver
		if webToken != "" {
			webTokenConfig, _ := cc.String("webServer.webToken")
			if webTokenConfig != webToken {
				blog.Errorf("handle request from webServer error, rid: %s", rid)
				rsp := metadata.BaseResp{
					Code:   common.CCErrAPINoPassSourceCertification,
					ErrMsg: "handle request from webServer error",
					Result: false,
				}
				resp.WriteAsJson(rsp)
				return
			}
			fchain.ProcessFilter(req, resp)
			return
		}

		// it means request by esb
		if err := s.doRequestFromEsb(req); err != nil {
			rsp := metadata.BaseResp{
				Code:   common.CCErrAPINoPassSourceCertification,
				ErrMsg: err.Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

func (s *service) doRequestFromEsb(req *restful.Request) error {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)

	if s.GetEsbPublicKey() == "" {
		esbResponse, err := esb.EsbClient().EsbSrv().GetApiPublicKey(req.Request.Context(), header)
		if err != nil {
			blog.Errorf("handle request from esb error, err: %v, rid: %s", err, rid)
			return err
		}
		key, ok := esbResponse.Data[publicKey].(string)
		if !ok {
			blog.Errorf("handle request from esb error, can not transfer publicKey %v to string type,"+
				" rid: %s", esbResponse.Data[publicKey], rid)
			return fmt.Errorf("can not transfer publicKey: %v to string type", esbResponse.Data[publicKey])
		}
		s.SetEsbPublicKey(key)
	}
	jwtToken := util.GetJWTToken(header)
	esbPublicKey := s.GetEsbPublicKey()
	_, err := ParseToken(jwtToken, esbPublicKey)
	if err != nil {
		blog.Errorf("handle request from esb error, jwtToken: %s, esbPublicKey: %s, err: %v, rid: %s",
			jwtToken, esbPublicKey, err, rid)
		return err
	}
	return nil
}

func (s *service) doRequestFromAPIGW(req *restful.Request) error {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)

	key, err := cc.String("apiServer.apiGateway.publicKey")
	if err != nil {
		blog.Errorf("handle request from API Gateway error, err: %v, rid: %s", err, rid)
		return err
	}
	token := util.GetJWTToken(header)
	jwtSecret := "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
	parseToken, err := ParseToken(token, jwtSecret)
	if err != nil {
		blog.Errorf("handle request from API Gateway error, token: %s, publicKey: %s, err: %v, rid: %s",
			token, jwtSecret, err, rid)
		return err
	}

	if parseToken.User.UserName == "" {
		blog.Errorf("handle request from API Gateway error, no setting bk_username, token: %s, "+
			"publicKey: %s, rid: %s", token, jwtSecret, rid)
		return fmt.Errorf("no setting bk_username")
	}
	if !util.SetLanguageFromApiGW(header) {
		blog.Errorf("handle request from API Gateway error, no setting HTTP-BLUEKING-LANGUAGE, "+
			"rid: %s", rid)
		return fmt.Errorf("no setting HTTP-BLUEKING-LANGUAGE")
	}
	if !util.SetOwnerIDFromApiGW(header) {
		blog.Errorf("handle request from API Gateway error, no setting HTTP-BLUEKING-SUPPLIER-ID,"+
			" rid: %s", rid)
		return fmt.Errorf("no setting HTTP-BLUEKING-SUPPLIER-ID")
	}

	header.Set(common.BKHTTPHeaderUser, parseToken.User.UserName)
	header.Set(common.BKHTTPRequestAppCode, parseToken.App.AppCode)
	return nil
}
