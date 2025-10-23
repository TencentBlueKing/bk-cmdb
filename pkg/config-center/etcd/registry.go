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
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// registry is the etcd config registry implementation.
type registry struct {
	// cli is the etcd client.
	cli *clientv3.Client
}

// NewRegistry creates a new config registry instance.
func NewRegistry(cli *clientv3.Client) (cc.Registry, error) {
	if cli == nil {
		log.Error(context.Background(), "new registry but etcd client is not set")
		return nil, fmt.Errorf("etcd client is not set")
	}

	return &registry{
		cli: cli,
	}, nil
}

// Write registers config item of specified key to the config center.
func (r *registry) Write(ctx context.Context, key string, data []byte) error {
	_, err := r.cli.Put(ctx, key, string(data))
	if err != nil {
		log.Error(ctx, "write config to etcd failed", log.E(err), "key", key, "value", string(data))
		return err
	}
	return nil
}

// Delete removes config item of specified key from the config center.
func (r *registry) Delete(ctx context.Context, key string) error {
	_, err := r.cli.Delete(ctx, key)
	if err != nil {
		log.Error(ctx, "delete config from etcd failed", "key", key, log.E(err))
		return err
	}
	return nil
}
