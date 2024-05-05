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

// Package fieldtemplate Package fieldtmpl defines field template api machinery.
package fieldtemplate

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// FieldTemplateInterface field template interface
type FieldTemplateInterface interface {
	// field template
	ListFieldTemplate(ctx context.Context, h http.Header, opt metadata.CommonQueryOption) (
		*metadata.ListFieldTemplateResp, errors.CCErrorCoder)
	CreateFieldTemplate(ctx context.Context, header http.Header, option metadata.CreateFieldTmplOption) (
		*metadata.CreateResult, errors.CCErrorCoder)
	FindFieldTemplateByID(ctx context.Context, header http.Header, fieldTemplateID int64) (
		*metadata.FieldTemplateResp, errors.CCErrorCoder)
	FieldTemplateBindObject(ctx context.Context, header http.Header, option metadata.FieldTemplateBindObjOpt) (
		*metadata.Response, errors.CCErrorCoder)
	FieldTemplateUnbindObject(ctx context.Context, header http.Header, option metadata.FieldTemplateUnbindObjOpt) (
		*metadata.Response, errors.CCErrorCoder)
	DeleteFieldTemplate(ctx context.Context, header http.Header, option metadata.DeleteFieldTmplOption) (
		*metadata.Response, errors.CCErrorCoder)
	CloneFieldTemplate(ctx context.Context, header http.Header, option metadata.CloneFieldTmplOption) (
		*metadata.CreateResult, errors.CCErrorCoder)
	UpdateFieldTemplate(ctx context.Context, header http.Header, option metadata.UpdateFieldTmplOption) (
		*metadata.Response, errors.CCErrorCoder)
	UpdateFieldTemplateInfo(ctx context.Context, header http.Header, option metadata.FieldTemplate) (
		*metadata.Response, errors.CCErrorCoder)

	// field template attribute
	ListFieldTemplateAttr(ctx context.Context, h http.Header, opt metadata.ListFieldTmplAttrOption) (
		*metadata.ListFieldTemplateAttrResp, errors.CCErrorCoder)
	CountFieldTemplateAttr(ctx context.Context, header http.Header, option metadata.CountFieldTmplResOption) (
		*metadata.CountFieldTemplateAttrResult, errors.CCErrorCoder)

	// field template unique
	ListFieldTemplateUnique(ctx context.Context, header http.Header, option metadata.ListFieldTmplUniqueOption) (
		*metadata.ListFieldTmplUniqueResp, errors.CCErrorCoder)

	// field template sync to object
	SyncFieldTemplateInfoToObjects(ctx context.Context, header http.Header, option metadata.FieldTemplateSyncOption) (
		*metadata.Response, errors.CCErrorCoder)
	ListObjFieldTmplRel(ctx context.Context, header http.Header, option metadata.ListObjFieldTmplRelOption) (
		*metadata.ListObjFieldTmplRelResp, errors.CCErrorCoder)
	ListFieldTmplByObj(ctx context.Context, header http.Header, option metadata.ListFieldTmplByObjOption) (
		*metadata.ListFieldTemplateResp, errors.CCErrorCoder)
	ListObjByFieldTmpl(ctx context.Context, header http.Header, option metadata.ListObjByFieldTmplOption) (
		*metadata.ReadModelResult, errors.CCErrorCoder)

	// field template task
	SyncFieldTemplateToObjectTask(ctx context.Context, header http.Header, option metadata.SyncObjectTask) (
		*metadata.Response, errors.CCErrorCoder)

	// compare field template with object
	CompareFieldTemplateAttr(ctx context.Context, header http.Header, option metadata.CompareFieldTmplAttrOption) (
		*metadata.CompareFieldTmplAttrsResResp, errors.CCErrorCoder)
	CompareFieldTemplateUnique(ctx context.Context, header http.Header, option metadata.CompareFieldTmplUniqueOption) (
		*metadata.CompareFieldTmplUniquesResResp, errors.CCErrorCoder)

	ListFieldTemplateTasksStatus(ctx context.Context, header http.Header,
		option metadata.ListFieldTmplTaskStatusOption) (
		*metadata.ListFieldTmplTaskStatusResultResp, errors.CCErrorCoder)
	ListFieldTemplateSyncStatus(ctx context.Context, header http.Header,
		option metadata.ListFieldTmpltSyncStatusOption) (
		*metadata.ListFieldTmpltSyncStatusResultResp, errors.CCErrorCoder)
	ListFieldTmplByUniqueTmplIDForUI(ctx context.Context, header http.Header,
		option metadata.ListTmplSimpleByUniqueOption) (*metadata.ListFieldTemplateSimpleResp, errors.CCErrorCoder)
	ListFieldTmplByObjectTmplIDForUI(ctx context.Context, header http.Header,
		option metadata.ListTmplSimpleByAttrOption) (*metadata.ListFieldTemplateSimpleResp, errors.CCErrorCoder)
	ListFieldTemplateModelStatus(ctx context.Context, header http.Header,
		option metadata.ListFieldTmplModelStatusOption) (*metadata.ListFieldTmplTaskSyncResultResp, errors.CCErrorCoder)
}

// NewFieldTemplateInterface TODO
func NewFieldTemplateInterface(client rest.ClientInterface) FieldTemplateInterface {
	return &FieldTemplate{client: client}
}

// FieldTemplate TODO
type FieldTemplate struct {
	client rest.ClientInterface
}
