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

// Package service TODO
package service

import (
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/ginservice"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/logplatform/opentelemetry"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/capability"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"
	"configcenter/src/web_server/middleware"
	"configcenter/src/web_server/service/excel"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Service TODO
type Service struct {
	*options.ServerOption
	Engine   *backbone.Engine
	CacheCli redis.Client
	*logics.Logics
	Config  *options.Config
	Session redis.RedisStore
}

// WebService TODO
func (s *Service) WebService() *gin.Engine {
	setGinMode()
	ws := gin.New()
	ws.Use(gin.Logger())

	ws.Use(middleware.RequestIDMiddleware)
	ws.Use(sessions.Sessions(s.Config.Session.Name, s.Session))
	ws.Use(middleware.ValidLogin(*s.Config, s.Discovery()))
	ws.Use(func(c *gin.Context) {
		defer func() {
			// suppresses logging of a stack when err is ErrAbortHandler, same as net/http
			if err := recover(); err != nil {
				if err != http.ErrAbortHandler {
					stack := make([]byte, 10000)
					nbytes := runtime.Stack(stack, false)
					if nbytes < len(stack) {
						stack = stack[:nbytes]
					}
					request, _ := httputil.DumpRequest(c.Request, false)
					blog.Errorf("[Recovery] recovered:\n%s\n%s\n%s", string(request), err, string(stack))
				}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	})

	opentelemetry.UseOtlpMiddleware(ws)

	middleware.Engine = s.Engine

	ws.Static("/static", s.Config.Site.HtmlRoot)
	ws.LoadHTMLFiles(s.Config.Site.HtmlRoot+"/index.html", s.Config.Site.HtmlRoot+"/login.html",
		s.Config.Site.HtmlRoot+"/"+webCommon.InaccessibleHtml)

	ws.GET("/hosts/:bk_host_id/listen_ip_options", s.ListenIPOptions)
	ws.POST("/logout", s.LogOutUser)
	ws.GET("/login", s.Login)
	ws.POST("/login", s.LoginUser)
	ws.POST("/object/exportmany", s.BatchExportObject)
	ws.POST("/object/importmany/analysis", s.BatchImportObjectAnalysis)
	ws.POST("/object/importmany", s.BatchImportObject)
	ws.GET("/user/list", s.GetUserList)
	// suggest move to  Organization
	ws.GET("/user/department", s.GetDepartment)
	ws.GET("/user/departmentprofile", s.GetDepartmentProfile)

	ws.GET("/organization/department", s.GetDepartment)
	ws.GET("/organization/departmentprofile", s.GetDepartmentProfile)

	ws.GET("/user/language/:language", s.UpdateUserLanguage)
	// get current login user info
	ws.GET("/userinfo", s.UserInfo)
	ws.PUT("/user/current/supplier/:id", s.UpdateSupplier)
	ws.POST("/biz/search/web", s.SearchBusiness)

	ws.GET("/", s.Index)

	ws.POST("/object/count", s.GetObjectInstanceCount)

	ws.POST("/regular/verify_regular_express", s.VerifyRegularExpress)
	ws.POST("/regular/verify_regular_content_batch", s.VerifyRegularContentBatch)

	ws.Any("/proxy/:method/:target/*target_url", s.ProxyRequest)

	// get changelog info
	ws.POST("findmany/changelog", s.GetVersionList)
	ws.POST("find/changelog/detail", s.GetVersionDetail)

	// common api
	ws.GET("/healthz", s.Healthz)
	ws.GET("/version", ginservice.Version)

	// table instance, only for ui, should be removed later
	s.initModelQuote(ws)

	// field template, only for ui
	s.initFieldTemplate(ws)

	c := &capability.Capability{
		Ws:     ws,
		Engine: s.Engine,
	}
	// init excel func
	excel.Init(c)

	// if no route, redirect to 404 page
	ws.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/#/404")
	})

	return ws
}

func setGinMode() {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(mode)
}

// Healthz TODO
func (s *Service) Healthz(c *gin.Context) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}

	meta.Items = append(meta.Items, zkItem)

	apiServer := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_APISERVER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_APISERVER); err != nil {
		apiServer.IsHealthy = false
		apiServer.Message = err.Error()
	}
	meta.Items = append(meta.Items, apiServer)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "web server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_WEBSERVER,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	answer.SetCommonResponse()
	c.JSON(200, answer)
}
