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
	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal"
)

type Service struct {
	Core     *backbone.Engine
	Instance dal.RDB
	Cache    *redis.Client
}

func (s *Service) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.Core.CCErr
	}
	api.Path("/object/{version}").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	api.Route(api.POST("/identifier/{obj_type}/search").To(s.SearchIdentifier))

	api.Route(api.POST("/insts/{obj_type}/search").To(s.SearchInstObjects))
	api.Route(api.POST("/insts/{obj_type}").To(s.CreateInstObject))
	api.Route(api.DELETE("/insts/{obj_type}").To(s.DeleteInstObject))
	api.Route(api.PUT("/insts/{obj_type}").To(s.UpdateInstObject))

	api.Route(api.POST("/meta/objects").To(s.SelectObjects))
	api.Route(api.DELETE("/meta/object/{id}").To(s.DeleteObject))
	api.Route(api.POST("/meta/object").To(s.CreateObject))
	api.Route(api.PUT("/meta/object/{id}").To(s.UpdateObject))

	api.Route(api.POST("/meta/objectassts").To(s.SelectObjectAssociations))
	api.Route(api.DELETE("/meta/objectasst/{id}").To(s.DeleteObjectAssociation))
	api.Route(api.POST("/meta/objectasst").To(s.CreateObjectAssociation))
	api.Route(api.POST("/meta/mainlineobjectasst").To(s.CreateMainlineObjectAssociation))
	api.Route(api.PUT("/meta/objectasst/{id}").To(s.UpdateObjectAssociation))

	api.Route(api.POST("/meta/objectatt/{id}").To(s.SelectObjectAttByID))
	api.Route(api.POST("/meta/objectatts").To(s.SelectObjectAttWithParams))
	api.Route(api.DELETE("/meta/objectatt/{id}").To(s.DeleteObjectAttByID))
	api.Route(api.POST("/meta/objectatt").To(s.CreateObjectAtt))
	api.Route(api.PUT("/meta/objectatt/{id}").To(s.UpdateObjectAttByID))

	api.Route(api.POST("/meta/objectatt/group/new").To(s.CreatePropertyGroup))
	api.Route(api.PUT("/meta/objectatt/group/update").To(s.UpdatePropertyGroup))
	api.Route(api.DELETE("/meta/objectatt/group/groupid/{id}").To(s.DeletePropertyGroup))
	api.Route(api.PUT("/meta/objectatt/group/property").To(s.UpdatePropertyGroupObjectAtt))
	api.Route(api.DELETE("/meta/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}").To(s.DeletePropertyGroupObjectAtt))
	api.Route(api.POST("/meta/objectatt/group/property/owner/{owner_id}/object/{object_id}").To(s.SelectPropertyGroupByObjectID))
	api.Route(api.POST("/meta/objectatt/group/search").To(s.SelectGroup))

	api.Route(api.POST("/meta/object/classification/{owner_id}/objects").To(s.SelectClassificationWithObject))
	api.Route(api.POST("/meta/object/classification/search").To(s.SelectClassifications))
	api.Route(api.DELETE("/meta/object/classification/{id}").To(s.DeleteClassification))
	api.Route(api.POST("/meta/object/classification").To(s.CreateClassification))
	api.Route(api.PUT("/meta/object/classification/{id}").To(s.UpdateClassification))

	api.Route(api.POST("/object/{bk_obj_id}/unique/action/create").To(s.CreateObjectUnique))
	api.Route(api.PUT("/object/{bk_obj_id}/unique/{id}/action/update").To(s.UpdateObjectUnique))
	api.Route(api.DELETE("/object/{bk_obj_id}/unique/{id}/action/delete").To(s.DeleteObjectUnique))
	api.Route(api.GET("/object/{bk_obj_id}/unique/action/search").To(s.SearchObjectUnique))

	// association api
	api.Route(api.POST("/association/action/search").To(s.SearchAssociationType))
	api.Route(api.POST("/association/action/create").To(s.CreateAssociationType))
	api.Route(api.PUT("/association/{id}/action/update").To(s.UpdateAssociationType))
	api.Route(api.DELETE("/association/{id}/action/delete").To(s.DeleteAssociationType))

	api.Route(api.POST("/object/association/action/search").To(s.SelectObjectAssociations))                 // optimization: new api path
	api.Route(api.POST("/object/association/action/create").To(s.CreateObjectAssociation))                  // optimization: new api path
	api.Route(api.POST("/object/association/mainline/action/create").To(s.CreateMainlineObjectAssociation)) // interface mainline association
	api.Route(api.PUT("/object/association/{id}/action/update").To(s.UpdateObjectAssociation))              // optimization: new api path
	api.Route(api.DELETE("/object/association/{id}/action/delete").To(s.DeleteObjectAssociation))           // optimization: new api path

	api.Route(api.POST("/inst/association/action/search").To(s.SearchInstAssociations))
	api.Route(api.POST("/inst/association/action/create").To(s.CreateInstAssociation))
	api.Route(api.DELETE("/inst/association/{association_id}/action/delete").To(s.DeleteInstAssociation))

	api.Route(api.POST("/topographics/search").To(s.SearchTopoGraphics))
	api.Route(api.POST("/topographics/update").To(s.UpdateTopoGraphics))

	api.Route(api.POST("/openapi/proc/getProcModule").To(s.GetProcessesByModuleName))
	api.Route(api.DELETE("/openapi/set/delhost").To(s.DeleteSetHost))

	api.Route(api.POST("/privilege/group/{bk_supplier_account}").To(s.CreateUserGroup))
	api.Route(api.PUT("/privilege/group/{bk_supplier_account}/{group_id}").To(s.UpdateUserGroup))
	api.Route(api.DELETE("/privilege/group/{bk_supplier_account}/{group_id}").To(s.DeleteUserGroup))
	api.Route(api.POST("/privilege/group/{bk_supplier_account}/search").To(s.SearchUserGroup))

	api.Route(api.POST("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.CreateUserGroupPrivi))
	api.Route(api.PUT("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.UpdateUserGroupPrivi))
	api.Route(api.GET("/privilege/group/detail/{bk_supplier_account}/{group_id}").To(s.GetUserGroupPrivi))

	api.Route(api.POST("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.CreateRolePri))
	api.Route(api.GET("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.GetRolePri))
	api.Route(api.PUT("/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}").To(s.UpdateRolePri))

	api.Route(api.GET("/system/{flag}/{bk_supplier_account}").To(s.GetSystemFlag))

	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
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
		AtTime:     metadata.Now(),
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
