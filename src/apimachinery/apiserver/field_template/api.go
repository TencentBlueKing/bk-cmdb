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
	ListFieldTemplateAttr(ctx context.Context, h http.Header, opt *metadata.ListFieldTmplAttrOption) (
		*metadata.FieldTemplateAttrInfo, errors.CCErrorCoder)
	ListObjFieldTmplRel(ctx context.Context, h http.Header, opt *metadata.ListObjFieldTmplRelOption) (
		*metadata.ObjFieldTmplRelInfo, errors.CCErrorCoder)
	CountFieldTemplateObj(ctx context.Context, h http.Header, opt *metadata.CountFieldTmplResOption) (
		[]metadata.FieldTmplResCount, errors.CCErrorCoder)
}

// New field template api client.
func New(client rest.ClientInterface) Interface {
	return &template{client: client}
}

type template struct {
	client rest.ClientInterface
}
