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

// Package fieldtmpl defines field template api machinery.
package fieldtmpl

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// Interface defines field template apis.
type Interface interface {
	ListFieldTemplate(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
		*metadata.FieldTemplateInfo, errors.CCErrorCoder)
	FieldTemplateBindObject(ctx context.Context, h http.Header,
		opt *metadata.FieldTemplateBindObjOpt) errors.CCErrorCoder
	FieldTemplateUnbindObject(ctx context.Context, h http.Header,
		opt *metadata.FieldTemplateUnbindObjOpt) errors.CCErrorCoder
	ListFieldTemplateAttr(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
		*metadata.FieldTemplateAttrInfo, errors.CCErrorCoder)
	ListFieldTemplateUnique(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
		*metadata.FieldTemplateUniqueInfo, errors.CCErrorCoder)
	CreateFieldTemplate(ctx context.Context, h http.Header, opt *metadata.FieldTemplate) (*metadata.RspID,
		errors.CCErrorCoder)
	CreateFieldTemplateAttrs(ctx context.Context, h http.Header, templateID int64, opt []metadata.FieldTemplateAttr) (
		*metadata.RspIDs, errors.CCErrorCoder)
	CreateFieldTemplateUniques(ctx context.Context, h http.Header, templateID int64,
		opt []metadata.FieldTemplateUnique) (*metadata.RspIDs, errors.CCErrorCoder)
	ListObjFieldTmplRel(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
		*metadata.ObjFieldTmplRelInfo, errors.CCErrorCoder)
	DeleteFieldTemplate(ctx context.Context, h http.Header, opt *metadata.DeleteOption) errors.CCErrorCoder
	DeleteFieldTemplateAttr(ctx context.Context, h http.Header, templateID int64,
		opt *metadata.DeleteOption) errors.CCErrorCoder
	DeleteFieldTemplateUnique(ctx context.Context, h http.Header, templateID int64,
		opt *metadata.DeleteOption) errors.CCErrorCoder
	UpdateFieldTemplate(ctx context.Context, h http.Header, opt *metadata.FieldTemplate) errors.CCErrorCoder
	UpdateFieldTemplateAttrs(ctx context.Context, h http.Header, templateID int64,
		opt []metadata.FieldTemplateAttr) errors.CCErrorCoder
	UpdateFieldTemplateUniques(ctx context.Context, h http.Header, templateID int64,
		opt []metadata.FieldTemplateUnique) errors.CCErrorCoder
	ListFieldTmplSimplyByUniqueTemplateID(ctx context.Context, h http.Header,
		opt *metadata.ListTmplSimpleByUniqueOption) (*metadata.ListTmplSimpleResult, errors.CCErrorCoder)
	ListFieldTmplSimplyByAttrTemplateID(ctx context.Context, h http.Header,
		opt *metadata.ListTmplSimpleByAttrOption) (*metadata.ListTmplSimpleResult, errors.CCErrorCoder)
}

// New field template api client.
func New(client rest.ClientInterface) Interface {
	return &template{client: client}
}

type template struct {
	client rest.ClientInterface
}
