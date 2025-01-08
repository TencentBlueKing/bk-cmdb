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
func (st *Kube) BatchCreateNode(ctx context.Context, header http.Header, data *types.CreateNodesOption) ([]int64,
	errors.CCErrorCoder) {

	ret := new(types.CreateNodesRsp)
	subPath := "/createmany/kube/node"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
func (st *Kube) BatchCreatePod(ctx context.Context, header http.Header, data *types.CreatePodsOption) ([]int64,
	errors.CCErrorCoder) {

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
func (st *Kube) SearchCluster(ctx context.Context, header http.Header, input *types.QueryClusterOption) (
	*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)

	subPath := "/findmany/kube/cluster"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
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
func (st *Kube) SearchNode(ctx context.Context, header http.Header, input *types.QueryNodeOption) (
	*metadata.Response, errors.CCErrorCoder) {

	ret := new(metadata.Response)

	subPath := "/findmany/kube/node"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
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
func (st *Kube) UpdateNodeFields(ctx context.Context, header http.Header,
	data *types.UpdateNodeOption) errors.CCErrorCoder {

	ret := new(metadata.Response)
	subPath := "/updatemany/kube/node"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}
	return nil
}

// UpdateClusterFields update cluster fields.
func (st *Kube) UpdateClusterFields(ctx context.Context, header http.Header,
	data *types.UpdateClusterOption) errors.CCErrorCoder {

	ret := new(metadata.Response)
	subPath := "/updatemany/kube/cluster"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}
	return nil
}

// CreateCluster create cluster.
func (st *Kube) CreateCluster(ctx context.Context, header http.Header,
	data *types.Cluster) (int64, errors.CCErrorCoder) {

	ret := new(types.CreateClusterRsp)
	subPath := "/create/kube/cluster"
	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
func (st *Kube) DeleteCluster(ctx context.Context, header http.Header,
	option *types.DeleteClusterOption) errors.CCErrorCoder {

	ret := new(metadata.Response)
	subPath := "/delete/kube/cluster"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// BatchDeleteNode delete node.
func (st *Kube) BatchDeleteNode(ctx context.Context, header http.Header,
	option *types.BatchDeleteNodeOption) errors.CCErrorCoder {

	ret := new(metadata.Response)
	subPath := "/deletemany/kube/node"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// CreateNamespace create namespace
func (st *Kube) CreateNamespace(ctx context.Context, header http.Header,
	option *types.NsCreateOption) (*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.NsCreateResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/kube/namespace").
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
func (st *Kube) UpdateNamespace(ctx context.Context, header http.Header,
	option *types.NsUpdateOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := st.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/kube/namespace").
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
func (st *Kube) DeleteNamespace(ctx context.Context, header http.Header,
	option *types.NsDeleteOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/kube/namespace").
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
func (st *Kube) ListNamespace(ctx context.Context, header http.Header, option *types.NsQueryOption) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/namespace").
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
func (st *Kube) CreateWorkload(ctx context.Context, header http.Header, kind types.WorkloadType,
	option *types.WlCreateOption) (*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.WlCreateResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/kube/workload/%s", kind).
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
func (st *Kube) UpdateWorkload(ctx context.Context, header http.Header, kind types.WorkloadType,
	option *types.WlUpdateOption) errors.CCErrorCoder {
	result := new(metadata.BaseResp)

	err := st.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/kube/workload/%s", kind).
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
func (st *Kube) DeleteWorkload(ctx context.Context, header http.Header, kind types.WorkloadType,
	option *types.WlDeleteOption) errors.CCErrorCoder {
	result := new(metadata.BaseResp)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/kube/workload/%s", kind).
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
func (st *Kube) ListWorkload(ctx context.Context, header http.Header, kind types.WorkloadType,
	option *types.WlQueryOption) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/workload/%s", kind).
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
func (st *Kube) ListPod(ctx context.Context, header http.Header, option *types.PodQueryOption) (
	*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/pod").
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
func (st *Kube) ListContainer(ctx context.Context, header http.Header,
	option *types.ContainerQueryOption) (*metadata.InstDataInfo, errors.CCErrorCoder) {

	result := new(metadata.ResponseInstData)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/findmany/kube/container").
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
func (st *Kube) FindPodPath(ctx context.Context, header http.Header, option *types.PodPathOption) (
	*types.PodPathData, errors.CCErrorCoder) {

	result := new(types.PodPathResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/find/kube/pod_path").
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

// DeletePods delete pods
func (st *Kube) DeletePods(ctx context.Context, header http.Header,
	params *types.DeletePodsOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := st.client.Delete().
		WithContext(ctx).
		Body(params).
		SubResourcef("/deletemany/kube/pod").
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

// ListContainerByTopo list container by topo
func (st *Kube) ListContainerByTopo(ctx context.Context, header http.Header, params *types.GetContainerByTopoOption) (
	*types.ContainerInfo, errors.CCErrorCoder) {

	result := new(types.ContainerWithTopoResp)

	err := st.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef("/findmany/kube/container/by_topo").
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

// UpdateClusterType update cluster type
func (st *Kube) UpdateClusterType(ctx context.Context, header http.Header, params *types.UpdateClusterTypeOpt) error {

	result := new(metadata.BaseResp)

	err := st.client.Put().
		WithContext(ctx).
		Body(params).
		SubResourcef("/update/kube/cluster/type").
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
