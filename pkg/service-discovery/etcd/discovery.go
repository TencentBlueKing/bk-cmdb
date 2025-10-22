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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
	"github.com/TencentBlueKing/bk-cmdb/pkg/version"
)

// discovery is the etcd service discovery implementation.
type discovery struct {
	// cli is the etcd client.
	cli *clientv3.Client

	// environment defines the environment to be discovered.
	environment string
	// supportedServices is used to check if the service name is supported for discovery.
	supportedServices map[config.ServiceName]struct{}

	// services defines the discovered service info.
	services *discoveredServices
}

// NewDiscovery creates a new service discovery instance.
func NewDiscovery(ctx context.Context, cli *clientv3.Client, opt *DiscoveryOption) (sd.Discovery, error) {
	return newDiscovery(ctx, cli, opt)
}

func newDiscovery(ctx context.Context, cli *clientv3.Client, opt *DiscoveryOption) (*discovery, error) {
	if cli == nil || opt == nil {
		logger.Error(ctx, "new discovery but etcd client or discovery options is not set")
		return nil, fmt.Errorf("etcd client or discovery options is not set")
	}

	if err := opt.Validate(); err != nil {
		logger.Error(ctx, "validate discovery options failed", logger.E(err), "opt", opt)
		return nil, err
	}

	d := &discovery{
		cli:               cli,
		environment:       opt.Environment,
		supportedServices: make(map[config.ServiceName]struct{}),
		services: &discoveredServices{
			services:      make(map[config.ServiceName][]sd.ServiceInstance),
			serviceIdxMap: make(map[config.ServiceName]map[string]int),
		},
	}
	for _, service := range opt.Services {
		d.supportedServices[service] = struct{}{}
		d.services.serviceIdxMap[service] = make(map[string]int)

		if err := d.runDiscovery(ctx, service); err != nil {
			logger.Error(ctx, "run discovery failed", "service", service, logger.E(err))
			return nil, err
		}
	}

	return d, nil
}

// runDiscovery runs service discovery for specified service.
func (d *discovery) runDiscovery(ctx context.Context, name config.ServiceName) error {
	path := sd.GenServiceDiscoveryPath(name)

	// initialize service instances
	err := d.refreshService(ctx, name, path)
	if err != nil {
		return err
	}

	// watch service instance change events and sync service instance info
	if err = d.watchService(ctx, name); err != nil {
		return err
	}

	// loop discover service instances in case watch encountered error
	go func() {
		time.Sleep(discoveryInterval)
		for {
			if err = d.refreshService(ctx, name, path); err != nil {
				logger.Error(ctx, "refresh service failed", "service", name, logger.E(err))
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(discoveryInterval)
		}
	}()
	return nil
}

// refreshService refreshes service instances from etcd.
func (d *discovery) refreshService(ctx context.Context, name config.ServiceName, path string) error {
	// get service instances from etcd
	resp, err := d.cli.Get(ctx, path, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		logger.Error(ctx, "get service from etcd failed", "service", name, logger.E(err))
		return err
	}

	// parse service instances
	services := make([]sd.ServiceInstance, 0)
	serviceIdxMap := make(map[string]int)
	for _, kv := range resp.Kvs {
		var service sd.ServiceInstance
		if err = json.Unmarshal(kv.Value, &service); err != nil {
			logger.Error(ctx, "unmarshal service instance failed", "value", string(kv.Value), logger.E(err))
			return err
		}

		if service.Environment != d.environment {
			logger.Info(ctx, "skip service instance with different environment", "value", string(kv.Value))
			continue
		}

		if !version.IsGreaterOrEqualVersion(service.Version) {
			logger.Info(ctx, "skip service instance with lower version", "value", string(kv.Value))
			continue
		}

		services = append(services, service)
		serviceIdxMap[service.UUID] = len(services) - 1
	}

	// update service instances
	d.services.updateServices(name, services, serviceIdxMap)

	return nil
}

// watchService watches service instance change events and sync service instance info.
func (d *discovery) watchService(ctx context.Context, name config.ServiceName) error {
	watchChan, err := d.Watch(ctx, name)
	if err != nil {
		logger.Error(ctx, "watch service failed", "service", name, logger.E(err))
		return err
	}

	go func() {
		for event := range watchChan {
			switch event.Type {
			case sd.UpsertEvent:
				d.services.upsertService(name, event.Instance)
			case sd.DeleteEvent:
				d.services.deleteService(name, event.Instance)
			}
		}
	}()
	return nil
}

// Discover service instances from registry center.
func (d *discovery) Discover(ctx context.Context, name config.ServiceName, opts ...sd.DiscoverOption) (
	[]sd.ServiceInstance, error) {

	if _, exists := d.supportedServices[name]; !exists {
		return nil, fmt.Errorf("service %s is not supported", name)
	}

	return d.services.getServices(name), nil
}

// Watch service instance change events.
func (d *discovery) Watch(ctx context.Context, name config.ServiceName, opts ...sd.DiscoverOption) (<-chan sd.Event,
	error) {

	if _, exists := d.supportedServices[name]; !exists {
		return nil, fmt.Errorf("service %s is not supported", name)
	}

	// watch service instance change events from etcd by service discovery path prefix
	path := sd.GenServiceDiscoveryPath(name)
	watchChan := d.cli.Watch(ctx, path, clientv3.WithPrefix(), clientv3.WithPrevKV())

	eventChan := make(chan sd.Event, 1)

	go func() {
		for resp := range watchChan {
			if resp.Err() != nil {
				logger.Error(ctx, "watch service failed", "service", name, "err", resp.Err())
				close(eventChan)
				return
			}

			// parse etcd events into cmdb service discovery events
			for _, event := range resp.Events {
				var eventType sd.EventType
				var eventValue []byte

				switch event.Type {
				case mvccpb.PUT:
					eventType = sd.UpsertEvent
					eventValue = event.Kv.Value
				case mvccpb.DELETE:
					eventType = sd.DeleteEvent
					eventValue = event.PrevKv.Value
				default:
					logger.Error(ctx, "event type is not supported", "type", event.Type)
					continue
				}

				// generate cmdb service discovery event from etcd event
				service := new(sd.ServiceInstance)
				if err := json.Unmarshal(eventValue, service); err != nil {
					logger.Error(ctx, "unmarshal service instance failed", "value", string(eventValue), logger.E(err))
					continue
				}

				if service.Environment != d.environment {
					logger.Info(ctx, "skip service instance with different environment", "value", string(eventValue))
					continue
				}

				if !version.IsGreaterOrEqualVersion(service.Version) {
					logger.Info(ctx, "skip service instance with lower version", "value", string(eventValue))
					continue
				}

				eventChan <- sd.Event{
					Type:     eventType,
					Instance: service,
				}
			}
		}
	}()

	return eventChan, nil
}

// discoveredServices defines the discovered service info.
type discoveredServices struct {
	// services is the service name to service instance map.
	services map[config.ServiceName][]sd.ServiceInstance
	// serviceIdxMap is the service name to service instance uuid to index map.
	serviceIdxMap map[config.ServiceName]map[string]int
	// servicesLock is used to lock the services and serviceIdxMap maps.
	servicesLock sync.RWMutex
}

// getServices gets the service instances of specified service name.
func (s *discoveredServices) getServices(name config.ServiceName) []sd.ServiceInstance {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()

	return s.services[name]
}

// updateServices updates the service info of specified service name.
func (s *discoveredServices) updateServices(name config.ServiceName, services []sd.ServiceInstance,
	serviceIdxMap map[string]int) {

	s.servicesLock.Lock()
	s.services[name] = services
	s.serviceIdxMap[name] = serviceIdxMap
	s.servicesLock.Unlock()
}

// upsertService updates or inserts a service instance of specified service name.
func (s *discoveredServices) upsertService(name config.ServiceName, service *sd.ServiceInstance) {
	s.servicesLock.Lock()
	defer s.servicesLock.Unlock()

	// get service instance index by uuid
	idx, exists := s.serviceIdxMap[name][service.UUID]
	if exists {
		// update service instance info if exists
		s.services[name][idx] = *service
		return
	}

	// insert service instance if not exists
	s.services[name] = append(s.services[name], *service)
	s.serviceIdxMap[name][service.UUID] = len(s.services[name]) - 1
}

// deleteService deletes a service instance of specified service name.
func (s *discoveredServices) deleteService(name config.ServiceName, service *sd.ServiceInstance) {
	s.servicesLock.Lock()
	defer s.servicesLock.Unlock()

	// delete service instance if exists
	idx, exists := s.serviceIdxMap[name][service.UUID]
	if exists {
		s.services[name] = append(s.services[name][:idx], s.services[name][idx+1:]...)
		delete(s.serviceIdxMap[name], service.UUID)
	}
}
