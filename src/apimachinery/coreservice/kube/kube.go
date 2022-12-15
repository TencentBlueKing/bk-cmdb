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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// KubeClientInterface the kube client interface
type KubeClientInterface interface {

	// CreateNamespace create namespace
	CreateNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsCreateOption) (
		*metadata.RspIDs, errors.CCErrorCoder)

	// UpdateNamespace update namespace
	UpdateNamespace(ctx context.Context, header http.Header, bizID int64,
		option *types.NsUpdateOption) errors.CCErrorCoder

	// DeleteNamespace delete namespace
	DeleteNamespace(ctx context.Context, header http.Header, bizID int64,
		option *types.NsDeleteOption) errors.CCErrorCoder

	// ListNamespace list namespace
	ListNamespace(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types.NsDataResp,
		errors.CCErrorCoder)

	// CreateWorkload create workload
	CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlCreateOption) (*metadata.RspIDs, errors.CCErrorCoder)

	// UpdateWorkload update workload
	UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlUpdateOption) errors.CCErrorCoder

	// DeleteWorkload delete workload
	DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlDeleteOption) errors.CCErrorCoder

	// ListWorkload list workload
	ListWorkload(ctx context.Context, header http.Header, input *metadata.QueryCondition, kind types.WorkloadType) (
		*types.WlDataResp, errors.CCErrorCoder)

	// ListPod list pod
	ListPod(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types.PodDataResp,
		errors.CCErrorCoder)

	// ListContainer list container
	ListContainer(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types.ContainerDataResp,
		errors.CCErrorCoder)

	// DeletePods delete pods
	DeletePods(ctx context.Context, h http.Header, opt *types.DeletePodsByIDsOption) errors.CCErrorCoder

	CreateCluster(ctx context.Context, h http.Header, bizID int64, option *types.Cluster) (
		*types.CreateClusterResult, errors.CCErrorCoder)
	UpdateClusterFields(ctx context.Context, header http.Header, bizID int64,
		data *types.UpdateClusterOption) errors.CCErrorCoder
	UpdateNodeFields(ctx context.Context, header http.Header, bizID int64,
		data *types.UpdateNodeOption) errors.CCErrorCoder
	SearchCluster(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
		*types.ResponseCluster, errors.CCErrorCoder)
	DeleteCluster(ctx context.Context, header http.Header, bizID int64,
		option *types.DeleteClusterOption) errors.CCErrorCoder
	BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
		option *types.BatchDeleteNodeOption) errors.CCErrorCoder
	BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
		data *types.CreateNodesOption) (*types.CreateNodesResult, errors.CCErrorCoder)
	BatchCreatePod(ctx context.Context, header http.Header, data *types.CreatePodsOption) (
		[]types.Pod, errors.CCErrorCoder)
	SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
		*types.SearchNodeRsp, errors.CCErrorCoder)
}

// NewKubeClientInterface new kube client interface
func NewKubeClientInterface(client rest.ClientInterface) KubeClientInterface {
	return &kube{client: client}
}

type kube struct {
	client rest.ClientInterface
}
