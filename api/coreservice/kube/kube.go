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

	types2 "configcenter/pkg/kube/types"

	"configcenter/api/rest"
	"configcenter/pkg/errors"
	"configcenter/pkg/metadata"
)

// KubeClientInterface the kube client interface
type KubeClientInterface interface {
	// FindInst find instance with table name and condition
	FindInst(ctx context.Context, header http.Header, option *types2.QueryReq) (*metadata.InstDataInfo,
		errors.CCErrorCoder)

	// CreateNamespace create namespace
	CreateNamespace(ctx context.Context, header http.Header, bizID int64, option *types2.NsCreateReq) (
		*types2.NsCreateRespData, errors.CCErrorCoder)

	// UpdateNamespace update namespace
	UpdateNamespace(ctx context.Context, header http.Header, bizID int64, option *types2.NsUpdateReq) errors.CCErrorCoder

	// DeleteNamespace delete namespace
	DeleteNamespace(ctx context.Context, header http.Header, bizID int64, option *types2.NsDeleteReq) errors.CCErrorCoder

	// ListNamespace list namespace
	ListNamespace(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types2.NsDataResp,
		errors.CCErrorCoder)

	// CreateWorkload create workload
	CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types2.WorkloadType,
		option *types2.WlCreateReq) (*types2.WlCreateRespData, errors.CCErrorCoder)

	// UpdateWorkload update workload
	UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types2.WorkloadType,
		option *types2.WlUpdateReq) errors.CCErrorCoder

	// DeleteWorkload delete workload
	DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types2.WorkloadType,
		option *types2.WlDeleteReq) errors.CCErrorCoder

	// ListWorkload list workload
	ListWorkload(ctx context.Context, header http.Header, input *metadata.QueryCondition, kind types2.WorkloadType) (
		*types2.WlDataResp, errors.CCErrorCoder)

	// ListPod list pod
	ListPod(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types2.PodDataResp,
		errors.CCErrorCoder)

	// ListContainer list container
	ListContainer(ctx context.Context, header http.Header, input *metadata.QueryCondition) (*types2.ContainerDataResp,
		errors.CCErrorCoder)

	// DeletePods delete pods
	DeletePods(ctx context.Context, h http.Header, opt *types2.DeletePodsByIDsOption) errors.CCErrorCoder
}

// NewKubeClientInterface new kube client interface
func NewKubeClientInterface(client rest.ClientInterface) KubeClientInterface {
	return &kube{client: client}
}

type kube struct {
	client rest.ClientInterface
}
