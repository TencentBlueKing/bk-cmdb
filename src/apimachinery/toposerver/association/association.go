// Package association TODO
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
	"net/http"

	"configcenter/src/common/metadata"
)

// SearchType TODO
func (asst *Association) SearchType(ctx context.Context, h http.Header, request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error) {
	resp = new(metadata.SearchAssociationTypeResult)
	subPath := "/find/associationtype"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// CreateType TODO
func (asst *Association) CreateType(ctx context.Context, h http.Header, request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error) {
	resp = new(metadata.CreateAssociationTypeResult)
	subPath := "/create/associationtype"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// UpdateType TODO
func (asst *Association) UpdateType(ctx context.Context, h http.Header, asstTypeID int, request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error) {
	resp = new(metadata.UpdateAssociationTypeResult)
	subPath := "/update/associationtype/%d"

	err = asst.client.Put().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath, asstTypeID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteType TODO
func (asst *Association) DeleteType(ctx context.Context, h http.Header, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult, err error) {
	resp = new(metadata.DeleteAssociationTypeResult)
	subPath := "/delete/associationtype/%d"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, asstTypeID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchObject TODO
func (asst *Association) SearchObject(ctx context.Context, h http.Header, request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	resp = new(metadata.SearchAssociationObjectResult)
	subPath := "/find/objectassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// CreateObject TODO
func (asst *Association) CreateObject(ctx context.Context, h http.Header, request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error) {
	resp = new(metadata.CreateAssociationObjectResult)
	subPath := "/create/objectassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// UpdateObject TODO
func (asst *Association) UpdateObject(ctx context.Context, h http.Header, asstID int, request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error) {
	resp = new(metadata.UpdateAssociationObjectResult)
	subPath := "/update/objectassociation/%d"

	err = asst.client.Put().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath, asstID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// DeleteObject TODO
func (asst *Association) DeleteObject(ctx context.Context, h http.Header, asstID int) (resp *metadata.DeleteAssociationObjectResult, err error) {
	resp = new(metadata.DeleteAssociationObjectResult)
	subPath := "/delete/objectassociation/%d"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, asstID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchInst TODO
func (asst *Association) SearchInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	resp = new(metadata.SearchAssociationInstResult)
	subPath := "/find/instassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchAssociationRelatedInst TODO
func (asst *Association) SearchAssociationRelatedInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationRelatedInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	resp = new(metadata.SearchAssociationInstResult)
	subPath := "/find/instassociation/related"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// CreateInst TODO
func (asst *Association) CreateInst(ctx context.Context, h http.Header, request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error) {
	resp = new(metadata.CreateAssociationInstResult)
	subPath := "/create/instassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// CreateManyInstAssociation TODO
func (asst *Association) CreateManyInstAssociation(ctx context.Context, header http.Header,
	request *metadata.CreateManyInstAsstRequest) (*metadata.CreateManyInstAsstResult, error) {
	resp := new(metadata.CreateManyInstAsstResult)
	subPath := "/createmany/instassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	return resp, err
}

// DeleteInst TODO
func (asst *Association) DeleteInst(ctx context.Context, h http.Header, objID string, assoID int64) (resp *metadata.DeleteAssociationInstResult, err error) {
	resp = new(metadata.DeleteAssociationInstResult)
	subPath := "/delete/instassociation/%s/%d"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID, assoID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// DeleteInstBatch TODO
func (asst *Association) DeleteInstBatch(ctx context.Context, h http.Header, assoIDs *metadata.DeleteAssociationInstBatchRequest) (resp *metadata.DeleteAssociationInstBatchResult, err error) {
	resp = new(metadata.DeleteAssociationInstBatchResult)
	subPath := "/delete/instassociation/batch"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(assoIDs).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchObjectAssoWithAssoKindList TODO
func (asst *Association) SearchObjectAssoWithAssoKindList(ctx context.Context, h http.Header, assoKindIDs metadata.AssociationKindIDs) (resp *metadata.ListAssociationsWithAssociationKindResult, err error) {
	resp = new(metadata.ListAssociationsWithAssociationKindResult)
	subPath := "/find/topoassociationtype"

	err = asst.client.Post().
		WithContext(ctx).
		Body(assoKindIDs).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchInstanceAssociations is search instance associations api in toposerver.
func (asst *Association) SearchInstanceAssociations(ctx context.Context, header http.Header,
	objID string, input *metadata.CommonSearchFilter) (*metadata.Response, error) {

	resp := new(metadata.Response)
	subPath := "/search/instance_associations/object/%s"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(header).
		Do().
		Into(resp)

	return resp, err
}

// CountInstanceAssociations is count instance associations api in toposerver.
func (asst *Association) CountInstanceAssociations(ctx context.Context, header http.Header,
	objID string, input *metadata.CommonCountFilter) (*metadata.Response, error) {

	resp := new(metadata.Response)
	subPath := "/count/instance_associations/object/%s"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(header).
		Do().
		Into(resp)

	return resp, err
}

// SearchInstAssocAndInstDetail is search instance associations and inst detail api in toposerver.
func (asst *Association) SearchInstAssocAndInstDetail(ctx context.Context, header http.Header, objID string,
	input *metadata.InstAndAssocRequest) (*metadata.InstAndAssocDetailResult, error) {

	resp := new(metadata.InstAndAssocDetailResult)
	subPath := "/find/instassociation/object/%s/inst/detail"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}
	return resp, nil
}
