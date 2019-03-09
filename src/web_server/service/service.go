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
	"plugin"
	"strings"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/logics"
	"configcenter/src/web_server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
	redis "gopkg.in/redis.v5"
)

type Service struct {
	VersionPlg *plugin.Plugin
	*options.ServerOption
	Engine   *backbone.Engine
	CacheCli *redis.Client
	*logics.Logics
	Disc   discovery.DiscoveryInterface
	Config options.Config
}

func (s *Service) WebService() *gin.Engine {
	ws := gin.Default()

	var store sessions.RedisStore
	var redisErr error
	if 0 == len(s.Config.Session.Address) {
		address := s.Config.Session.Host + ":" + s.Config.Session.Port
		store, redisErr = sessions.NewRedisStore(10, "tcp", address, s.Config.Session.Secret, []byte("secret"))
		if redisErr != nil {
			blog.Fatal("failed to create new redis store, error info is %v", redisErr)
		}
	} else {
		address := strings.Split(s.Config.Session.Address, ";")
		store, redisErr = sessions.NewRedisStoreWithSentinel(address, 10, s.Config.Session.MasterName, "tcp", s.Config.Session.Secret, []byte("secret"))
		if redisErr != nil {
			blog.Fatal("failed to create new redis store, error info is %v", redisErr)
		}
	}

	ws.Use(sessions.Sessions(s.Config.Session.Name, store))
	middleware.Engine = s.Engine
	ws.Use(middleware.ValidLogin(s.Config, s.Disc))

	ws.Static("/static", s.Config.Site.HtmlRoot)
	ws.LoadHTMLFiles(s.Config.Site.HtmlRoot + "/index.html")

	ws.POST("/hosts/import", s.ImportHost)
	ws.POST("/hosts/export", s.ExportHost)
	ws.GET("/importtemplate/:bk_obj_id", s.BuildDownLoadExcelTemplate)
	ws.POST("/insts/owner/:bk_supplier_account/object/:bk_obj_id/import", s.ImportInst)
	ws.POST("/insts/owner/:bk_supplier_account/object/:bk_obj_id/export", s.ExportInst)
	ws.POST("/logout", s.LogOutUser)
	ws.POST("/object/owner/:bk_supplier_account/object/:bk_obj_id/import", s.ImportObject)
	ws.POST("/object/owner/:bk_supplier_account/object/:bk_obj_id/export", s.ExportObject)
	ws.GET("/user/list", s.GetUserList)
	ws.GET("/user/language/:language", s.UpdateUserLanguage)
	ws.GET("/userinfo", s.UserInfo)
	ws.PUT("/user/current/supplier/:id", s.UpdateSupplier)

	ws.GET("/healthz", s.Healthz)
	ws.GET("/", s.Index)
	return ws
}

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
	c.JSON(200, answer)
}
