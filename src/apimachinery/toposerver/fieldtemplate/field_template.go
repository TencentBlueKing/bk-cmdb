/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package fieldtemplate

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ListFieldTemplate list field templates
func (ft FieldTemplate) ListFieldTemplate(ctx context.Context, h http.Header, opt metadata.CommonQueryOption) (
	*metadata.ListFieldTemplateResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTemplateResp)
	subPath := "/findmany/field_template"

	err := ft.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// CreateFieldTemplate create field template
func (ft FieldTemplate) CreateFieldTemplate(ctx context.Context, header http.Header,
	option metadata.CreateFieldTmplOption) (*metadata.CreateResult, errors.CCErrorCoder) {

	ret := new(metadata.CreateResult)
	subPath := "/create/field_template"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// FindFieldTemplateByID find field template by id
func (ft FieldTemplate) FindFieldTemplateByID(ctx context.Context, header http.Header, fieldTemplateID int64) (
	*metadata.FieldTemplateResp, errors.CCErrorCoder) {

	ret := new(metadata.FieldTemplateResp)
	subPath := "/find/field_template/%d"

	err := ft.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, fieldTemplateID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// FieldTemplateBindObject field template bind object
func (ft FieldTemplate) FieldTemplateBindObject(ctx context.Context, header http.Header,
	option metadata.FieldTemplateBindObjOpt) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/update/field_template/bind/object"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// FieldTemplateUnbindObject field template unbind object
func (ft FieldTemplate) FieldTemplateUnbindObject(ctx context.Context, header http.Header,
	option metadata.FieldTemplateUnbindObjOpt) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/update/field_template/unbind/object"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// DeleteFieldTemplate delete field template
func (ft FieldTemplate) DeleteFieldTemplate(ctx context.Context, header http.Header,
	option metadata.DeleteFieldTmplOption) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/delete/field_template"

	err := ft.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// CloneFieldTemplate clone field template
func (ft FieldTemplate) CloneFieldTemplate(ctx context.Context, header http.Header,
	option metadata.CloneFieldTmplOption) (*metadata.CreateResult, errors.CCErrorCoder) {

	ret := new(metadata.CreateResult)
	subPath := "/create/field_template/clone"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// UpdateFieldTemplate update field template
func (ft FieldTemplate) UpdateFieldTemplate(ctx context.Context, header http.Header,
	option metadata.UpdateFieldTmplOption) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/update/field_template"

	err := ft.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// UpdateFieldTemplateInfo update field template info
func (ft FieldTemplate) UpdateFieldTemplateInfo(ctx context.Context, header http.Header,
	option metadata.FieldTemplate) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/update/field_template/info"

	err := ft.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTemplateAttr list field template attributes
func (ft FieldTemplate) ListFieldTemplateAttr(ctx context.Context, h http.Header,
	opt metadata.ListFieldTmplAttrOption) (*metadata.ListFieldTemplateAttrResp, errors.CCErrorCoder) {

	resp := new(metadata.ListFieldTemplateAttrResp)
	subPath := "/findmany/field_template/attribute"

	err := ft.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp, nil
}

// CountFieldTemplateAttr count field template attribute
func (ft FieldTemplate) CountFieldTemplateAttr(ctx context.Context, header http.Header,
	option metadata.CountFieldTmplResOption) (*metadata.CountFieldTemplateAttrResult, errors.CCErrorCoder) {

	ret := new(metadata.CountFieldTemplateAttrResult)
	subPath := "/findmany/field_template/attribute/count"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTemplateUnique list field template unique
func (ft FieldTemplate) ListFieldTemplateUnique(ctx context.Context, header http.Header,
	option metadata.ListFieldTmplUniqueOption) (*metadata.ListFieldTmplUniqueResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTmplUniqueResp)
	subPath := "/findmany/field_template/unique"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// SyncFieldTemplateInfoToObjects sync field template info to objects
func (ft FieldTemplate) SyncFieldTemplateInfoToObjects(ctx context.Context, header http.Header,
	option metadata.FieldTemplateSyncOption) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/update/topo/field_template/sync"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListObjFieldTmplRel list object and field template relation
func (ft FieldTemplate) ListObjFieldTmplRel(ctx context.Context, header http.Header,
	option metadata.ListObjFieldTmplRelOption) (*metadata.ListObjFieldTmplRelResp, errors.CCErrorCoder) {

	ret := new(metadata.ListObjFieldTmplRelResp)
	subPath := "/findmany/field_template/object/relation"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTmplByObj list field template by object
func (ft FieldTemplate) ListFieldTmplByObj(ctx context.Context, header http.Header,
	option metadata.ListFieldTmplByObjOption) (*metadata.ListFieldTemplateResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTemplateResp)
	subPath := "/findmany/field_template/by_object"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListObjByFieldTmpl list object by field template
func (ft FieldTemplate) ListObjByFieldTmpl(ctx context.Context, header http.Header,
	option metadata.ListObjByFieldTmplOption) (*metadata.ReadModelResult, errors.CCErrorCoder) {

	ret := new(metadata.ReadModelResult)
	subPath := "/findmany/object/by_field_template"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// SyncFieldTemplateToObjectTask sync field template to object task
func (ft FieldTemplate) SyncFieldTemplateToObjectTask(ctx context.Context, header http.Header,
	option metadata.SyncObjectTask) (*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)
	subPath := "/sync/field_template/object/task"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// CompareFieldTemplateAttr compare field template attribute
func (ft FieldTemplate) CompareFieldTemplateAttr(ctx context.Context, header http.Header,
	option metadata.CompareFieldTmplAttrOption) (*metadata.CompareFieldTmplAttrsResResp, errors.CCErrorCoder) {

	ret := new(metadata.CompareFieldTmplAttrsResResp)
	subPath := "/find/field_template/attribute/difference"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// CompareFieldTemplateUnique compare field template unique
func (ft FieldTemplate) CompareFieldTemplateUnique(ctx context.Context, header http.Header,
	option metadata.CompareFieldTmplUniqueOption) (*metadata.CompareFieldTmplUniquesResResp, errors.CCErrorCoder) {

	ret := new(metadata.CompareFieldTmplUniquesResResp)
	subPath := "/find/field_template/unique/difference"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTemplateTasksStatus list field template task status
func (ft FieldTemplate) ListFieldTemplateTasksStatus(ctx context.Context, header http.Header,
	option metadata.ListFieldTmplTaskStatusOption) (*metadata.ListFieldTmplTaskStatusResultResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTmplTaskStatusResultResp)
	subPath := "/find/field_template/tasks_status"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTemplateSyncStatus list field template sync status
func (ft FieldTemplate) ListFieldTemplateSyncStatus(ctx context.Context, header http.Header,
	option metadata.ListFieldTmpltSyncStatusOption) (*metadata.ListFieldTmpltSyncStatusResultResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTmpltSyncStatusResultResp)
	subPath := "/find/field_template/sync/status"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTmplByUniqueTmplIDForUI list field template by unique template for ui
func (ft FieldTemplate) ListFieldTmplByUniqueTmplIDForUI(ctx context.Context, header http.Header,
	option metadata.ListTmplSimpleByUniqueOption) (*metadata.ListFieldTemplateSimpleResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTemplateSimpleResp)
	subPath := "/find/field_template/simplify/by_unique_template_id"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTmplByObjectTmplIDForUI list field template by attribute template for ui
func (ft FieldTemplate) ListFieldTmplByObjectTmplIDForUI(ctx context.Context, header http.Header,
	option metadata.ListTmplSimpleByAttrOption) (*metadata.ListFieldTemplateSimpleResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTemplateSimpleResp)
	subPath := "/find/field_template/simplify/by_attr_template_id"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListFieldTemplateModelStatus list field template model status
func (ft FieldTemplate) ListFieldTemplateModelStatus(ctx context.Context, header http.Header,
	option metadata.ListFieldTmplModelStatusOption) (*metadata.ListFieldTmplTaskSyncResultResp, errors.CCErrorCoder) {

	ret := new(metadata.ListFieldTmplTaskSyncResultResp)
	subPath := "/find/field_template/model/status"

	err := ft.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := ret.CCError(); err != nil {
		return nil, err
	}

	return ret, nil
}
