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

package kube

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// FindInst find instance with table name and condition
func (k *kube) FindInst(ctx context.Context, header http.Header, option *types.QueryReq) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	resp := new(metadata.QueryConditionResult)

	err := k.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/find/inst").
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &resp.Data, nil
}

// CreateNamespace create namespace
func (k *kube) CreateNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsCreateReq) (
	*types.NsCreateRespData, errors.CCErrorCoder) {

	result := types.NsCreateResp{}

	err := k.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// UpdateNamespace update namespace
func (k *kube) UpdateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsUpdateReq) errors.CCErrorCoder {

	result := metadata.BaseResp{}

	err := k.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// DeleteNamespace delete namespace
func (k *kube) DeleteNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsDeleteReq) errors.CCErrorCoder {

	result := metadata.BaseResp{}

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// CreateWorkload create workload
func (k *kube) CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlCreateReq) (*types.WlCreateRespData, errors.CCErrorCoder) {

	result := types.WlCreateResp{}

	err := k.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// UpdateWorkload update workload
func (k *kube) UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlUpdateReq) errors.CCErrorCoder {
	result := metadata.BaseResp{}

	err := k.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// DeleteWorkload delete workload
func (k *kube) DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlDeleteReq) errors.CCErrorCoder {
	result := metadata.BaseResp{}

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(&result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}
