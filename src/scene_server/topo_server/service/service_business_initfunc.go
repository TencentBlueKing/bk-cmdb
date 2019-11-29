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

	s.addAction(http.MethodPost, "/createmany/object", s.CreateObjectBatch, nil)
	s.addAction(http.MethodPost, "/findmany/object", s.SearchObjectBatch, nil)
	s.addAction(http.MethodPost, "/create/object", s.CreateObject, nil)
	s.addAction(http.MethodPost, "/find/object", s.SearchObject, nil)
	s.addAction(http.MethodPut, "/update/object/{id}", s.UpdateObject, nil)
	s.addAction(http.MethodDelete, "/delete/object/{id}", s.DeleteObject, nil)
	s.addAction(http.MethodPut, "/find/objecttopology", s.SearchObjectTopo, nil)
}

func (s *Service) initBusinessClassification() {
	s.addAction(http.MethodPost, "/create/objectclassification", s.CreateClassification, nil)
	s.addAction(http.MethodPost, "/find/classificationobject", s.SearchClassificationWithObjects, nil)
	s.addAction(http.MethodPost, "/find/objectclassification", s.SearchClassification, nil)
	s.addAction(http.MethodPut, "/update/objectclassification/{id}", s.UpdateClassification, nil)
	s.addAction(http.MethodDelete, "/delete/objectclassification/{id}", s.DeleteClassification, nil)
}

func (s *Service) initBusinessObjectAttribute() {
	s.addAction(http.MethodPost, "/create/objectattr", s.CreateObjectAttribute, nil)
	s.addAction(http.MethodPost, "/find/objectattr", s.SearchObjectAttribute, nil)
	s.addAction(http.MethodPost, "/find/objectattr/host", s.ListHostModelAttribute, nil)
	s.addAction(http.MethodPut, "/update/objectattr/{id}", s.UpdateObjectAttribute, nil)
	s.addAction(http.MethodDelete, "/delete/objectattr/{id}", s.DeleteObjectAttribute, nil)
}

func (s *Service) initBusinessObjectUnique() {
	s.addAction(http.MethodPost, "/create/objectunique/object/{bk_obj_id}", s.CreateObjectUnique, nil)
	s.addAction(http.MethodPut, "/update/objectunique/object/{bk_obj_id}/unique/{id}", s.UpdateObjectUnique, nil)
	s.addAction(http.MethodPost, "/delete/objectunique/object/{bk_obj_id}/unique/{id}", s.DeleteObjectUnique, nil)
	s.addAction(http.MethodPost, "/find/objectunique/object/{bk_obj_id}", s.SearchObjectUnique, nil)
}

func (s *Service) initBusinessObjectAttrGroup() {
	s.addAction(http.MethodPost, "/create/objectattgroup", s.CreateObjectGroup, nil)
	s.addAction(http.MethodPut, "/update/objectattgroup", s.UpdateObjectGroup, nil)
	s.addAction(http.MethodDelete, "/delete/objectattgroup/{id}", s.DeleteObjectGroup, nil)
	s.addAction(http.MethodPut, "/update/objectattgroupproperty", s.UpdateObjectAttributeGroupProperty, s.ParseUpdateObjectAttributeGroupPropertyInput)
	s.addAction(http.MethodPost, "/find/objectattgroup/object/{bk_obj_id}", s.SearchGroupByObject, nil)
}

func (s *Service) initBusinessGraphics() {
	s.addAction(http.MethodPost, "/find/objecttopo/scope_type/{scope_type}/scope_id/{scope_id}", s.SelectObjectTopoGraphics, nil)
	s.addAction(http.MethodPost, "/update/objecttopo/scope_type/{scope_type}/scope_id/{scope_id}", s.UpdateObjectTopoGraphicsNew, nil)
}

func (s *Service) initBusinessAssociation() {

	// mainline topo methods
	s.addAction(http.MethodPost, "/create/topomodelmainline", s.CreateMainLineObject, nil)
	s.addAction(http.MethodDelete, "/delete/topomodelmainline/object/{bk_obj_id}", s.DeleteMainLineObject, nil)
	s.addAction(http.MethodPost, "/find/topomodelmainline", s.SearchMainLineObjectTopo, nil)
	s.addAction(http.MethodPost, "/find/topoinst/biz/{bk_biz_id}", s.SearchBusinessTopo, nil)
	s.addAction(http.MethodPost, "/find/topoinst_with_statistics/biz/{bk_biz_id}", s.SearchBusinessTopoWithStatistics, nil)
	s.addAction(http.MethodPost, "/find/topopath/biz/{bk_biz_id}", s.SearchTopoPath, nil)

	// association type methods ,NOT SUPPORT BUSINESS
	s.addAction(http.MethodPost, "/find/topoassociationtype", s.SearchObjectAssocWithAssocKindList, nil)
	s.addAction(http.MethodPost, "/find/associationtype", s.SearchAssociationType, nil)
	s.addAction(http.MethodPost, "/create/associationtype", s.CreateAssociationType, nil)
	s.addAction(http.MethodPut, "/update/associationtype/{id}", s.UpdateAssociationType, nil)
	s.addAction(http.MethodDelete, "/delete/associationtype/{id}", s.DeleteAssociationType, nil)

	// object association methods
	s.addAction(http.MethodPost, "/find/objectassociation", s.SearchObjectAssociation, nil)
	s.addAction(http.MethodPost, "/create/objectassociation", s.CreateObjectAssociation, nil)
	s.addAction(http.MethodPut, "/update/objectassociation/{id}", s.UpdateObjectAssociation, nil)
	s.addAction(http.MethodDelete, "/delete/objectassociation/{id}", s.DeleteObjectAssociation, nil)

	// inst association methods
	s.addAction(http.MethodPost, "/find/instassociation", s.SearchAssociationInst, nil)
	s.addAction(http.MethodPost, "/create/instassociation", s.CreateAssociationInst, nil)
	s.addAction(http.MethodDelete, "/delete/instassociation/{association_id}", s.DeleteAssociationInst, nil)

	// topo search methods
	s.addAction(http.MethodPost, "/find/instassociation/object/{bk_obj_id}", s.SearchInstByAssociation, nil)
	s.addAction(http.MethodPost, "/find/instassttopo/object/{bk_obj_id}/inst/{inst_id}", s.SearchInstTopo, nil)

	// ATTENTION: the following methods is not recommended
	s.addAction(http.MethodPost, "/find/insttopo/object/{bk_obj_id}/inst/{inst_id}", s.SearchInstChildTopo, nil)
	s.addAction(http.MethodPost, "/import/instassociation/{bk_obj_id}", s.ImportInstanceAssociation, nil)
}

func (s *Service) initBusinessInst() {
	s.addAction(http.MethodPost, "/create/instance/object/{bk_obj_id}", s.CreateInst, nil)
	s.addAction(http.MethodDelete, "/delete/instance/object/{bk_obj_id}/inst/{inst_id}", s.DeleteInst, nil)
	s.addAction(http.MethodDelete, "/deletemany/instance/object/{bk_obj_id}", s.DeleteInsts, nil)
	s.addAction(http.MethodPut, "/update/instance/object/{bk_obj_id}/inst/{inst_id}", s.UpdateInst, nil)
	s.addAction(http.MethodPut, "/updatemany/instance/object/{bk_obj_id}", s.UpdateInsts, nil)
	s.addAction(http.MethodPost, "/find/instance/object/{bk_obj_id}", s.SearchInstAndAssociationDetail, nil)
	s.addAction(http.MethodPost, "/find/instdetail/object/{bk_obj_id}/inst/{inst_id}", s.SearchInstByInstID, nil)
}
