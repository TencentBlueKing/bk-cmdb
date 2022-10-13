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

// KubeOperationInterface the kube implements the interface
type KubeOperationInterface interface {
	CreateCluster(ctx context.Context, h http.Header, bizID int64, option *types.Cluster) (
		*metadata.Response, errors.CCErrorCoder)
	UpdateClusterFields(ctx context.Context, header http.Header, bizID int64,
		data *types.UpdateClusterOption) (*metadata.Response, errors.CCErrorCoder)
	SearchCluster(ctx context.Context, header http.Header, bizID int64, input *types.QueryClusterOption) (
		*metadata.Response, errors.CCErrorCoder)
	DeleteCluster(ctx context.Context, header http.Header, bizID int64,
		option *types.DeleteClusterOption) (*metadata.Response, errors.CCErrorCoder)
	BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
		option *types.BatchDeleteNodeOption) (*metadata.Response, errors.CCErrorCoder)
	BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
		data *types.CreateNodesOption) (*metadata.Response, errors.CCErrorCoder)
	UpdateNodeFields(ctx context.Context, header http.Header, bizID int64, data *types.UpdateNodeOption) (
		*metadata.Response, errors.CCErrorCoder)
	SearchNode(ctx context.Context, header http.Header, bizID int64, input *types.QueryNodeOption) (
		*metadata.Response, errors.CCErrorCoder)
	BatchCreatePod(ctx context.Context, header http.Header, data *types.CreatePodsOption) (
		[]types.Pod, errors.CCErrorCoder)
}

// NewKubeOperationInterface initialize the container client object
func NewKubeOperationInterface(client rest.ClientInterface) KubeOperationInterface {
	return &Kube{client: client}
}

// Kube container object
type Kube struct {
	client rest.ClientInterface
}
