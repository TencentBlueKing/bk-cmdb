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

package modelquote

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// BatchCreateQuotedInstance batch create quoted instances
func (q quote) BatchCreateQuotedInstance(ctx context.Context, h http.Header, objID string, data []mapstr.MapStr) (
	[]uint64, errors.CCErrorCoder) {

	resp := new(metadata.BatchCreateResp)

	err := q.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef("/createmany/quoted/model/%s/instance", objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data.IDs, nil
}

// ListQuotedInstance list quoted instances
func (q quote) ListQuotedInstance(ctx context.Context, h http.Header, objID string, opt *metadata.CommonQueryOption) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	resp := new(metadata.ResponseInstData)

	err := q.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/quoted/model/%s/instance", objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// BatchUpdateQuotedInstance batch update quoted instance
func (q quote) BatchUpdateQuotedInstance(ctx context.Context, h http.Header, objID string,
	opt *metadata.CommonUpdateOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := q.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/updatemany/quoted/model/%s/instance", objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// BatchDeleteQuotedInstance batch delete quoted instances
func (q quote) BatchDeleteQuotedInstance(ctx context.Context, h http.Header, objID string,
	opt *metadata.CommonFilterOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := q.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/deletemany/quoted/model/%s/instance", objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}
