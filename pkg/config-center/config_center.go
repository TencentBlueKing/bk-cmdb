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

// Package cc is cmdb's config center.
package cc

import "context"

// ConfigCenter defines config related operations.
type ConfigCenter interface {
	// Get gets config item of specified key from config center.
	Get(conf ConfigType, key string) (any, bool)
}

// Registry defines config registry related operations.
type Registry interface {
	// Write registers config item of specified key to the config center.
	Write(ctx context.Context, key string, data []byte) error
	// Delete removes config item of specified key from the config center.
	Delete(ctx context.Context, key string) error
}

// Discovery defines config discovery related operations.
type Discovery interface {
	// Read reads config items of specified key from the config center.
	Read(ctx context.Context, key string) ([]byte, error)
	// Watch watches config item change events of specified key from the config center.
	Watch(ctx context.Context, key string) (<-chan DiscoveryEvent, error)
}
