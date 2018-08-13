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
	"net/http"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cfnc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/storage"
	"configcenter/src/storage/mgoclient"
	"configcenter/src/storage/redisclient"
)

type ProctrlServer struct {
	Core       *backbone.Engine
	DbInstance storage.DI
	CacheDI    *redis.Client
	MongoCfg   *mgoclient.MongoConfig
	RedisCfg   *redisclient.RedisConfig
}

func (ps *ProctrlServer) WebService() http.Handler {

	container := restful.NewContainer()
	getErrFun := func() errors.CCErrorIf {
		return ps.Core.CCErr
	}
	// v3
	ws := new(restful.WebService)
	ws.Path("/process/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.DELETE("/module").To(ps.DeleteProc2Module))
	ws.Route(ws.POST("/module").To(ps.CreateProc2Module))
	ws.Route(ws.POST("/module/search").To(ps.GetProc2Module))

	ws.Route(ws.POST("/conftemp").To(ps.CreateConfigTemp))
	ws.Route(ws.PUT("/conftemp").To(ps.UpdateConfigTemp))
	ws.Route(ws.DELETE("/conftemp").To(ps.DeleteConfigTemp))
	ws.Route(ws.POST("/conftemp/search").To(ps.QueryConfigTemp))

	ws.Route(ws.POST("/instance/model").To(ps.CreateProcInstanceModel))
	ws.Route(ws.POST("/instance/model/search").To(ps.GetProcInstanceModel))
	ws.Route(ws.DELETE("/instance/model").To(ps.DeleteProcInstanceModel))
	ws.Route(ws.GET("/healthz").To(ps.Healthz))

	container.Add(ws)

	return container
}

func (ps *ProctrlServer) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := ps.Core.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongodb status
	mongoItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityMongo}
	if err := ps.DbInstance.Ping(); err != nil {
		mongoItem.IsHealthy = false
		mongoItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, mongoItem)

	// redis status
	redisItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityRedis}
	if err := ps.CacheDI.Ping().Err(); err != nil {
		redisItem.IsHealthy = false
		redisItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, redisItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "audit controller is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_PROCCONTROLLER,
		HealthMeta: meta,
		AtTime:     types.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	resp.WriteJson(answer, "application/json")
}

func (ps *ProctrlServer) OnProcessConfUpdate(previous, current cfnc.ProcessConfig) {
	prefix := storage.DI_MONGO
	ps.MongoCfg = &mgoclient.MongoConfig{
		Address:      current.ConfigMap[prefix+".host"],
		User:         current.ConfigMap[prefix+".usr"],
		Password:     current.ConfigMap[prefix+".pwd"],
		Database:     current.ConfigMap[prefix+".database"],
		Port:         current.ConfigMap[prefix+".port"],
		MaxOpenConns: current.ConfigMap[prefix+".maxOpenConns"],
		MaxIdleConns: current.ConfigMap[prefix+".maxIDleConns"],
		Mechanism:    current.ConfigMap[prefix+".mechanism"],
	}

	prefix = storage.DI_REDIS
	ps.RedisCfg = &redisclient.RedisConfig{
		Address:  current.ConfigMap[prefix+".host"],
		User:     current.ConfigMap[prefix+".usr"],
		Password: current.ConfigMap[prefix+".pwd"],
		Database: current.ConfigMap[prefix+".database"],
		Port:     current.ConfigMap[prefix+".port"],
	}
}
