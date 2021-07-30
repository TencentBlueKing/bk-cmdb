/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package registerdiscover

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common/backbone/service_mange/etcd"
	"configcenter/src/common/blog"

	"github.com/coreos/etcd/clientv3"
)

const (
	// defaultPingInterval is default ping action interval.
	defaultPingInterval = 3 * time.Second
)

// etcdRegDiscv has methods for service discovery and registration
type etcdRegDiscv struct {
	// cancel method
	cancel  context.CancelFunc
	rootCxt context.Context
	// etcd client
	etcdCli *etcd.EtcdCli
	// service keepalive ttl.
	ttl int64
	// etcd leaseid
	leaseid clientv3.LeaseID
	// etcd kv pair for service instance.
	key   string
	value string
	// is it clear to register
	isClearRegister bool
}

// NewEtcdRegDiscv create a object of etcdRegDiscv
func NewEtcdRegDiscv(client *etcd.EtcdCli) RegDiscvServer {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &etcdRegDiscv{
		cancel:  ctxCancel,
		rootCxt: ctx,
		etcdCli: client,
	}
}

func (e *etcdRegDiscv) Ping() error {
	for !e.isClearRegister {
		if err := e.etcdCli.Ping(e.leaseid); err != nil {
			return err
		}
		select {
		case <-e.rootCxt.Done():
			break
		default:
			time.Sleep(defaultPingInterval)
		}
	}
	return nil
}

func (e *etcdRegDiscv) RegisterAndWatch(path string, data []byte) error {
	blog.Infof("register server. path(%s), data(%s)", path, string(data))
	e.key = path
	e.value = string(data)
	// TODO 拿到ttl配置，这里先设置默认值
	leaseID, err := e.etcdCli.PutAndBindLease(e.key, e.value, 10)
	if err != nil {
		return err
	}
	e.leaseid = leaseID
	return nil
}

func (e *etcdRegDiscv) GetServNodes(key string) ([]string, error) {
	panic("implement me")
}

func (e *etcdRegDiscv) Discover(path string) (<-chan *DiscoverEvent, error) {
	fmt.Printf("begin to discover by watch node of path(%s)\n", path)
	discvCtx := e.rootCxt
	env := make(chan *DiscoverEvent, 1)
	// loop compare server info in case watch encountered error
	go func() {
		var oldServer map[string]bool
		for {
			event := e.getServerInfoByPath(path)
			isUpdated := false
			newServer := make(map[string]bool)
			if len(event.Server) != len(oldServer) {
				isUpdated = true
			}
			for _, server := range event.Server {
				if !isUpdated && !oldServer[server] {
					isUpdated = true
				}
				newServer[server] = true
			}
			oldServer = newServer
			if isUpdated {
				env <- event
			}
			select {
			case <-discvCtx.Done():
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	return env, nil
}

func (e *etcdRegDiscv) getServerInfoByPath(path string) *DiscoverEvent {
	// get server infos
	serverInfos, err := e.etcdCli.Get(path)
	if err != nil {
		fmt.Errorf("fail to get server info from etcd by path(%s), err:%s\n", path, err)
	}
	discvEnv := &DiscoverEvent{
		Key:    path,
		Server: serverInfos,
	}
	return discvEnv
}

func (e *etcdRegDiscv) Cancel() {
	e.cancel()
}

func (e *etcdRegDiscv) ClearRegisterPath() error {
	e.isClearRegister = true
	return nil
}
