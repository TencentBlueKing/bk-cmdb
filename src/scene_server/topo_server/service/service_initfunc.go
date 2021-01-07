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

	"configcenter/src/common"
)

func (s *topoService) initHealth() {
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/healthz", HandlerFunc: s.Health})
}

func (s *topoService) initAssociation() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/model/mainline", HandlerFunc: s.CreateMainLineObject})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/topo/model/mainline/owners/{owner_id}/objectids/{obj_id}", HandlerFunc: s.DeleteMainLineObject})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/model/{owner_id}", HandlerFunc: s.SearchMainLineOBjectTopo})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/model/{owner_id}/{cls_id}/{obj_id}", HandlerFunc: s.SearchObjectByClassificationID})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/inst/{owner_id}/{app_id}", HandlerFunc: s.SearchBusinessTopo})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", HandlerFunc: s.SearchMainLineChildInstTopo})
}

func (s *topoService) initAuditLog() {

	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/audit/search", HandlerFunc: s.AuditQuery})
}

func (s *topoService) initCompatiblev2() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/app/searchAll", HandlerFunc: s.SearchAllApp})

	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/openapi/set/multi/{appid}", HandlerFunc: s.UpdateMultiSet})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/openapi/set/multi/{appid}", HandlerFunc: s.DeleteMultiSet})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/openapi/set/setHost/{appid}", HandlerFunc: s.DeleteSetHost})

	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", HandlerFunc: s.UpdateMultiModule})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/openapi/module/searchByApp/{" + common.BKAppIDField + "}", HandlerFunc: s.SearchModuleByApp})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/openapi/module/searchByProperty/{" + common.BKAppIDField + "}", HandlerFunc: s.SearchModuleBySetProperty})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/openapi/module/multi", HandlerFunc: s.AddMultiModule})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", HandlerFunc: s.DeleteMultiModule})

}

func (s *topoService) initBusiness() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/app/{owner_id}", HandlerFunc: s.CreateBusiness})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/app/{owner_id}/{app_id}", HandlerFunc: s.DeleteBusiness})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/app/{owner_id}/{app_id}", HandlerFunc: s.UpdateBusiness})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/app/status/{flag}/{owner_id}/{app_id}", HandlerFunc: s.UpdateBusinessStatus})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/app/search/{owner_id}", HandlerFunc: s.SearchBusiness})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}/search", HandlerFunc: s.SearchDefaultBusiness})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}", HandlerFunc: s.CreateDefaultBusiness})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/internal/{owner_id}/{app_id}", HandlerFunc: s.GetInternalModule})
}

func (s *topoService) initModule() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/module/{app_id}/{set_id}", HandlerFunc: s.CreateModule})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: s.DeleteModule})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: s.UpdateModule})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/module/search/{owner_id}/{app_id}/{set_id}", HandlerFunc: s.SearchModule})

}

func (s *topoService) initSet() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/{app_id}", HandlerFunc: s.CreateSet})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/set/{app_id}/{set_id}", HandlerFunc: s.DeleteSet})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/set/{app_id}/batch", HandlerFunc: s.DeleteSets})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/set/{app_id}/{set_id}", HandlerFunc: s.UpdateSet})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/search/{owner_id}/{app_id}", HandlerFunc: s.SearchSet})

}

func (s *topoService) initInst() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/{owner_id}/{obj_id}", HandlerFunc: s.CreateInst})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: s.DeleteInst})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/inst/{owner_id}/{obj_id}/batch", HandlerFunc: s.DeleteInsts})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: s.UpdateInst})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/inst/{owner_id}/{obj_id}/batch/update", HandlerFunc: s.UpdateInsts})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}", HandlerFunc: s.SearchInsts})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}/detail", HandlerFunc: s.SearchInstAndAssociationDetail})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: s.SearchInstByObject})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: s.SearchInstByAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: s.SearchInstByInstID})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: s.SearchInstChildTopo})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: s.SearchInstTopo})
}

func (s *topoService) initObjectAttribute() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectattr", HandlerFunc: s.CreateObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectattr/search", HandlerFunc: s.SearchObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectattr/{id}", HandlerFunc: s.UpdateObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectattr/{id}", HandlerFunc: s.DeleteObjectAttribute})
}

func (s *topoService) initObjectClassification() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classification", HandlerFunc: s.CreateClassification})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classification/{owner_id}/objects", HandlerFunc: s.SearchClassificationWithObjects})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classifications", HandlerFunc: s.SearchClassification})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/object/classification/{id}", HandlerFunc: s.UpdateClassification})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/object/classification/{id}", HandlerFunc: s.DeleteClassification})
}

func (s *topoService) initObjectGroup() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectatt/group/new", HandlerFunc: s.CreateObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectatt/group/update", HandlerFunc: s.UpdateObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/groupid/{id}", HandlerFunc: s.DeleteObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectatt/group/property", HandlerFunc: s.UpdateObjectAttributeGroup, HandlerParseOriginDataFunc: s.ParseUpdateObjectAttributeGroupInput})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", HandlerFunc: s.DeleteObjectAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", HandlerFunc: s.SearchGroupByObject})
}

func (s *topoService) initObject() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/batch", HandlerFunc: s.CreateObjectBatch})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/search/batch", HandlerFunc: s.SearchObjectBatch})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object", HandlerFunc: s.CreateObject})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objects", HandlerFunc: s.SearchObject})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objects/topo", HandlerFunc: s.SearchObjectTopo})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/object/{id}", HandlerFunc: s.UpdateObject})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/object/{id}", HandlerFunc: s.DeleteObject})
}
func (s *topoService) initPrivilegeGroup() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/privilege/group/{bk_supplier_account}", HandlerFunc: s.CreateUserGroup})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/topo/privilege/group/{bk_supplier_account}/{group_id}", HandlerFunc: s.DeleteUserGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/topo/privilege/group/{bk_supplier_account}/{group_id}", HandlerFunc: s.UpdateUserGroup})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/privilege/group/{bk_supplier_account}/search", HandlerFunc: s.SearchUserGroup})
}

func (s *topoService) initPrivigeRole() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", HandlerFunc: s.CreatePrivilege, HandlerParseOriginDataFunc: s.ParseCreateRolePrivilegeOriginData})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", HandlerFunc: s.GetPrivilege})
}

func (s *topoService) initPrivilege() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/privilege/group/detail/{bk_supplier_account}/{group_id}", HandlerFunc: s.UpdateUserGroupPrivi})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/privilege/group/detail/{bk_supplier_account}/{group_id}", HandlerFunc: s.GetUserGroupPrivi})
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/privilege/user/detail/{bk_supplier_account}/{user_name}", HandlerFunc: s.GetUserPrivi})
}

func (s *topoService) initGraphics() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/search", HandlerFunc: s.SelectObjectTopoGraphics})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/update", HandlerFunc: s.UpdateObjectTopoGraphics, HandlerParseOriginDataFunc: s.ParseOriginGraphicsUpdateInput})
}
func (s *topoService) initIdentifier() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/identifier/{obj_type}/search", HandlerFunc: s.SearchIdentifier, HandlerParseOriginDataFunc: s.ParseSearchIdentifierOriginData})
}

func (s *topoService) initService() {
	s.initHealth()
	s.initAssociation()
	s.initAuditLog()
	s.initCompatiblev2()
	s.initBusiness()
	s.initInst()
	s.initModule()
	s.initSet()
	s.initObject()
	s.initObjectAttribute()
	s.initObjectClassification()
	s.initObjectGroup()
	s.initPrivilegeGroup()
	s.initPrivigeRole()
	s.initPrivilege()
	s.initGraphics()
	s.initIdentifier()
}
