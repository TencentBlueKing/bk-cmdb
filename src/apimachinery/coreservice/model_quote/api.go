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
	ListModelQuoteRelation(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
		*metadata.ListModelQuoteRelRes, errors.CCErrorCoder)
	CreateModelQuoteRelation(ctx context.Context, h http.Header, data []metadata.ModelQuoteRelation) errors.CCErrorCoder
	DeleteModelQuoteRelation(ctx context.Context, h http.Header, opt *metadata.CommonFilterOption) errors.CCErrorCoder

	BatchCreateQuotedInstance(ctx context.Context, h http.Header, objID string, data []mapstr.MapStr) (
		[]uint64, errors.CCErrorCoder)
	ListQuotedInstance(ctx context.Context, h http.Header, objID string, opt *metadata.CommonQueryOption) (
		*metadata.InstDataInfo, errors.CCErrorCoder)
	BatchUpdateQuotedInstance(ctx context.Context, h http.Header, objID string,
		opt *metadata.CommonUpdateOption) errors.CCErrorCoder
	BatchDeleteQuotedInstance(ctx context.Context, h http.Header, objID string,
		opt *metadata.CommonFilterOption) errors.CCErrorCoder
}

// New model quote api client.
func New(client rest.ClientInterface) Interface {
	return &quote{client: client}
}

type quote struct {
	client rest.ClientInterface
}
