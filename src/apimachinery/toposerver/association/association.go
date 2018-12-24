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
package association

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (asst *Association) SearchType(ctx context.Context, h http.Header, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
	resp = new(metadata.SearchAssociationTypeResult)
	subPath := fmt.Sprintf("/topo/association/type/action/search")

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) CreateType(ctx context.Context, h http.Header, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {
	resp = new(metadata.CreateAssociationTypeResult)
	subPath := "/topo/association/type/action/create"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
func (asst *Association) UpdateType(ctx context.Context, h http.Header, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	resp = new(metadata.UpdateAssociationTypeResult)
	subPath := fmt.Sprintf("/topo/association/type/%d/action/update", asstTypeID)

	err = asst.client.Put().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *Association) DeleteType(ctx context.Context, h http.Header, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error) {
	resp = new(metadata.DeleteAssociationTypeResult)
	subPath := fmt.Sprintf("/topo/association/type/%d/action/delete", asstTypeID)

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) SearchObject(ctx context.Context, h http.Header, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	resp = new(metadata.SearchAssociationObjectResult)
	subPath := "/object/association/action/search"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) CreateObject(ctx context.Context, h http.Header, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	resp = new(metadata.CreateAssociationObjectResult)
	subPath := "/object/association/action/create"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) UpdateObject(ctx context.Context, h http.Header, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	resp = new(metadata.UpdateAssociationObjectResult)
	subPath := fmt.Sprintf("/object/association/%d/action/update", asstID)

	err = asst.client.Put().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) DeleteObject(ctx context.Context, h http.Header, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {
	resp = new(metadata.DeleteAssociationObjectResult)
	subPath := fmt.Sprintf("/object/association/%d/action/delete", asstID)

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
func (asst *Association) SearchInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	resp = new(metadata.SearchAssociationInstResult)
	subPath := "/inst/association/action/search"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
func (asst *Association) CreateInst(ctx context.Context, h http.Header, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	resp = new(metadata.CreateAssociationInstResult)
	subPath := "/inst/association/action/create"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
func (asst *Association) DeleteInst(ctx context.Context, h http.Header, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error) {
	resp = new(metadata.DeleteAssociationInstResult)
	subPath := fmt.Sprintf("/inst/association/%d/action/delete", assoID)

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (asst *Association) SearchObjectAssoWithAssoKindList(ctx context.Context, h http.Header, assoKindIDs metadata.AssociationKindIDs) (resp *metadata.ListAssociationsWithAssociationKindResult, err error) {
	resp = new(metadata.ListAssociationsWithAssociationKindResult)
	subPath := "/topo/association/type/action/search/batch"

	err = asst.client.Post().
		WithContext(ctx).
		Body(assoKindIDs).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
