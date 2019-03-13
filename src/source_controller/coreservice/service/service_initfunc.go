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

func (s *coreService) initHealth() {
	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/healthz", HandlerFunc: s.Health})
}

func (s *coreService) initModelClassification() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model/classification", HandlerFunc: s.CreateOneModelClassification})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/model/classification", HandlerFunc: s.CreateManyModelClassification})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/setmany/model/classification", HandlerFunc: s.SetManyModelClassificaiton})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/model/classification", HandlerFunc: s.SetOneModelClassificaition})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/classification", HandlerFunc: s.UpdateModelClassification})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/classification", HandlerFunc: s.DeleteModelClassification})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/classification/cascade", HandlerFunc: s.CascadeDeleteModelClassification})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/classification", HandlerFunc: s.SearchModelClassification})
}

func (s *coreService) initModel() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model", HandlerFunc: s.CreateModel})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/model", HandlerFunc: s.SetModel})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model", HandlerFunc: s.UpdateModel})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model", HandlerFunc: s.DeleteModel})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/cascade", HandlerFunc: s.CascadeDeleteModel})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model", HandlerFunc: s.SearchModel})

	// init model attribute groups methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model/{bk_obj_id}/group", HandlerFunc: s.CreateModelAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/model/{bk_obj_id}/group", HandlerFunc: s.SetModelAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/{bk_obj_id}/group", HandlerFunc: s.UpdateModelAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/group", HandlerFunc: s.UpdateModelAttributeGroupByCondition})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/{bk_obj_id}/group", HandlerFunc: s.DeleteModelAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/group", HandlerFunc: s.DeleteModelAttributeGroupByCondition})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/{bk_obj_id}/group", HandlerFunc: s.SearchModelAttributeGroup})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/group", HandlerFunc: s.SearchModelAttributeGroupByCondition})

	// init attributes methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model/{bk_obj_id}/attributes", HandlerFunc: s.CreateModelAttributes})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/model/{bk_obj_id}/attributes", HandlerFunc: s.SetModelAttributes})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/{bk_obj_id}/attributes", HandlerFunc: s.UpdateModelAttributes})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/attributes", HandlerFunc: s.UpdateModelAttributesByCondition})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/{bk_obj_id}/attributes", HandlerFunc: s.DeleteModelAttribute})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/{bk_obj_id}/attributes", HandlerFunc: s.SearchModelAttributes})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/attributes", HandlerFunc: s.SearchModelAttributesByCondition})

}

func (s *coreService) initAttrUnique() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/attributes/unique", HandlerFunc: s.SearchModelAttrUnique})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model/{bk_obj_id}/attributes/unique", HandlerFunc: s.CreateModelAttrUnique})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/{bk_obj_id}/attributes/unique/{id}", HandlerFunc: s.UpdateModelAttrUnique})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/{bk_obj_id}/attributes/unique/{id}", HandlerFunc: s.DeleteModelAttrUnique})
}

func (s *coreService) initModelInstances() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/model/{bk_obj_id}/instance", HandlerFunc: s.CreateOneModelInstance})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/model/{bk_obj_id}/instance", HandlerFunc: s.CreateManyModelInstances})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/model/{bk_obj_id}/instance", HandlerFunc: s.UpdateModelInstances})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/model/{bk_obj_id}/instances", HandlerFunc: s.SearchModelInstances})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/{bk_obj_id}/instance", HandlerFunc: s.DeleteModelInstances})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/model/{bk_obj_id}/instance/cascade", HandlerFunc: s.CascadeDeleteModelInstances})
}

func (s *coreService) initAssociationKind() {

	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/associationkind", HandlerFunc: s.CreateOneAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/associationkind", HandlerFunc: s.CreateManyAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/associationkind", HandlerFunc: s.SetOneAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/setmany/associationkind", HandlerFunc: s.SetManyAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/associationkind", HandlerFunc: s.UpdateAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/associationkind", HandlerFunc: s.DeleteAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/associationkind/cascade", HandlerFunc: s.CascadeDeleteAssociationKind})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/associationkind", HandlerFunc: s.SearchAssociationKind})

}

func (s *coreService) initModelAssociation() {

	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/modelassociation", HandlerFunc: s.CreateModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/mainlinemodelassociation", HandlerFunc: s.CreateMainlineModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/set/modelassociation", HandlerFunc: s.SetModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/modelassociation", HandlerFunc: s.UpdateModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/modelassociation", HandlerFunc: s.SearchModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/modelassociation", HandlerFunc: s.DeleteModelAssociation})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/modelassociation/cascade", HandlerFunc: s.DeleteModelAssociation})
}

func (s *coreService) initInstanceAssociation() {

	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/instanceassociation", HandlerFunc: s.CreateOneInstanceAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/instanceassociation", HandlerFunc: s.CreateManyInstanceAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/read/instanceassociation", HandlerFunc: s.SearchInstanceAssociation})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/instanceassociation", HandlerFunc: s.DeleteInstanceAssociation})
}

func (s *coreService) initService() {
	s.initHealth()
	s.initModelClassification()
	s.initModel()
	s.initAssociationKind()
	s.initAttrUnique()
	s.initModelAssociation()
	s.initModelInstances()
	s.initInstanceAssociation()
	s.initDataSynchronize()
}
