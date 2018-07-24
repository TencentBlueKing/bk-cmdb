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
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"

	"configcenter/src/storage"
)

type Service struct {
	Core     *backbone.Engine
	Instance storage.DI
	Cache    *redis.Client
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.Core.CCErr
	}

	ws.Path("/object/{version}").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)
	//restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/identifier/{obj_type}/search").To(s.SearchIdentifier))

	ws.Route(ws.POST("/insts/{obj_type}/search").To(s.SearchInstObjects))
	ws.Route(ws.POST("/insts/{obj_type}").To(s.CreateInstObject))
	ws.Route(ws.DELETE("/insts/{obj_type}").To(s.DeleteInstObject))
	ws.Route(ws.PUT("/insts/{obj_type}").To(s.UpdateInstObject))

	ws.Route(ws.POST("/meta/objects").To(s.SelectObjects))
	ws.Route(ws.DELETE("/meta/object/{id}").To(s.DeleteObject))
	ws.Route(ws.POST("/meta/object").To(s.CreateObject))
	ws.Route(ws.PUT("/meta/object/{id}").To(s.UpdateObject))

	ws.Route(ws.POST("/meta/objectassts").To(s.SelectObjectAssociations))
	ws.Route(ws.DELETE("/meta/objectasst/{id}").To(s.DeleteObjectAssociation))
	ws.Route(ws.POST("/meta/objectasst").To(s.CreateObjectAssociation))
	ws.Route(ws.PUT("/meta/objectasst/{id}").To(s.UpdateObjectAssociation))

	ws.Route(ws.POST("/meta/objectatt/{id}").To(s.SelectObjectAttByID))
	ws.Route(ws.POST("/meta/objectatts").To(s.SelectObjectAttWithParams))
	ws.Route(ws.DELETE("/meta/objectatt/{id}").To(s.DeleteObjectAttByID))
	ws.Route(ws.POST("/meta/objectatt").To(s.CreateObjectAtt))
	ws.Route(ws.PUT("/meta/objectatt/{id}").To(s.UpdateObjectAttByID))

	ws.Route(ws.POST("/meta/objectatt/group/new").To(s.CreatePropertyGroup))
	ws.Route(ws.PUT("/meta/objectatt/group/update").To(s.UpdatePropertyGroup))
	ws.Route(ws.DELETE("/meta/objectatt/group/groupid/{id}").To(s.DeletePropertyGroup))
	ws.Route(ws.PUT("/meta/objectatt/group/property").To(s.UpdatePropertyGroupObjectAtt))
	ws.Route(ws.DELETE("/meta/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}").To(s.DeletePropertyGroupObjectAtt))
	ws.Route(ws.POST("/meta/objectatt/group/property/owner/{owner_id}/object/{object_id}").To(s.SelectPropertyGroupByObjectID))
	ws.Route(ws.POST("/meta/objectatt/group/search").To(s.SelectGroup))

	ws.Route(ws.POST("/meta/object/classification/{owner_id}/objects").To(s.SelectClassificationWithObject))
	ws.Route(ws.POST("/meta/object/classification/search").To(s.SelectClassifications))
	ws.Route(ws.DELETE("/meta/object/classification/{id}").To(s.DeleteClassification))
	ws.Route(ws.POST("/meta/object/classification").To(s.CreateClassification))
	ws.Route(ws.PUT("/meta/object/classification/{id}").To(s.UpdateClassification))

	ws.Route(ws.POST("/topographics/search").To(s.SearchTopoGraphics))
	ws.Route(ws.POST("/topographics/update").To(s.UpdateTopoGraphics))

	ws.Route(ws.POST("/openapi/proc/getProcModule").To(s.GetProcessesByModuleName))
	ws.Route(ws.DELETE("/openapi/set/delhost").To(s.DeleteSetHost))

	ws.Route(ws.POST("/privilege/group/{bk_supplier_account}").To(s.CreateUserGroup))
	ws.Route(ws.PUT("/privilege/group/{bk_supplier_account}/{group_id}").To(s.UpdateUserGroup))
	ws.Route(ws.DELETE("/privilege/group/{bk_supplier_account}/{group_id}").To(s.DeleteUserGroup))
	ws.Route(ws.POST("/privilege/group/{bk_supplier_account}/search").To(s.SearchUserGroup))

	ws.Route(ws.POST("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.CreateUserGroupPrivi))
	ws.Route(ws.PUT("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.UpdateUserGroupPrivi))
	ws.Route(ws.GET("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.GetUserGroupPrivi))

	ws.Route(ws.POST("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.CreateRolePri))
	ws.Route(ws.GET("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.GetRolePri))
	ws.Route(ws.PUT("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.UpdateRolePri))

	ws.Route(ws.GET("/system/{flag}/{bk_supplier_account}").To(s.GetSystemFlag))

	ws.Route(ws.GET("/healthz").To(s.Healthz))
	return ws
}

func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Core.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongo health status
	mongoItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityMongo}
	if err := s.Instance.Ping(); nil != err {
		mongoItem.IsHealthy = false
		mongoItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, mongoItem)

	// redis status
	redisItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityRedis}
	if err := s.Cache.Ping().Err(); err != nil {
		redisItem.IsHealthy = false
		redisItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, redisItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "object controller is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_OBJECTCONTROLLER,
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
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}
