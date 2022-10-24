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

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateAssociationType api of create association type
func (asst *association) CreateAssociationType(ctx context.Context, h http.Header,
	input *metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error) {

	resp := new(metadata.CreatedOneOptionResult)
	subPath := "/create/associationkind"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateManyAssociation TODO
func (asst *association) CreateManyAssociation(ctx context.Context, h http.Header,
	input *metadata.CreateManyAssociationKind) (*metadata.CreateManyDataResult, error) {

	resp := new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/associationkind"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// SetAssociation TODO
func (asst *association) SetAssociation(ctx context.Context, h http.Header, input *metadata.SetAssociationKind) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SetManyAssociation TODO
func (asst *association) SetManyAssociation(ctx context.Context, h http.Header, input *metadata.SetManyAssociationKind) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/setmany/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateAssociationType api of update association type
func (asst *association) UpdateAssociationType(ctx context.Context, h http.Header, input *metadata.UpdateOption) (
	*metadata.UpdatedCount, error) {

	resp := new(metadata.UpdatedOptionResult)
	subPath := "/update/associationkind"

	err := asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteAssociationType api of delete association type
func (asst *association) DeleteAssociationType(ctx context.Context, h http.Header, input *metadata.DeleteOption) (
	*metadata.DeletedCount, error) {

	resp := new(metadata.DeletedOptionResult)
	subPath := "/delete/associationkind"

	err := asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteAssociationCascade TODO
func (asst *association) DeleteAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/associationkind/cascade"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ReadAssociationType api of search association type
func (asst *association) ReadAssociationType(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
	*metadata.SearchAssociationType, error) {

	resp := new(metadata.SearchAssociationTypeResult)
	subPath := "/read/associationkind"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateModelAssociation api of create object association
func (asst *association) CreateModelAssociation(ctx context.Context, h http.Header,
	input *metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {

	resp := new(metadata.CreatedOneOptionResult)
	subPath := "/create/modelassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateMainlineModelAssociation api of create mainline object association
func (asst *association) CreateMainlineModelAssociation(ctx context.Context, h http.Header,
	input *metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {

	resp := new(metadata.CreatedOneOptionResult)
	subPath := "/create/mainlinemodelassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// SetModelAssociation TODO
func (asst *association) SetModelAssociation(ctx context.Context, h http.Header, input *metadata.SetModelAssociation) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/modelassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateModelAssociation api of update object association
func (asst *association) UpdateModelAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (
	*metadata.UpdatedCount, error) {

	resp := new(metadata.UpdatedOptionResult)
	subPath := "/update/modelassociation"

	err := asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// ReadModelAssociation api of search object association
func (asst *association) ReadModelAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
	*metadata.QueryModelAssociationResult, error) {

	resp := new(metadata.ReadModelAssociationResult)
	subPath := "/read/modelassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteModelAssociation api of delete object association
func (asst *association) DeleteModelAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (
	*metadata.DeletedCount, error) {

	resp := new(metadata.DeletedOptionResult)
	subPath := "/delete/modelassociation"

	err := asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteModelAssociationCascade TODO
func (asst *association) DeleteModelAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/modelassociation/cascade"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CreateInstAssociation api of create object instances association
func (asst *association) CreateInstAssociation(ctx context.Context, h http.Header,
	input *metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {

	resp := new(metadata.CreatedOneOptionResult)
	subPath := "/create/instanceassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateManyInstAssociation api of batch create object instances association
func (asst *association) CreateManyInstAssociation(ctx context.Context, header http.Header,
	input *metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error) {

	resp := new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/instanceassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// SetInstAssociation TODO
func (asst *association) SetInstAssociation(ctx context.Context, h http.Header, input *metadata.SetOneInstanceAssociation) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/instanceassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateInstAssociation TODO
func (asst *association) UpdateInstAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/instanceassociation"

	err = asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ReadInstAssociation api of search object instances association
func (asst *association) ReadInstAssociation(ctx context.Context, h http.Header,
	input *metadata.InstAsstQueryCondition) (*metadata.QueryInstAssociationResult, error) {

	resp := new(metadata.ReadInstAssociationResult)
	subPath := "/read/instanceassociation"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CountInstanceAssociations counts model instance associations num.
func (asst *association) CountInstanceAssociations(ctx context.Context, header http.Header, objID string,
	input *metadata.Condition) (*metadata.CountResponseContent, error) {

	resp := new(metadata.CountResponse)
	subPath := "/count/instanceassociation/model/%s"

	err := asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteInstAssociation api of delete object instances associationS
func (asst *association) DeleteInstAssociation(ctx context.Context, h http.Header,
	input *metadata.InstAsstDeleteOption) (*metadata.DeletedCount, error) {

	resp := new(metadata.DeletedOptionResult)
	subPath := "/delete/instanceassociation"

	err := asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
