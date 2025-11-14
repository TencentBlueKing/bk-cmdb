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

// Package sd defines service discovery related operations.
package sd

import (
	"context"

	"google.golang.org/grpc/resolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// Registry defines service registry related operations.
type Registry interface {
	// Register service instance to registry center.
	Register(ctx context.Context, opts ...RegisterOption) error
	// Deregister service instance from registry center.
	Deregister(ctx context.Context) error
}

// Discovery defines service discovery related operations.
type Discovery interface {
	// Discover service instances from registry center.
	Discover(ctx context.Context, name config.ServiceName, opts ...DiscoverOption) ([]ServiceInstance, error)
	// Watch service instance change events.
	Watch(ctx context.Context, name config.ServiceName, opts ...DiscoverOption) (<-chan Event, error)
	// Builder is grpc resolver builder.
	resolver.Builder
}

// State defines service state related operations.
type State interface {
	// IsMaster check if current service instance is master.
	IsMaster() bool
}

// RegistryWithState defines service registry and service state related operations.
type RegistryWithState interface {
	Registry
	State
}

// ServiceDiscovery defines service registry and discovery related operations.
type ServiceDiscovery interface {
	RegistryWithState
	Discovery
}
