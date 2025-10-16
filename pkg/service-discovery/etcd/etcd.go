/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package etcd defines etcd service discovery related operations.
package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"

	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// serviceDiscovery is the etcd service registry and discovery implementation.
type serviceDiscovery struct {
	*registry
	*discovery
}

// NewServiceDiscovery creates a new etcd service registry and discovery instance.
func NewServiceDiscovery(ctx context.Context, cli *clientv3.Client, regOpt *RegistryOption, disOpt *DiscoveryOption) (
	sd.ServiceDiscovery, error) {

	reg, err := newRegistry(ctx, cli, regOpt, true)
	if err != nil {
		return nil, err
	}

	dis, err := newDiscovery(ctx, cli, disOpt)
	if err != nil {
		return nil, err
	}

	return &serviceDiscovery{
		registry:  reg,
		discovery: dis,
	}, nil
}
