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

// CreateNamespace create namespace
func (k *kube) CreateNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsCreateOption) (
	*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.NsCreateResp)

	err := k.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/namespace/bk_biz_id/%d", bizID).
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
func (k *kube) UpdateNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsUpdateOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := k.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/namespace/bk_biz_id/%d", bizID).
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
func (k *kube) DeleteNamespace(ctx context.Context, header http.Header, bizID int64,
	option *types.NsDeleteOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/namespace/bk_biz_id/%d", bizID).
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
func (k *kube) ListNamespace(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.NsDataResp, errors.CCErrorCoder) {

	result := new(types.NsInstResp)

	subPath := "/findmany/namespace"
	err := k.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
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
func (k *kube) CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlCreateOption) (*metadata.RspIDs, errors.CCErrorCoder) {

	result := new(types.WlCreateResp)

	err := k.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef("/createmany/workload/%s/%d", kind, bizID).
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
func (k *kube) UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlUpdateOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := k.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef("/updatemany/workload/%s/%d", kind, bizID).
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
func (k *kube) DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
	option *types.WlDeleteOption) errors.CCErrorCoder {

	result := new(metadata.BaseResp)

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef("/deletemany/workload/%s/%d", kind, bizID).
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
func (k *kube) ListWorkload(ctx context.Context, header http.Header, input *metadata.QueryCondition,
	kind types.WorkloadType) (*types.WlDataResp, errors.CCErrorCoder) {
	result := types.WlInstResp{
		Data: types.WlDataResp{
			Kind: kind,
		},
	}

	subPath := "/findmany/workload/%s"
	err := k.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, kind).
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

// ListPod list Pod
func (k *kube) ListPod(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.PodDataResp, errors.CCErrorCoder) {

	result := new(types.PodInstResp)

	subPath := "/findmany/pod"
	err := k.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
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

// ListContainer list Container
func (k *kube) ListContainer(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.ContainerDataResp, errors.CCErrorCoder) {

	result := new(types.ContainerInstResp)

	subPath := "/findmany/container"
	err := k.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
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

// DeletePods delete pods api
func (k *kube) DeletePods(ctx context.Context, h http.Header, opt *types.DeletePodsByIDsOption) errors.CCErrorCoder {
	resp := new(metadata.Response)
	subPath := "/deletemany/pod"

	err := k.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
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

// BatchCreateNode batch create nodes
func (k *kube) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types.CreateNodesOption) (*types.CreateNodesResult, errors.CCErrorCoder) {
	ret := new(types.CreateNodesResult)
	subPath := "/createmany/kube/node/%d"

	err := k.client.Post().
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

// BatchCreatePod batch create pod.
func (k *kube) BatchCreatePod(ctx context.Context, header http.Header,
	data *types.CreatePodsOption) ([]types.Pod, errors.CCErrorCoder) {
	ret := new(types.CreatePodsResult)
	subPath := "/createmany/kube/pod"

	err := k.client.Post().
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

	return ret.Info, nil

}

// SearchNsClusterRelation search namespace cluster relation in the shared cluster scenario.
func (k *kube) SearchNsClusterRelation(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.ResponseNsClusterRelation, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data types.ResponseNsClusterRelation `json:"data"`
	}{}

	subPath := "/findmany/kube/ns/relation"
	err := k.client.Post().
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

	return &ret.Data, nil
}

// SearchCluster create cluster.
func (k *kube) SearchCluster(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.ResponseCluster, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data types.ResponseCluster `json:"data"`
	}{}

	subPath := "/findmany/kube/cluster"
	err := k.client.Post().
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

	return &ret.Data, nil
}

// SearchNode search node.
func (k *kube) SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.SearchNodeRsp, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data types.SearchNodeRsp `json:"data"`
	}{}

	subPath := "/findmany/kube/node"
	err := k.client.Post().
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

	return &ret.Data, nil
}

// UpdateNodeFields update node fields.
func (k *kube) UpdateNodeFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateNodeOption) errors.CCErrorCoder {

	ret := new(metadata.BaseResp)
	subPath := "/updatemany/kube/node/%d"
	err := k.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	return nil
}

// UpdateClusterFields update cluster fields.
func (k *kube) UpdateClusterFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateClusterOption) errors.CCErrorCoder {

	ret := new(metadata.BaseResp)
	subPath := "/updatemany/kube/cluster/%d"
	err := k.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	return nil
}

// CreateCluster create cluster.
func (k *kube) CreateCluster(ctx context.Context, header http.Header, bizID int64,
	data *types.Cluster) (*types.CreateClusterResult, errors.CCErrorCoder) {
	ret := new(types.CreateClusterResult)
	subPath := "/create/kube/cluster/%d"
	err := k.client.Post().
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

// DeleteCluster delete cluster.
func (k *kube) DeleteCluster(ctx context.Context, header http.Header, bizID int64,
	option *types.DeleteClusterOption) errors.CCErrorCoder {
	ret := new(metadata.Response)
	subPath := "/deletemany/kube/cluster/%d"

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	return nil
}

// BatchDeleteNode delete cluster.
func (k *kube) BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
	option *types.BatchDeleteNodeOption) errors.CCErrorCoder {
	ret := new(metadata.Response)
	subPath := "/deletemany/kube/node/%d"

	err := k.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}

	return nil
}
