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
)

func (s *coreService) initModelClassification() {
	s.addAction(http.MethodPost, "/create/model/classification", s.CreateOneModelClassification, nil)
	s.addAction(http.MethodPost, "/createmany/model/classification", s.CreateManyModelClassification, nil)
	s.addAction(http.MethodPost, "/setmany/model/classification", s.SetManyModelClassificaiton, nil)
	s.addAction(http.MethodPost, "/set/model/classification", s.SetOneModelClassificaition, nil)
	s.addAction(http.MethodPut, "/update/model/classification", s.UpdateModelClassification, nil)
	s.addAction(http.MethodDelete, "/delete/model/classification", s.DeleteModelClassification, nil)
	s.addAction(http.MethodDelete, "/delete/model/classification/cascade", s.CascadeDeleteModelClassification, nil)
	s.addAction(http.MethodPost, "/read/model/classification", s.SearchModelClassification, nil)
}

func (s *coreService) initModel() {
	s.addAction(http.MethodPost, "/create/model", s.CreateModel, nil)
	s.addAction(http.MethodPost, "/set/model", s.SetModel, nil)
	s.addAction(http.MethodPut, "/update/model", s.UpdateModel, nil)
	s.addAction(http.MethodDelete, "/delete/model", s.DeleteModel, nil)
	s.addAction(http.MethodDelete, "/delete/model/cascade", s.CascadeDeleteModel, nil)
	s.addAction(http.MethodPost, "/read/model", s.SearchModel, nil)

	// init model attribute groups methods
	s.addAction(http.MethodPost, "/create/model/{bk_obj_id}/group", s.CreateModelAttributeGroup, nil)
	s.addAction(http.MethodPost, "/set/model/{bk_obj_id}/group", s.SetModelAttributeGroup, nil)
	s.addAction(http.MethodPut, "/update/model/{bk_obj_id}/group", s.UpdateModelAttributeGroup, nil)
	s.addAction(http.MethodPut, "/update/model/group", s.UpdateModelAttributeGroupByCondition, nil)
	s.addAction(http.MethodDelete, "/delete/model/{bk_obj_id}/group", s.DeleteModelAttributeGroup, nil)
	s.addAction(http.MethodDelete, "/delete/model/group", s.DeleteModelAttributeGroupByCondition, nil)
	s.addAction(http.MethodPost, "/read/model/{bk_obj_id}/group", s.SearchModelAttributeGroup, nil)
	s.addAction(http.MethodPost, "/read/model/group", s.SearchModelAttributeGroupByCondition, nil)

	// init attributes methods
	s.addAction(http.MethodPost, "/create/model/{bk_obj_id}/attributes", s.CreateModelAttributes, nil)
	s.addAction(http.MethodPost, "/set/model/{bk_obj_id}/attributes", s.SetModelAttributes, nil)
	s.addAction(http.MethodPut, "/update/model/{bk_obj_id}/attributes", s.UpdateModelAttributes, nil)
	s.addAction(http.MethodPut, "/update/model/attributes", s.UpdateModelAttributesByCondition, nil)
	s.addAction(http.MethodDelete, "/delete/model/{bk_obj_id}/attributes", s.DeleteModelAttribute, nil)
	s.addAction(http.MethodPost, "/read/model/{bk_obj_id}/attributes", s.SearchModelAttributes, nil)
	s.addAction(http.MethodPost, "/read/model/attributes", s.SearchModelAttributesByCondition, nil)

}

func (s *coreService) initAttrUnique() {
	s.addAction(http.MethodPost, "/read/model/attributes/unique", s.SearchModelAttrUnique, nil)
	s.addAction(http.MethodPost, "/create/model/{bk_obj_id}/attributes/unique", s.CreateModelAttrUnique, nil)
	s.addAction(http.MethodPut, "/update/model/{bk_obj_id}/attributes/unique/{id}", s.UpdateModelAttrUnique, nil)
	s.addAction(http.MethodDelete, "/delete/model/{bk_obj_id}/attributes/unique/{id}", s.DeleteModelAttrUnique, nil)
}

func (s *coreService) initModelInstances() {
	s.addAction(http.MethodPost, "/create/model/{bk_obj_id}/instance", s.CreateOneModelInstance, nil)
	s.addAction(http.MethodPost, "/createmany/model/{bk_obj_id}/instance", s.CreateManyModelInstances, nil)
	s.addAction(http.MethodPut, "/update/model/{bk_obj_id}/instance", s.UpdateModelInstances, nil)
	s.addAction(http.MethodPost, "/read/model/{bk_obj_id}/instances", s.SearchModelInstances, nil)
	s.addAction(http.MethodDelete, "/delete/model/{bk_obj_id}/instance", s.DeleteModelInstances, nil)
	s.addAction(http.MethodDelete, "/delete/model/{bk_obj_id}/instance/cascade", s.CascadeDeleteModelInstances, nil)
}

func (s *coreService) initAssociationKind() {

	s.addAction(http.MethodPost, "/create/associationkind", s.CreateOneAssociationKind, nil)
	s.addAction(http.MethodPost, "/createmany/associationkind", s.CreateManyAssociationKind, nil)
	s.addAction(http.MethodPost, "/set/associationkind", s.SetOneAssociationKind, nil)
	s.addAction(http.MethodPost, "/setmany/associationkind", s.SetManyAssociationKind, nil)
	s.addAction(http.MethodPut, "/update/associationkind", s.UpdateAssociationKind, nil)
	s.addAction(http.MethodDelete, "/delete/associationkind", s.DeleteAssociationKind, nil)
	s.addAction(http.MethodDelete, "/delete/associationkind/cascade", s.CascadeDeleteAssociationKind, nil)
	s.addAction(http.MethodPost, "/read/associationkind", s.SearchAssociationKind, nil)

}

func (s *coreService) initModelAssociation() {

	s.addAction(http.MethodPost, "/create/modelassociation", s.CreateModelAssociation, nil)
	s.addAction(http.MethodPost, "/create/mainlinemodelassociation", s.CreateMainlineModelAssociation, nil)
	s.addAction(http.MethodPost, "/set/modelassociation", s.SetModelAssociation, nil)
	s.addAction(http.MethodPut, "/update/modelassociation", s.UpdateModelAssociation, nil)
	s.addAction(http.MethodPost, "/read/modelassociation", s.SearchModelAssociation, nil)
	s.addAction(http.MethodDelete, "/delete/modelassociation", s.DeleteModelAssociation, nil)
	s.addAction(http.MethodDelete, "/delete/modelassociation/cascade", s.DeleteModelAssociation, nil)
}

func (s *coreService) initInstanceAssociation() {

	s.addAction(http.MethodPost, "/create/instanceassociation", s.CreateOneInstanceAssociation, nil)
	s.addAction(http.MethodPost, "/createmany/instanceassociation", s.CreateManyInstanceAssociation, nil)
	s.addAction(http.MethodPost, "/read/instanceassociation", s.SearchInstanceAssociation, nil)
	s.addAction(http.MethodDelete, "/delete/instanceassociation", s.DeleteInstanceAssociation, nil)
}

func (s *coreService) initMainline() {
	// add handler for model topo and business topo
	s.addAction(http.MethodPost, "/read/mainline/model", s.SearchMainlineModelTopo, nil)
	s.addAction(http.MethodPost, "/read/mainline/instance/{bk_biz_id}", s.SearchMainlineInstanceTopo, nil)
}

func (s *coreService) host() {
	s.addAction(http.MethodPost, "/set/module/host/relation/inner/module", s.TransferHostToDefaultModule, nil)
	s.addAction(http.MethodPost, "/set/module/host/relation/module", s.TransferHostModule, nil)
	s.addAction(http.MethodPost, "/set/module/host/relation/cross/business", s.TransferHostCrossBusiness, nil)
	s.addAction(http.MethodPost, "/read/module/host/relation", s.GetHostModuleRelation, nil)
	s.addAction(http.MethodDelete, "/delete/host", s.DeleteHost, nil)

}

func (s *coreService) audit() {
	s.addAction(http.MethodPost, "/create/auditlog", s.CreateAuditLog, nil)
	s.addAction(http.MethodPost, "/read/auditlog", s.SearchAuditLog, nil)
}

func (s *coreService) initService() {
	s.initModelClassification()
	s.initModel()
	s.initAssociationKind()
	s.initAttrUnique()
	s.initModelAssociation()
	s.initModelInstances()
	s.initInstanceAssociation()
	s.initDataSynchronize()
	s.initMainline()
	s.host()
	s.audit()
}
