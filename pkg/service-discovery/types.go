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

package sd

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// ServiceInstance is the service instance used for service discovery.
type ServiceInstance struct {
	// Address is the service address.
	Address string `json:"address"`
	// UUID is used to distinguish which service is master.
	UUID string `json:"uuid"`
	// Cluster is the server's cluster, servers can only discover other servers in the same cluster.
	Cluster string `json:"cluster"`
	// Version is the server's version, servers can only discover other servers with greater or equal version.
	Version string `json:"version"`
}

// Validate service instance.
func (s *ServiceInstance) Validate() error {
	if len(s.Address) == 0 {
		return errors.New("service address cannot be empty")
	}

	if len(s.UUID) == 0 {
		s.UUID = uuid.New().String()
	}

	return nil
}

// RegisterOption is the option for registering service.
type RegisterOption struct {
	// TTL is the ttl of the registered service instance.
	TTL int64
}

// DiscoverOption is the option for discovering service. Reserved for future use, etc. service filter.
type DiscoverOption struct{}

// EventType is the service instance change event type.
type EventType string

const (
	// UpsertEvent is the create and update event type.
	UpsertEvent EventType = "upsert"
	// DeleteEvent is the delete event type.
	DeleteEvent EventType = "delete"
)

// Event is the service instance change event.
type Event struct {
	Type     EventType
	Instance *ServiceInstance
}

const (
	// servicePath is the service register path prefix.
	servicePath = "/cc/services/endpoints"
)

// GenServiceDiscoveryPath generate service register path prefix by service name.
func GenServiceDiscoveryPath(name config.ServiceName) string {
	return fmt.Sprintf("%s/%s", servicePath, name)
}

// GetServiceRegisterPath generate service register path by service name and uuid.
func GetServiceRegisterPath(name config.ServiceName, uuid string) string {
	return fmt.Sprintf("%s/%s", GenServiceDiscoveryPath(name), uuid)
}

// GenGrpcServiceDiscoveryPath generate grpc service discovery path by service name.
func GenGrpcServiceDiscoveryPath(name config.ServiceName) string {
	return fmt.Sprintf("etcd:///%s", name)
}
