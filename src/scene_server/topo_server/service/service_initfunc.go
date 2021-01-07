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

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful"
)

func (s *Service) initAssociation(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	// mainline topo methods
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/topo/model/mainline", Handler: s.CreateMainLineObject})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/topo/model/mainline/owners/{owner_id}/objectids/{bk_obj_id}", Handler: s.DeleteMainLineObject})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/model/{owner_id}", Handler: s.SearchMainLineObjectTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/model/{owner_id}/{cls_id}/{bk_obj_id}", Handler: s.SearchObjectByClassificationID})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/inst/{owner_id}/{bk_biz_id}", Handler: s.SearchBusinessTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/topo/tree/brief/biz/{bk_biz_id}", Handler: s.SearchBriefBizTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/topo/cache/topotree", Handler: s.SearchTopologyTree})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/topo/cache/topo/node_path/biz/{bk_biz_id}",
		Handler: s.SearchTopologyNodePath})

	// TODO: delete this api, it's not used by front.
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", Handler: s.SearchMainLineChildInstTopo})

	// association type methods
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/topo/association/type/action/search/batch", Handler: s.SearchObjectAssocWithAssocKindList})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/topo/association/type/action/search", Handler: s.SearchAssociationType})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/topo/association/type/action/create", Handler: s.CreateAssociationType})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/topo/association/type/{id}/action/update", Handler: s.UpdateAssociationType})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/topo/association/type/{id}/action/delete", Handler: s.DeleteAssociationType})

	// object association methods
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/association/action/search", Handler: s.SearchObjectAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/association/action/create", Handler: s.CreateObjectAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/object/association/{id}/action/update", Handler: s.UpdateObjectAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/object/association/{id}/action/delete", Handler: s.DeleteObjectAssociation})

	// inst association methods
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/action/search", Handler: s.SearchAssociationInst})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/related/action/search", Handler: s.SearchAssociationRelatedInst})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/action/create", Handler: s.CreateAssociationInst})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/inst/association/{association_id}/action/delete", Handler: s.DeleteAssociationInst})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/inst/association/batch/action/delete", Handler: s.DeleteAssociationInstBatch})

	// topo search methods
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/search/owner/{owner_id}/object/{bk_obj_id}", Handler: s.SearchInstByAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{bk_obj_id}/inst/{inst_id}", Handler: s.SearchInstTopo})

	// ATTENTION: the following methods is not recommended
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/search/topo/owner/{owner_id}/object/{bk_obj_id}/inst/{inst_id}", Handler: s.SearchInstChildTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/association/action/{bk_obj_id}/import", Handler: s.ImportInstanceAssociation})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initAuditLog(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/audit_dict", Handler: s.SearchAuditDict})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/audit_list", Handler: s.SearchAuditList})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/audit", Handler: s.SearchAuditDetail})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initBusiness(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/app/search/{owner_id}", Handler: s.SearchBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/app/{owner_id}", Handler: s.CreateBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/app/{owner_id}/{app_id}", Handler: s.DeleteBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/app/{owner_id}/{app_id}", Handler: s.UpdateBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/app/status/{flag}/{owner_id}/{app_id}", Handler: s.UpdateBusinessStatus})
	// utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/app/search/{owner_id}", Handler: s.SearchBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/app/{app_id}/basic_info", Handler: s.GetBusinessBasicInfo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/app/default/{owner_id}/search", Handler: s.SearchOwnerResourcePoolBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/app/default/{owner_id}", Handler: s.CreateDefaultBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/internal/{owner_id}/{app_id}", Handler: s.GetInternalModule})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/topo/internal/{owner_id}/{app_id}/with_statistics", Handler: s.GetInternalModuleWithStatistics})
	// find reduced business list with only few fields for business itself.
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/app/with_reduced", Handler: s.SearchReducedBusinessList})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/app/simplify", Handler: s.ListAllBusinessSimplify})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initModule(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/module/{app_id}/{set_id}", Handler: s.CreateModule})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/module/{app_id}/{set_id}/{module_id}", Handler: s.DeleteModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/module/{app_id}/{set_id}/{module_id}", Handler: s.UpdateModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/module/search/{owner_id}/{bk_biz_id}/{bk_set_id}", Handler: s.SearchModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/module/biz/{bk_biz_id}", Handler: s.SearchModuleByCondition})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/module/bk_biz_id/{bk_biz_id}", Handler: s.SearchModuleBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/module/with_relation/biz/{bk_biz_id}", Handler: s.SearchModuleWithRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/module/bk_biz_id/{bk_biz_id}/service_template_id/{service_template_id}", Handler: s.ListModulesByServiceTemplateID})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/module/host_apply_enable_status/bk_biz_id/{bk_biz_id}/bk_module_id/{bk_module_id}", Handler: s.UpdateModuleHostApplyEnableStatus})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initSet(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/{app_id}", Handler: s.CreateSet})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/{app_id}/batch", Handler: s.BatchCreateSet})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/set/{app_id}/{set_id}", Handler: s.DeleteSet})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/set/{app_id}/{set_id}", Handler: s.UpdateSet})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/search/{owner_id}/{app_id}", Handler: s.SearchSet})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/set/bk_biz_id/{bk_biz_id}", Handler: s.SearchSetBatch})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initInst(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/{owner_id}/{bk_obj_id}", Handler: s.CreateInst})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/inst/{owner_id}/{bk_obj_id}/{inst_id}", Handler: s.DeleteInst})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/inst/{owner_id}/{bk_obj_id}/batch", Handler: s.DeleteInsts})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/inst/{owner_id}/{bk_obj_id}/{inst_id}", Handler: s.UpdateInst})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/inst/{owner_id}/{bk_obj_id}/batch/update", Handler: s.UpdateInsts})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/search/{owner_id}/{bk_obj_id}", Handler: s.SearchInsts})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{bk_obj_id}/detail", Handler: s.SearchInstAndAssociationDetail})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{bk_obj_id}", Handler: s.SearchInstByObject})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inst/search/{owner_id}/{bk_obj_id}/{inst_id}", Handler: s.SearchInstByInstID})
	// 2019-09-30 废弃接口
	// utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/inst/association/object/{bk_obj_id}/inst_id/{id}/offset/{start}/limit/{limit}", Handler: s.SearchInstAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/inst/association/object/{bk_obj_id}/inst_id/{id}/offset/{start}/limit/{limit}/web", Handler: s.SearchInstAssociationUI})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/inst/association/association_object/inst_base_info", Handler: s.SearchInstAssociationWithOtherObject})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initObjectAttribute(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objectattr", Handler: s.CreateObjectAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objectattr/search", Handler: s.SearchObjectAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/objectattr/{id}", Handler: s.UpdateObjectAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/objectattr/{id}", Handler: s.DeleteObjectAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/objectattr/index/{bk_obj_id}/{id}", Handler: s.UpdateObjectAttributeIndex})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initObjectClassification(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/classification", Handler: s.CreateClassification})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/classification/{owner_id}/objects", Handler: s.SearchClassificationWithObjects})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/classifications", Handler: s.SearchClassification})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/object/classification/{id}", Handler: s.UpdateClassification})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/object/classification/{id}", Handler: s.DeleteClassification})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initObjectObjectUnique(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/{bk_obj_id}/unique/action/create", Handler: s.CreateObjectUnique})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/object/{bk_obj_id}/unique/{id}/action/update", Handler: s.UpdateObjectUnique})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/object/{bk_obj_id}/unique/{id}/action/delete", Handler: s.DeleteObjectUnique})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/object/{bk_obj_id}/unique/action/search", Handler: s.SearchObjectUnique})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initObjectGroup(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objectatt/group/new", Handler: s.CreateObjectGroup})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/objectatt/group/update", Handler: s.UpdateObjectGroup})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/objectatt/group/groupid/{id}", Handler: s.DeleteObjectGroup})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/objectatt/group/property", Handler: s.UpdateObjectAttributeGroupProperty})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/objectatt/group/owner/{owner_id}/object/{bk_object_id}/propertyids/{property_id}/groupids/{group_id}", Handler: s.DeleteObjectAttributeGroup})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{bk_obj_id}", Handler: s.SearchGroupByObject})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initObject(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/batch", Handler: s.CreateObjectBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object/search/batch", Handler: s.SearchObjectBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/object", Handler: s.CreateObject})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objects", Handler: s.SearchObject})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objects/topo", Handler: s.SearchObjectTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/object/{id}", Handler: s.UpdateObject})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/object/{id}", Handler: s.DeleteObject})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/object/statistics", Handler: s.GetModelStatistics})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initGraphics(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/search", Handler: s.SelectObjectTopoGraphics})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/update", Handler: s.UpdateObjectTopoGraphics})

	utility.AddToRestfulWebService(web)
}
func (s *Service) initIdentifier(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/identifier/{obj_type}/search", Handler: s.SearchIdentifier})

	utility.AddToRestfulWebService(web)
}

// 全文索引
func (s *Service) initFullTextSearch(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/full_text", Handler: s.FullTextFind})

	utility.AddToRestfulWebService(web)
}

// 资源池目录
func (s *Service) initResourceDirectory(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/resource/directory", Handler: s.CreateResourceDirectory})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/resource/directory/{bk_module_id}", Handler: s.UpdateResourceDirectory})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/resource/directory", Handler: s.SearchResourceDirectory})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/resource/directory/{bk_module_id}", Handler: s.DeleteResourceDirectory})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initService(web *restful.WebService) {
	s.initAssociation(web)
	s.initAuditLog(web)
	s.initBusiness(web)
	s.initInst(web)
	s.initModule(web)
	s.initSet(web)
	s.initObject(web)
	s.initObjectAttribute(web)
	s.initObjectClassification(web)
	s.initObjectGroup(web)
	s.initGraphics(web)
	s.initIdentifier(web)
	s.initObjectObjectUnique(web)

	s.initBusinessObject(web)
	s.initBusinessClassification(web)
	s.initBusinessObjectAttribute(web)
	s.initBusinessObjectUnique(web)
	s.initBusinessObjectAttrGroup(web)
	s.initBusinessAssociation(web)
	s.initBusinessGraphics(web)
	s.initBusinessInst(web)

	s.initFullTextSearch(web)
	s.initSetTemplate(web)
	s.initInternalTask(web)

	s.initResourceDirectory(web)
}
