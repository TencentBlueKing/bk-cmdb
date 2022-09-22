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

package container

import (
	"context"
	"net/http"

	types2 "configcenter/pkg/kube/types"

	"configcenter/api/rest"
	"configcenter/pkg/errors"
	"configcenter/pkg/metadata"
)

// ContainerInterface the container implements the interface
type ContainerInterface interface {
	CreateCluster(ctx context.Context, h http.Header, bizID int64, option *types2.ClusterBaseFields) (
		*types2.CreateClusterResult, errors.CCErrorCoder)
	UpdateClusterFields(ctx context.Context, header http.Header, supplierAccount string, bizID int64,
		data *types2.UpdateClusterOption) errors.CCErrorCoder
	UpdateNodeFields(ctx context.Context, header http.Header, supplierAccount string, bizID int64,
		data *types2.UpdateNodeOption) errors.CCErrorCoder
	SearchCluster(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
		*types2.ResponseCluster, errors.CCErrorCoder)
	DeleteCluster(ctx context.Context, header http.Header, bizID int64,
		option *types2.DeleteClusterOption) errors.CCErrorCoder
	BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
		option *types2.BatchDeleteNodeOption) errors.CCErrorCoder
	BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
		data *types2.CreateNodesOption) (*types2.CreateNodesResult, errors.CCErrorCoder)
	BatchCreatePod(ctx context.Context, header http.Header, bizID int64,
		data *types2.CreatePodsOption) ([]types2.Pod, errors.CCErrorCoder)
	SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
		*types2.SearchNodeRsp, errors.CCErrorCoder)
}

// NewContainerInterface initialize the container client object
func NewContainerInterface(client rest.ClientInterface) ContainerInterface {
	return &Container{client: client}
}

// Container container object
type Container struct {
	client rest.ClientInterface
}
