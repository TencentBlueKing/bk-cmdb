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

func (s *Service) initBusinessObject() {

	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/object", HandlerFunc: s.CreateObjectBatch})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/findmany/object", HandlerFunc: s.SearchObjectBatch})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/object", HandlerFunc: s.CreateObject})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/object", HandlerFunc: s.SearchObject})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/object/{id}", HandlerFunc: s.UpdateObject})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/object/{id}", HandlerFunc: s.DeleteObject})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/find/objecttopology", HandlerFunc: s.SearchObjectTopo})

}

func (s *Service) initBusinessClassification() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/objectclassification", HandlerFunc: s.CreateClassification})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/classificationobject", HandlerFunc: s.SearchClassificationWithObjects})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objectclassification", HandlerFunc: s.SearchClassification})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectclassification/{id}", HandlerFunc: s.UpdateClassification})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/objectclassification/{id}", HandlerFunc: s.DeleteClassification})

}

func (s *Service) initBusinessObjectAttribute() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/objectattr", HandlerFunc: s.CreateObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objectattr", HandlerFunc: s.SearchObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectattr/{id}", HandlerFunc: s.UpdateObjectAttribute})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/objectattr/{id}", HandlerFunc: s.DeleteObjectAttribute})
}

func (s *Service) initBusinessObjectUnique() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/objectunique/object/{bk_obj_id}", HandlerFunc: s.CreateObjectUnique})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectunique/object/{bk_obj_id}/unique/{id}", HandlerFunc: s.UpdateObjectUnique})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/objectunique/object/{bk_obj_id}/unique/{id}", HandlerFunc: s.DeleteObjectUnique})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objectunique/object/{bk_obj_id}", HandlerFunc: s.SearchObjectUnique})
}

func (s *Service) initBusinessObjectAttrGroup() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/objectattgroup", HandlerFunc: s.CreateObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectattgroup", HandlerFunc: s.UpdateObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/objectattgroup/{id}", HandlerFunc: s.DeleteObjectGroup})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectattgroupproperty", HandlerFunc: s.UpdateObjectAttributeGroupProperty, HandlerParseOriginDataFunc: s.ParseUpdateObjectAttributeGroupPropertyInput})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objectattgroup/object/{bk_obj_id}", HandlerFunc: s.SearchGroupByObject})
}

func (s *Service) initBusinessGraphics() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objecttopo/scope_type/{scope_type}/scope_id/{scope_id}", HandlerFunc: s.SelectObjectTopoGraphics})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/update/objecttopo/scope_type/{scope_type}/scope_id/{scope_id}", HandlerFunc: s.UpdateObjectTopoGraphicsNew})
}

func (s *Service) initBusinessAssociation() {

	// mainline topo methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/topomodelmainline", HandlerFunc: s.CreateMainLineObject})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/topomodelmainline/object/{bk_obj_id}", HandlerFunc: s.DeleteMainLineObject})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/topomodelmainline", HandlerFunc: s.SearchMainLineObjectTopo})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/topoinst/biz/{bk_biz_id}", HandlerFunc: s.SearchBusinessTopo})

	// association type methods ,NOT SUPPORT BUSINESS
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/topoassociationtype", HandlerFunc: s.SearchObjectAssoWithAssoKindList})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/associationtype", HandlerFunc: s.SearchAssociationType})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/associationtype", HandlerFunc: s.CreateAssociationType})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/associationtype/{id}", HandlerFunc: s.UpdateAssociationType})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/associationtype/{id}", HandlerFunc: s.DeleteAssociationType})

	// object association methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/objectassociation", HandlerFunc: s.SearchObjectAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/objectassociation", HandlerFunc: s.CreateObjectAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/objectassociation/{id}", HandlerFunc: s.UpdateObjectAssociation})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/objectassociation/{id}", HandlerFunc: s.DeleteObjectAssociation})

	// inst association methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/instassociation", HandlerFunc: s.SearchAssociationInst})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/instassociation", HandlerFunc: s.CreateAssociationInst})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/instassociation/{association_id}", HandlerFunc: s.DeleteAssociationInst})

	// topo search methods
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/instassociation/object/{bk_obj_id}", HandlerFunc: s.SearchInstByAssociation})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/instassttopo/object/{bk_obj_id}/inst/{inst_id}", HandlerFunc: s.SearchInstTopo})

	// ATTENTION: the following methods is not recommended
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/insttopo/object/{bk_obj_id}/inst/{inst_id}", HandlerFunc: s.SearchInstChildTopo})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/import/instassociation/{bk_obj_id}", HandlerFunc: s.ImportInstanceAssociation})

}

func (s *Service) initBusinessInst() {
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/instance/object/{bk_obj_id}", HandlerFunc: s.CreateInst})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/instance/object/{bk_obj_id}/inst/{inst_id}", HandlerFunc: s.DeleteInst})
	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/deletemany/instance/object/{bk_obj_id}", HandlerFunc: s.DeleteInsts})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/instance/object/{bk_obj_id}/inst/{inst_id}", HandlerFunc: s.UpdateInst})
	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/updatemany/instance/object/{bk_obj_id}", HandlerFunc: s.UpdateInsts})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/instance/object/{bk_obj_id}", HandlerFunc: s.SearchInstAndAssociationDetail})
	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/instdetail/object/{bk_obj_id}/inst/{inst_id}", HandlerFunc: s.SearchInstByInstID})

}
