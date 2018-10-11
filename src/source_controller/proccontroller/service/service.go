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
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

type ProctrlServer struct {
	Core     *backbone.Engine
	Instance dal.RDB
	Cache    *redis.Client
}

func (ps *ProctrlServer) WebService() *restful.WebService {

	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return ps.Core.CCErr
	}
	// v3

	ws.Path("/process/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.DELETE("/module").To(ps.DeleteProc2Module))
	ws.Route(ws.POST("/module").To(ps.CreateProc2Module))
	ws.Route(ws.POST("/module/search").To(ps.GetProc2Module))

	ws.Route(ws.POST("/instance/model").To(ps.CreateProcInstanceModel))
	ws.Route(ws.POST("/instance/model/search").To(ps.GetProcInstanceModel))
	ws.Route(ws.DELETE("/instance/model").To(ps.DeleteProcInstanceModel))
	ws.Route(ws.GET("/healthz").To(ps.Healthz))

	return ws
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
	if err := ps.Instance.Ping(); err != nil {
		mongoItem.IsHealthy = false
		mongoItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, mongoItem)

	// redis status
	redisItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityRedis}
	if err := ps.Cache.Ping().Err(); err != nil {
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
		AtTime:     metadata.Now(),
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
