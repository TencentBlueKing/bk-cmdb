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

func (asst *association) CreateAssociationType(ctx context.Context, h http.Header, input *metadata.CreateAssociationKind) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) CreateManyAssociation(ctx context.Context, h http.Header, input *metadata.CreateManyAssociationKind) (resp *metadata.CreatedManyOptionResult, err error) {
	resp = new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) SetAssociation(ctx context.Context, h http.Header, input *metadata.SetAssociationKind) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) SetManyAssociation(ctx context.Context, h http.Header, input *metadata.SetManyAssociationKind) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/setmany/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) UpdateAssociationType(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/associationkind"

	err = asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) DeleteAssociationType(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/associationkind"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) DeleteAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/associationkind/cascade"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) ReadAssociationType(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.SearchAssociationTypeResult, err error) {
	resp = new(metadata.SearchAssociationTypeResult)
	subPath := "/read/associationkind"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) CreateModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/modelassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) CreateMainlineModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/mainlinemodelassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) SetModelAssociation(ctx context.Context, h http.Header, input *metadata.SetModelAssociation) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/modelassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) UpdateModelAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/modelassociation"

	err = asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) ReadModelAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelAssociationResult, err error) {
	resp = new(metadata.ReadModelAssociationResult)
	subPath := "/read/modelassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) DeleteModelAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/modelassociation"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) DeleteModelAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/modelassociation/cascade"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) CreateInstAssociation(ctx context.Context, h http.Header, input *metadata.CreateOneInstanceAssociation) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/instanceassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) SetInstAssociation(ctx context.Context, h http.Header, input *metadata.SetOneInstanceAssociation) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/instanceassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) UpdateInstAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/instanceassociation"

	err = asst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) ReadInstAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadInstAssociationResult, err error) {
	resp = new(metadata.ReadInstAssociationResult)
	subPath := "/read/instanceassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (asst *association) DeleteInstAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/instanceassociation"

	err = asst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
