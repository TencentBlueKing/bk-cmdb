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

package etcd

import (
	"errors"
	"time"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// DiscoveryOption is the service discovery options.
type DiscoveryOption struct {
	// Environment defines the environment to be discovered.
	Environment string
	// Services defines the services to be discovered.
	Services []config.ServiceName
}

// Validate the service discovery options.
func (o *DiscoveryOption) Validate() error {
	if len(o.Services) == 0 {
		return errors.New("services cannot be empty")
	}
	return nil
}

// RegistryOption is the service registry options.
type RegistryOption struct {
	// Service defines the service to be registered.
	Service *config.ServerInfo
}

// Validate the service registry options.
func (o *RegistryOption) Validate() error {
	if o.Service == nil {
		return errors.New("service cannot be empty")
	}

	if err := o.Service.Validate(); err != nil {
		return err
	}

	return nil
}

const (
	defaultRegisterTTL = int64(10)
	discoveryInterval  = time.Duration(defaultRegisterTTL) * time.Second
)
