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

// Package modelquote defines model quote api machinery.
package modelquote

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// Interface defines model quote apis.
type Interface interface {
	BatchCreateQuotedInstance(ctx context.Context, h http.Header, opt *metadata.BatchCreateQuotedInstOption) (
		[]uint64, errors.CCErrorCoder)
	ListQuotedInstance(ctx context.Context, h http.Header, opt *metadata.ListQuotedInstOption) (*metadata.InstDataInfo,
		errors.CCErrorCoder)
	BatchUpdateQuotedInstance(ctx context.Context, h http.Header,
		opt *metadata.BatchUpdateQuotedInstOption) errors.CCErrorCoder
	BatchDeleteQuotedInstance(ctx context.Context, h http.Header,
		opt *metadata.BatchDeleteQuotedInstOption) errors.CCErrorCoder

	GetObjectAttrWithTable(ctx context.Context, h http.Header, params mapstr.MapStr) ([]metadata.Attribute, error)
}

// New model quote api client.
func New(client rest.ClientInterface) Interface {
	return &quote{client: client}
}

type quote struct {
	client rest.ClientInterface
}
