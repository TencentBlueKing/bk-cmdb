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

// Package etcd defines etcd config register and discovery related operations.
package etcd

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// discovery is the etcd config discovery implementation.
type discovery struct {
	// cli is the etcd client.
	cli *clientv3.Client
}

// NewDiscovery creates a new config discovery instance.
func NewDiscovery(cli *clientv3.Client) (cc.Discovery, error) {
	if cli == nil {
		log.Error(context.Background(), "new discovery but etcd client is not set")
		return nil, fmt.Errorf("etcd client is not set")
	}

	return &discovery{
		cli: cli,
	}, nil
}

// Read reads config items of specified key from the config center.
func (d *discovery) Read(ctx context.Context, key string) ([]byte, error) {
	resp, err := d.cli.Get(ctx, key, clientv3.WithSerializable())
	if err != nil {
		log.Error(ctx, "read config from etcd failed", "key", key, log.E(err))
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		log.Info(ctx, "read no config from etcd", "key", key)
		return make([]byte, 0), nil
	}

	return resp.Kvs[0].Value, nil
}

// Watch watches config item change events of specified key from the config center.
func (d *discovery) Watch(ctx context.Context, key string) (<-chan cc.DiscoveryEvent, error) {
	watchChan := d.cli.Watch(ctx, key, clientv3.WithPrefix())
	eventChan := make(chan cc.DiscoveryEvent, 1)

	go func() {
		for resp := range watchChan {
			if resp.Err() != nil {
				log.Error(ctx, "watch config failed", "key", key, "err", resp.Err())
				close(eventChan)
				return
			}

			// parse etcd events into cmdb config discovery events
			for _, event := range resp.Events {
				var eventType cc.EventType
				var eventValue []byte

				switch event.Type {
				case mvccpb.PUT:
					eventType = cc.UpsertEvent
					eventValue = event.Kv.Value
				case mvccpb.DELETE:
					eventType = cc.DeleteEvent
					eventValue = event.PrevKv.Value
				default:
					log.Info(ctx, "event type is not supported, skip", "key", key, "event type", event.Type)
					continue
				}

				eventChan <- cc.DiscoveryEvent{
					Type: eventType,
					Key:  string(event.Kv.Key),
					Data: eventValue,
				}
			}
		}
	}()

	return eventChan, nil
}
