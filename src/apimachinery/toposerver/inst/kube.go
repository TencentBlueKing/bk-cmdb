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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// CreateNamespace create namespace
func (t *instanceClient) CreateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsCreateReq) (*types.NsCreateRespData, errors.CCErrorCoder) {

	result := types.NsCreateResp{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/createmany/namespace/bk_biz_id/%d", bizID).
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
func (t *instanceClient) UpdateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsUpdateReq) errors.CCErrorCoder {

	result := metadata.BaseResp{}

	err := t.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/updatemany/namespace/bk_biz_id/%d", bizID).
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
func (t *instanceClient) DeleteNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsDeleteReq) errors.CCErrorCoder {

	result := metadata.BaseResp{}

	err := t.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/deletemany/namespace/bk_biz_id/%d", bizID).
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

// ListNamespace list namespace
func (t *instanceClient) ListNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsQueryReq) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := metadata.ResponseInstData{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/findmany/namespace/bk_biz_id/%d", bizID).
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

// CreateWorkload create workload
func (t *instanceClient) CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlCreateReq) (*types.WlCreateRespData, errors.CCErrorCoder) {

	result := types.WlCreateResp{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/createmany/workload/%s/%d", kind, bizID).
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
func (t *instanceClient) UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlUpdateReq) errors.CCErrorCoder {
	result := metadata.BaseResp{}

	err := t.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/updatemany/workload/%s/%d", kind, bizID).
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
func (t *instanceClient) DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlDeleteReq) errors.CCErrorCoder {
	result := metadata.BaseResp{}

	err := t.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/deletemany/workload/%s/%d", kind, bizID).
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

// ListWorkload list workload
func (t *instanceClient) ListWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlQueryReq) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := metadata.ResponseInstData{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/findmany/workload/%s/%d", kind, bizID).
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

// ListPod list pod
func (t *instanceClient) ListPod(ctx context.Context, header http.Header, bizID int64, option *types.PodQueryReq) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := metadata.ResponseInstData{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/findmany/pod/bk_biz_id/%d", bizID).
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

// ListContainer list container
func (t *instanceClient) ListContainer(ctx context.Context, header http.Header, bizID int64,
	option *types.ContainerQueryReq) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := metadata.ResponseInstData{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/findmany/container/bk_biz_id/%d", bizID).
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

// FindNodePathForHost find node path for host
func (t *instanceClient) FindNodePathForHost(ctx context.Context, header http.Header, option *types.HostPathReq) (
	*types.HostPathData, errors.CCErrorCoder) {

	result := types.HostPathResp{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/find/host_node_path").
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

// FindPodPath find pod path
func (t *instanceClient) FindPodPath(ctx context.Context, header http.Header, bizID int64, option *types.PodPathReq) (
	*types.PodPathData, errors.CCErrorCoder) {

	result := types.PodPathResp{}

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/kube/find/pod_path/bk_biz_id/%d", bizID).
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
