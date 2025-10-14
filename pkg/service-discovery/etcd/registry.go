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
	"encoding/json/v2"
	"fmt"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// registry is the etcd service registry implementation.
type registry struct {
	// cli is the etcd client.
	cli *clientv3.Client

	// service is the service instance to be registered.
	service *sd.ServiceInstance
	// serviceName is the service name.
	serviceName config.ServiceName

	// cancel is the cancel function of the keep alive context.
	cancel context.CancelFunc

	// masterKey is the etcd key of the mater service instance.
	masterKey atomic.Value
}

// NewRegistry creates a new service registry instance.
func NewRegistry(cli *clientv3.Client, opt *RegistryOption) (sd.Registry, error) {
	return newRegistry(cli, opt, false)
}

// NewService creates a new service registry with state instance.
func NewService(cli *clientv3.Client, opt *RegistryOption) (sd.Service, error) {
	return newRegistry(cli, opt, true)
}

func newRegistry(cli *clientv3.Client, opt *RegistryOption, withState bool) (*registry, error) {
	ctx := context.Background()

	if cli == nil || opt == nil {
		logger.Error(ctx, "new registry but etcd client or registry options is not set")
		return nil, fmt.Errorf("etcd client or registry options is not set")
	}

	if err := opt.Validate(); err != nil {
		logger.Error(ctx, "validate registry options failed", "err", err, "opt", opt)
		return nil, err
	}

	service := &sd.ServiceInstance{
		Address:     opt.Service.RegisterAddress(),
		UUID:        opt.Service.UUID,
		Environment: opt.Service.Environment,
	}
	if err := service.Validate(); err != nil {
		logger.Error(ctx, "validate registered service instance failed", "err", err, "service", service)
		return nil, err
	}

	r := &registry{
		cli:         cli,
		service:     service,
		serviceName: opt.Service.Name,
	}

	// start run service state sync logics.
	if withState {
		if err := r.runServiceStateSync(ctx); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Register service instance to registry center.
func (r *registry) Register(ctx context.Context, opts ...sd.RegisterOption) error {
	// generate service register path and value.
	key := sd.GetServiceRegisterPath(r.serviceName, r.service.UUID)
	serviceJson, err := json.Marshal(r.service)
	if err != nil {
		logger.Error(ctx, "marshal service instance failed", "err", err, "service", r.service)
		return err
	}
	value := string(serviceJson)

	// stop the previous keep alive logics.
	if r.cancel != nil {
		r.cancel()
	}

	// create etcd lease and register service instance.
	lease := clientv3.NewLease(r.cli)
	ttl := defaultRegisterTTL
	if len(opts) > 0 && opts[0].TTL > 0 {
		ttl = opts[0].TTL
	}
	leaseID, err := r.registerService(ctx, lease, ttl, key, value)
	if err != nil {
		return err
	}

	// run keep alive logics for service instance.
	go r.runKeepAlive(lease, leaseID, ttl, key, value)

	return nil
}

// registerService create etcd lease with ttl and register service instance with it.
func (r *registry) registerService(ctx context.Context, lease clientv3.Lease, ttl int64, key string, value string) (
	clientv3.LeaseID, error) {

	grantRes, err := lease.Grant(ctx, ttl)
	if err != nil {
		logger.Error(ctx, "grant lease failed", "ttl", ttl, "err", err)
		return 0, err
	}
	leaseID := grantRes.ID

	_, err = r.cli.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		logger.Error(ctx, "put etcd with lease failed", "key", key, "value", value, "err", err)
		return 0, err
	}

	return leaseID, nil
}

// runKeepAlive run keep alive for service related lease.
func (r *registry) runKeepAlive(lease clientv3.Lease, leaseID clientv3.LeaseID, ttl int64, key, value string) {
	defer func(lease clientv3.Lease) {
		if err := lease.Close(); err != nil {
			logger.Error(context.Background(), "close lease failed", "err", err)
		}
	}(lease)

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// call etcd keep alive function with the lease id.
	channel, err := r.cli.KeepAlive(ctx, leaseID)
	if err != nil {
		logger.Error(ctx, "keep alive failed", "key", key, "err", err, "lease id", leaseID)
		leaseID = 0
	}

	for {
		select {
		// context is done, cancel keep alive
		case <-ctx.Done():
			logger.Info(ctx, "cancel keep alive", "key", key)
			return
		default:
			// create a new lease if no lease is available and register service with it, keep it alive
			if leaseID == 0 {
				leaseID, err = r.registerService(ctx, lease, ttl, key, value)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				channel, err = r.cli.KeepAlive(ctx, leaseID)
				if err != nil {
					logger.Error(ctx, "keep alive failed", "key", key, "err", err, "lease id", leaseID)
					leaseID = 0
					time.Sleep(time.Second)
					continue
				}
			}

			select {
			case <-ctx.Done():
				logger.Info(ctx, "cancel keep alive", "key", key)
				return
			case resp, ok := <-channel:
				if !ok {
					logger.Info(ctx, "keep alive channel closed, register again later", "key", key, "lease id", leaseID)
					leaseID = 0
					time.Sleep(time.Second)
					continue
				}

				logger.Trace(ctx, "received keep alive response", "key", key, "lease", leaseID, "ttl", resp.TTL)
			}
		}
	}
}

// Deregister service instance from registry center.
func (r *registry) Deregister(ctx context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}

	if _, err := r.cli.Delete(ctx, sd.GetServiceRegisterPath(r.serviceName, r.service.UUID)); err != nil {
		logger.Error(ctx, "deregister service failed", "service", r.serviceName, "uuid", r.service.UUID, "err", err)
		return err
	}
	return nil
}
