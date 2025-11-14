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

	"google.golang.org/grpc/resolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// Build creates a new resolver for the given target.
func (d *discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {

	service := config.ServiceName(target.Endpoint())
	closeChan := make(chan struct{})

	// creates a new resolver with the grpc connection and stores the resolver
	r := &grpcResolver{
		cc:        cc,
		closeChan: closeChan,
	}
	r.updateServices(d.services.getServices(service))
	d.resolvers.Store(service, r)

	// removes the resolver when it's closed
	go func() {
		<-closeChan
		d.resolvers.Delete(service)
	}()
	return r, nil
}

// Scheme returns the scheme supported by this resolver.
func (d *discovery) Scheme() string {
	return "etcd"
}

// grpcResolver is the grpc resolver implementation.
type grpcResolver struct {
	// cc is the client connection that is used to update addresses.
	cc resolver.ClientConn
	// closeChan is used to notify the discovery instance that this resolver is closed.
	closeChan chan struct{}
}

// ResolveNow will be called by gRPC to try to resolve the target name again.
// It's just a hint, resolver can ignore this if it's not necessary. It could be called multiple times concurrently.
func (r *grpcResolver) ResolveNow(resolver.ResolveNowOptions) {}

// Close closes the resolver.
func (r *grpcResolver) Close() {
	close(r.closeChan)
}

// updateServices updates the service info of specified service name.
func (r *grpcResolver) updateServices(services []sd.ServiceInstance) {
	endpoints := make([]resolver.Endpoint, 0, len(services))
	for _, s := range services {
		endpoints = append(endpoints, resolver.Endpoint{
			Addresses: []resolver.Address{{Addr: s.Address}},
		})
	}

	state := resolver.State{
		Endpoints: endpoints,
	}
	if err := r.cc.UpdateState(state); err != nil {
		log.Error(context.Background(), "update grpc resolver state failed", log.E(err))
	}
}
