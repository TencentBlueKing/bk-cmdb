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

//import (
//	"net/http"

//	"configcenter/src/common"
//)

//func (s *topoService) initBusinessObject() {

//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/createmany/object", HandlerFunc: s.CreateObjectBatch})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/findmany/object", HandlerFunc: s.SearchObjectBatch})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/create/object", HandlerFunc: s.CreateOneObject})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/find/object", HandlerFunc: s.SearchObject})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/update/object", HandlerFunc: s.UpdateObject})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/delete/object", HandlerFunc: s.DeleteObject})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/find/object-topology", HandlerFunc: s.SearchObjectTopo})

//}

//func (s *topoService) initBusinessClassification() {
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classification", HandlerFunc: s.CreateClassification})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classification/{owner_id}/objects", HandlerFunc: s.SearchClassificationWithObjects})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/classifications", HandlerFunc: s.SearchClassification})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/object/classification/{id}", HandlerFunc: s.UpdateClassification})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/object/classification/{id}", HandlerFunc: s.DeleteClassification})

//}

//func (s *topoService) initBusinessObjectAttribute() {
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectattr", HandlerFunc: s.CreateObjectAttribute})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectattr/search", HandlerFunc: s.SearchObjectAttribute})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectattr/{id}", HandlerFunc: s.UpdateObjectAttribute})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectattr/{id}", HandlerFunc: s.DeleteObjectAttribute})
//}

//func (s *topoService) initBusinessObjectUnique() {
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/{bk_obj_id}/unique/action/create", HandlerFunc: s.CreateObjectUnique})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/object/{bk_obj_id}/unique/{id}/action/update", HandlerFunc: s.UpdateObjectUnique})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/object/{bk_obj_id}/unique/{id}/action/delete", HandlerFunc: s.DeleteObjectUnique})
//	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/object/{bk_obj_id}/unique/action/search", HandlerFunc: s.SearchObjectUnique})
//}

//func (s *topoService) initBusinessObjectAttrGroup() {
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectatt/group/new", HandlerFunc: s.CreateObjectGroup})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectatt/group/update", HandlerFunc: s.UpdateObjectGroup})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/groupid/{id}", HandlerFunc: s.DeleteObjectGroup})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/objectatt/group/property", HandlerFunc: s.UpdateObjectAttributeGroup, HandlerParseOriginDataFunc: s.ParseUpdateObjectAttributeGroupInput})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", HandlerFunc: s.DeleteObjectAttributeGroup})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", HandlerFunc: s.SearchGroupByObject})
//}

//func (s *topoService) initBusinessAssociation() {

//	// mainline topo methods
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/model/mainline", HandlerFunc: s.CreateMainLineObject})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/topo/model/mainline/owners/{owner_id}/objectids/{obj_id}", HandlerFunc: s.DeleteMainLineObject})
//	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/model/{owner_id}", HandlerFunc: s.SearchMainLineObjectTopo})
//	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/model/{owner_id}/{cls_id}/{obj_id}", HandlerFunc: s.SearchObjectByClassificationID})
//	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/inst/{owner_id}/{app_id}", HandlerFunc: s.SearchBusinessTopo})
//	// TODO: delete this api, it's not used by front.
//	s.actions = append(s.actions, action{Method: http.MethodGet, Path: "/topo/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", HandlerFunc: s.SearchMainLineChildInstTopo})

//	// association type methods
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/association/type/action/search/batch", HandlerFunc: s.SearchObjectAssoWithAssoKindList})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/association/type/action/search", HandlerFunc: s.SearchAssociationType})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/topo/association/type/action/create", HandlerFunc: s.CreateAssociationType})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/topo/association/type/{id}/action/update", HandlerFunc: s.UpdateAssociationType})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/topo/association/type/{id}/action/delete", HandlerFunc: s.DeleteAssociationType})

//	// object association methods
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/association/action/search", HandlerFunc: s.SearchObjectAssociation})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/object/association/action/create", HandlerFunc: s.CreateObjectAssociation})
//	s.actions = append(s.actions, action{Method: http.MethodPut, Path: "/object/association/{id}/action/update", HandlerFunc: s.UpdateObjectAssociation})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/object/association/{id}/action/delete", HandlerFunc: s.DeleteObjectAssociation})

//	// inst association methods
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/action/search", HandlerFunc: s.SearchAssociationInst})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/action/create", HandlerFunc: s.CreateAssociationInst})
//	s.actions = append(s.actions, action{Method: http.MethodDelete, Path: "/inst/association/{association_id}/action/delete", HandlerFunc: s.DeleteAssociationInst})

//	// topo search methods
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: s.SearchInstByAssociation})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: s.SearchInstTopo})

//	// ATTENTION: the following methods is not recommended
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: s.SearchInstChildTopo})
//	s.actions = append(s.actions, action{Method: http.MethodPost, Path: "/inst/association/action/{obj_id}/import", HandlerFunc: s.ImportInstanceAssociation})

//}
