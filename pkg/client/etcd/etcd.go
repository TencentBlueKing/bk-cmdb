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

// Package etcd defines etcd client related types and operations.
package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// New creates a new etcd client.
func New(conf *Config) (*clientv3.Client, error) {
	ctx := context.Background()
	if err := conf.Validate(); err != nil {
		logger.Error(ctx, "validate etcd config failed", "conf", conf, logger.E(err))
		return nil, err
	}

	// convert cmdb etcd config to etcd client config
	etcdConf := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	}

	if conf.TLS != nil {
		tlsConf, enabled, err := conf.TLS.ToClientConf()
		if err != nil {
			logger.Error(ctx, "parse etcd tls config failed", "conf", conf.TLS, logger.E(err))
			return nil, err
		}
		if enabled {
			etcdConf.TLS = tlsConf
		}
	}

	// new etcd client
	etcdCli, err := clientv3.New(etcdConf)
	if err != nil {
		logger.Error(ctx, "new etcd client failed", "conf", conf, logger.E(err))
		return nil, err
	}

	return etcdCli, nil
}
