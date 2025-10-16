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

package etcd_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	etcdcli "github.com/TencentBlueKing/bk-cmdb/pkg/client/etcd"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
	"github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery/etcd"
)

func initServiceDiscovery(t *testing.T) (sd.ServiceDiscovery, sd.ServiceDiscovery, sd.ServiceDiscovery) {
	ctx := context.Background()

	// create a new etcd client
	// NOTE: etcdConf is the etcd config for test, change it to your own etcd config.
	etcdConf := &etcdcli.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "",
		Password:  "",
		TLS:       nil,
	}

	etcdCli, err := etcdcli.New(etcdConf)
	if err != nil {
		t.Fatal(err)
	}

	// create etcd service discovery instances for testing
	registryOpt := &etcd.RegistryOption{
		Service: &config.ServerInfo{
			Name:        config.ApiServer,
			IP:          "127.0.0.1",
			Port:        11111,
			Environment: "a",
			UUID:        "1",
		},
	}
	discoveryOpt := &etcd.DiscoveryOption{
		Environment: "a",
		Services:    []config.ServiceName{config.ApiServer},
	}
	sd1, err := etcd.NewServiceDiscovery(ctx, etcdCli, registryOpt, discoveryOpt)
	if err != nil {
		t.Fatal(err)
	}

	registryOpt.Service.Port = 22222
	registryOpt.Service.UUID = "2"
	sd2, err := etcd.NewServiceDiscovery(ctx, etcdCli, registryOpt, discoveryOpt)
	if err != nil {
		t.Fatal(err)
	}

	registryOpt.Service.Port = 33333
	registryOpt.Service.Environment = "b"
	registryOpt.Service.UUID = "3"
	discoveryOpt.Environment = "b"
	sd3, err := etcd.NewServiceDiscovery(ctx, etcdCli, registryOpt, discoveryOpt)
	if err != nil {
		t.Fatal(err)
	}

	return sd1, sd2, sd3
}

func TestEtcdServiceDiscovery(t *testing.T) {
	ctx := context.Background()
	sd1, sd2, sd3 := initServiceDiscovery(t)

	// test watch
	watchCh, err := sd1.Watch(ctx, config.ApiServer)
	if err != nil {
		t.Fatal(err)
	}

	// test register and discover for sd1
	if err = sd1.Register(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	testWatchEvent(t, watchCh, sd.UpsertEvent, "http://127.0.0.1:11111")
	testDiscover(t, ctx, sd1, []string{"http://127.0.0.1:11111"})
	testDiscover(t, ctx, sd2, []string{"http://127.0.0.1:11111"})
	testDiscover(t, ctx, sd3, make([]string, 0))

	// test register and discover for sd2
	if err = sd2.Register(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	testWatchEvent(t, watchCh, sd.UpsertEvent, "http://127.0.0.1:22222")
	testDiscover(t, ctx, sd1, []string{"http://127.0.0.1:11111", "http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd2, []string{"http://127.0.0.1:11111", "http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd3, make([]string, 0))

	// test register and discover for sd3
	if err = sd3.Register(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	testDiscover(t, ctx, sd1, []string{"http://127.0.0.1:11111", "http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd2, []string{"http://127.0.0.1:11111", "http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd3, []string{"http://127.0.0.1:33333"})

	// test is master
	assert.True(t, sd1.IsMaster())
	assert.False(t, sd2.IsMaster())
	assert.False(t, sd3.IsMaster())

	// test deregister sd1
	if err = sd1.Deregister(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	testWatchEvent(t, watchCh, sd.DeleteEvent, "http://127.0.0.1:11111")
	testDiscover(t, ctx, sd1, []string{"http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd2, []string{"http://127.0.0.1:22222"})
	testDiscover(t, ctx, sd3, []string{"http://127.0.0.1:33333"})
	assert.False(t, sd1.IsMaster())
	assert.True(t, sd2.IsMaster())
	assert.False(t, sd3.IsMaster())

	// test register sd1 again
	if err = sd1.Register(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	testWatchEvent(t, watchCh, sd.UpsertEvent, "http://127.0.0.1:11111")
	testDiscover(t, ctx, sd1, []string{"http://127.0.0.1:22222", "http://127.0.0.1:11111"})
	testDiscover(t, ctx, sd2, []string{"http://127.0.0.1:22222", "http://127.0.0.1:11111"})
	testDiscover(t, ctx, sd3, []string{"http://127.0.0.1:33333"})
	assert.False(t, sd1.IsMaster())
	assert.True(t, sd2.IsMaster())
	assert.False(t, sd3.IsMaster())
}

func testDiscover(t *testing.T, ctx context.Context, d sd.Discovery, addresses []string) {
	instances, err := d.Discover(ctx, config.ApiServer)
	if err != nil {
		t.Fatal(err)
	}

	if len(instances) != len(addresses) {
		t.Fatalf("discovered instance %+v len is not expected, expected addresses: %+v", instances, addresses)
	}

	for i, instance := range instances {
		if instance.Address != addresses[i] {
			t.Fatalf("discovered instance %+v is not expected, expected addresses: %+v", instances, addresses)
		}
	}
}

func testWatchEvent(t *testing.T, watchCh <-chan sd.Event, eventType sd.EventType, eventAddr string) {
	select {
	case event := <-watchCh:
		if event.Type != eventType || event.Instance.Address != eventAddr {
			t.Fatalf("watch event(type: %s, instance %+v) is not expected, expected type: %v, address: %v",
				event.Type, *event.Instance, eventType, eventAddr)
		}
	default:
		t.Fatal("get no watch event")
	}
}
