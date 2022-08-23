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
	// FindInst find instance with table name and condition
	FindInst(ctx context.Context, header http.Header, option *types.QueryReq) (
		*metadata.InstDataInfo, errors.CCErrorCoder)

	// CreateNamespace create namespace
	CreateNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsCreateReq) (
		*types.NsCreateRespData, errors.CCErrorCoder)

	// UpdateNamespace update namespace
	UpdateNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsUpdateReq) errors.CCErrorCoder

	// DeleteNamespace delete namespace
	DeleteNamespace(ctx context.Context, header http.Header, bizID int64, option *types.NsDeleteReq) errors.CCErrorCoder

	// CreateWorkload create workload
	CreateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlCreateReq) (*types.WlCreateRespData, errors.CCErrorCoder)

	// UpdateWorkload update workload
	UpdateWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlUpdateReq) errors.CCErrorCoder

	// DeleteWorkload delete workload
	DeleteWorkload(ctx context.Context, header http.Header, bizID int64, kind types.WorkloadType,
		option *types.WlDeleteReq) errors.CCErrorCoder
}

// NewKubeClientInterface new kube client interface
func NewKubeClientInterface(client rest.ClientInterface) KubeClientInterface {
	return &kube{client: client}
}

type kube struct {
	client rest.ClientInterface
}
