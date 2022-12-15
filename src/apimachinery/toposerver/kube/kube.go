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

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// BatchCreateNode batch create node
func (st *Kube) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types.CreateNodesOption) ([]int64, errors.CCErrorCoder) {
	ret := new(types.CreateNodesRsp)
	subPath := "/createmany/kube/node/bk_biz_id/%d"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data.IDs, nil
}

// BatchCreatePod batch create pod.
func (st *Kube) BatchCreatePod(ctx context.Context, header http.Header,
	data *types.CreatePodsOption) ([]int64, errors.CCErrorCoder) {
	ret := new(types.CreatePodsRsp)
	subPath := "/createmany/kube/pod"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("batch create node failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data.IDs, nil
}

// SearchCluster search cluster.
func (st *Kube) SearchCluster(ctx context.Context, header http.Header, bizID int64, input *types.QueryClusterOption) (
	*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)

	subPath := "/findmany/kube/cluster/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// SearchNode search node.
func (st *Kube) SearchNode(ctx context.Context, header http.Header, bizID int64, input *types.QueryNodeOption) (
	*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)

	subPath := "/findmany/kube/node/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// UpdateNodeFields update node fields.
func (st *Kube) UpdateNodeFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateNodeOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/updatemany/kube/node/bk_biz_id/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// UpdateClusterFields update cluster fields.
func (st *Kube) UpdateClusterFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateClusterOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/updatemany/kube/cluster/bk_biz_id/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// CreateCluster create cluster.
func (st *Kube) CreateCluster(ctx context.Context, header http.Header, bizID int64,
	data *types.Cluster) (int64, errors.CCErrorCoder) {
	ret := new(types.CreateClusterRsp)
	subPath := "/create/kube/cluster/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return 0, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return 0, ret.CCError()
	}

	return ret.Data.ID, nil
}

// DeleteCluster delete cluster.
func (st *Kube) DeleteCluster(ctx context.Context, header http.Header, bizID int64,
	option *types.DeleteClusterOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/delete/kube/cluster/bk_biz_id/%d"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}

// BatchDeleteNode delete node.
func (st *Kube) BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
	option *types.BatchDeleteNodeOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/deletemany/kube/node/bk_biz_id/%d"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}

// CreateNamespace create namespace
func (st *Kube) CreateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsCreateOption) (*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.NsCreateResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/kube/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// UpdateNamespace update namespace
func (st *Kube) UpdateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsUpdateOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := st.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/kube/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// DeleteNamespace delete namespace
func (st *Kube) DeleteNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsDeleteOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/kube/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// ListNamespace list namespace
func (st *Kube) ListNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsQueryOption) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/namespace/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// CreateWorkload create workload
func (st *Kube) CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlCreateOption) (*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.WlCreateResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/kube/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// UpdateWorkload update workload
func (st *Kube) UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlUpdateOption) errors.CCErrorCoder {
	result := new(metadata.BaseResp)

	err := st.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/kube/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// DeleteWorkload delete workload
func (st *Kube) DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlDeleteOption) errors.CCErrorCoder {
	result := new(metadata.BaseResp)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/kube/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// ListWorkload list workload
func (st *Kube) ListWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlQueryOption) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/workload/%s/%d", kind, bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// ListPod list pod
func (st *Kube) ListPod(ctx context.Context, header http.Header, bizID int64, option *types.PodQueryOption) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/pod/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// ListContainer list container
func (st *Kube) ListContainer(ctx context.Context, header http.Header, bizID int64,
	option *types.ContainerQueryOption) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/container/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// FindNodePathForHost find node path for host
func (st *Kube) FindNodePathForHost(ctx context.Context, header http.Header, option *types.HostPathOption) (
	*types.HostPathData, errors.CCErrorCoder) {

	result := new(types.HostPathResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/find/kube/host_node_path").
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}

// FindPodPath find pod path
func (st *Kube) FindPodPath(ctx context.Context, header http.Header, bizID int64, option *types.PodPathOption) (
	*types.PodPathData, errors.CCErrorCoder) {

	result := new(types.PodPathResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/find/kube/pod_path/bk_biz_id/%d", bizID).
		WithHeaders(header).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &result.Data, nil
}
