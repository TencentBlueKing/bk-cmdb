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

package model

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (m *model) CreateManyModelClassification(ctx context.Context, h http.Header, input *metadata.CreateManyModelClassifiaction) (resp *metadata.CreatedManyOptionResult, err error) {
	resp = new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/model/classification"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) CreateModelClassification(ctx context.Context, h http.Header, input *metadata.CreateOneModelClassification) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/model/classification"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) SetManyModelClassification(ctx context.Context, h http.Header, input *metadata.SetManyModelClassification) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/setmany/model/classification"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) SetModelClassification(ctx context.Context, h http.Header, input *metadata.SetOneModelClassification) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/model/classification"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) UpdateModelClassification(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/model/classification"

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) DeleteModelClassification(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/model/classification"

	err = m.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) ReadModelClassification(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelClassifitionResult, err error) {
	resp = new(metadata.ReadModelClassifitionResult)
	subPath := "/read/model/classification"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) CreateModel(ctx context.Context, h http.Header, input *metadata.CreateModel) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/model"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) SetModel(ctx context.Context, h http.Header, input *metadata.SetModel) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/set/model"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) UpdateModel(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/model"

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) DeleteModel(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/model"

	err = m.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) DeleteModelCascade(ctx context.Context, h http.Header, modelID int64) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := fmt.Sprintf("/delete/model/%d/cascade", modelID)

	err = m.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) ReadModel(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelResult, err error) {
	resp = new(metadata.ReadModelResult)
	subPath := "/read/model"

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) CreateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelAttributes) (resp *metadata.CreatedManyOptionResult, err error) {
	resp = new(metadata.CreatedManyOptionResult)
	subPath := fmt.Sprintf("/create/model/%s/attributes", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) UpdateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := fmt.Sprintf("/update/model/%s/attributes", objID)

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) UpdateModelAttrsByCondition(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := fmt.Sprintf("/update/model/attributes")

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) SetModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.SetModelAttributes) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := fmt.Sprintf("/set/model/%s/attributes", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) DeleteModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := fmt.Sprintf("/delete/model/%s/attributes", objID)

	err = m.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) ReadModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (resp *metadata.ReadModelAttrResult, err error) {
	resp = new(metadata.ReadModelAttrResult)
	subPath := fmt.Sprintf("/read/model/%s/attributes", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) ReadModelAttrByCondition(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelAttrResult, err error) {
	resp = new(metadata.ReadModelAttrResult)
	subPath := fmt.Sprintf("/read/model/attributes")

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *model) ReadAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.QueryCondition) (resp metadata.ReadModelAttributeGroupResult, err error) {
	subPath := fmt.Sprintf("/read/model/%s/group", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) ReadAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.QueryCondition) (resp metadata.ReadModelAttributeGroupResult, err error) {
	subPath := fmt.Sprintf("/read/model/group")

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) CreateAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.CreateModelAttributeGroup) (resp metadata.CreatedOneOptionResult, err error) {
	subPath := fmt.Sprintf("/create/model/%s/group", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) UpdateAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.UpdateOption) (resp metadata.UpdatedOptionResult, err error) {
	subPath := fmt.Sprintf("/update/model/%s/group", objID)

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) UpdateAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.UpdateOption) (resp metadata.UpdatedOptionResult, err error) {
	subPath := fmt.Sprintf("/update/model/group")

	err = m.client.Put().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) SetAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.SetModelAttributes) (resp metadata.SetOptionResult, err error) {
	subPath := fmt.Sprintf("/set/model/%s/group", objID)

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) DeleteAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.DeleteOption) (resp metadata.DeletedOptionResult, err error) {
	subPath := fmt.Sprintf("/delete/model/%s/group", objID)

	err = m.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) DeleteAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.DeleteOption) (resp metadata.DeletedOptionResult, err error) {
	subPath := fmt.Sprintf("/delete/model/group")

	err = m.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) CreateModelAttrUnique(ctx context.Context, h http.Header, objID string, data metadata.CreateModelAttrUnique) (resp *metadata.CreatedOneOptionResult, err error) {
	subPath := fmt.Sprintf("/create/model/%s/attributes/unique", objID)
	err = m.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) UpdateModelAttrUnique(ctx context.Context, h http.Header, objID string, id uint64, data metadata.UpdateModelAttrUnique) (resp *metadata.UpdatedOptionResult, err error) {
	subPath := fmt.Sprintf("/update/model/%s/attributes/unique/%d", objID, id)

	err = m.client.Put().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) DeleteModelAttrUnique(ctx context.Context, h http.Header, objID string, id uint64, meta metadata.DeleteModelAttrUnique) (resp *metadata.DeletedOptionResult, err error) {
	subPath := fmt.Sprintf("/delete/model/%s/attributes/unique/%d", objID, id)

	err = m.client.Delete().
		WithContext(ctx).
		Body(meta).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

func (m *model) ReadModelAttrUnique(ctx context.Context, h http.Header, inputParam metadata.QueryCondition) (resp *metadata.ReadModelUniqueResult, err error) {
	subPath := fmt.Sprintf("/read/model/attributes/unique")

	err = m.client.Post().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Body(inputParam).
		Do().
		Into(&resp)
	return
}

// GetModelStatistics 统计各个模型的实例数
func (m *model) GetModelStatistics(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/read/model/statistics"

	err = m.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
