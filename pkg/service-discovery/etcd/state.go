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
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// IsMaster check if current service instance is master.
func (r *registry) IsMaster() bool {
	return r.masterKey.Load() == sd.GetServiceRegisterPath(r.serviceName, r.service.UUID)
}

// runServiceStateSync run service state sync logics.
func (r *registry) runServiceStateSync(ctx context.Context) error {
	path := sd.GenServiceDiscoveryPath(r.serviceName)

	// get master key from etcd
	if err := r.updateMaster(ctx, path); err != nil {
		return err
	}

	// watch service instance change events and refresh master key if current master key is deleted
	watchChan := r.cli.Watch(ctx, path, clientv3.WithPrefix())
	go func() {
		for resp := range watchChan {
			if resp.Err() != nil {
				log.Error(ctx, "watch service master key failed", "service", r.serviceName, "err", resp.Err())
				return
			}

			masterKey := r.masterKey.Load()
			if masterKey == "" {
				if err := r.updateMaster(ctx, path); err != nil {
					log.Error(ctx, "update service master key failed", "service", r.serviceName, log.E(err))
					time.Sleep(time.Second)
				}
				continue
			}

			for _, event := range resp.Events {
				if event.Type == clientv3.EventTypeDelete && string(event.Kv.Key) == masterKey {
					if err := r.updateMaster(ctx, path); err != nil {
						log.Error(ctx, "update service master key failed", "service", r.serviceName, log.E(err))
						time.Sleep(time.Second)
					}
					break
				}
			}
		}
	}()

	// loop refreshing master key in case watch encountered error
	go func() {
		time.Sleep(discoveryInterval)
		for {
			if err := r.updateMaster(ctx, path); err != nil {
				log.Error(ctx, "update service master key failed", "service", r.serviceName, log.E(err))
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(discoveryInterval)
		}
	}()

	return nil
}

// updateMaster update master key as the first created key in etcd.
func (r *registry) updateMaster(ctx context.Context, path string) error {
	// get master key from etcd
	resp, err := r.cli.Get(ctx, path, append(clientv3.WithFirstCreate(), clientv3.WithSerializable(),
		clientv3.WithKeysOnly())...)
	if err != nil {
		log.Error(ctx, "get service master key from etcd failed", "service", r.serviceName, log.E(err))
		return err
	}

	masterKey := ""
	if len(resp.Kvs) == 0 {
		log.Info(ctx, "service has no keys", "service", r.serviceName, "path", path)
	} else {
		masterKey = string(resp.Kvs[0].Key)
	}

	// update master key
	r.masterKey.Store(masterKey)

	log.Trace(ctx, "update master key", "service", r.serviceName, "key", masterKey)

	return nil
}
