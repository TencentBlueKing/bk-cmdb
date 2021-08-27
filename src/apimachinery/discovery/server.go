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

package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

func newServerDiscover(disc *registerdiscover.RegDiscv, path, name string) (*server, error) {
	discoverChan, eventErr := disc.Watch(context.Background(), path)
	if nil != eventErr {
		return nil, eventErr
	}

	svr := &server{
		path:         path,
		name:         name,
		servers:      make(map[string]*types.ServerInfo, 0),
		discoverChan: discoverChan,
		serversChan:  make(chan []string, 1),
		rd:           disc,
	}

	svr.run()
	return svr, nil
}

type server struct {
	sync.RWMutex
	next int
	// server's name
	name         string
	path         string
	servers      map[string]*types.ServerInfo
	discoverChan <-chan *registerdiscover.DiscoverEvent
	serversChan  chan []string
	rd           *registerdiscover.RegDiscv
}

func (s *server) GetServers() ([]string, error) {
	if s == nil {
		return []string{}, nil
	}

	s.Lock()
	num := len(s.servers)
	if num == 0 {
		s.Unlock()
		return []string{}, fmt.Errorf("oops, there is no %s can be used", s.name)
	}

	var servers_array []*types.ServerInfo
	for _, server := range s.servers {
		servers_array = append(servers_array, server)
	}

	var infos []*types.ServerInfo
	if s.next < num-1 {
		s.next = s.next + 1
		infos = append(servers_array[s.next-1:], servers_array[:s.next-1]...)
	} else {
		s.next = 0
		infos = append(servers_array[num-1:], servers_array[:num-1]...)
	}
	s.Unlock()

	servers := make([]string, 0)
	for _, server := range infos {
		servers = append(servers, server.RegisterAddress())
	}

	return servers, nil
}

func (s *server) run() {
	go s.loopWatchUpdates()
	// loop compare server info in case watch encountered error
	go s.loopSyncUpdates()
}

func (s *server) loopWatchUpdates() {
	blog.Infof("start to discover cc component from register and discover, path:[%s]", s.path)
	for event := range s.discoverChan {
		blog.Infof("received one event from path: %s", s.path)
		if event.Err != nil {
			blog.Errorf("discovery received event err: %v", event.Err)
			continue
		}

		switch event.Type {
		case registerdiscover.EVENT_PUT:
			s.updateServer(event.Key, event.Value)
			s.setServersChan()
		case registerdiscover.EVENT_DEL:
			s.removeServer(event.Key, event.Value)
			s.setServersChan()
		default:
			blog.Errorf("discovery received unknown event type: %v", event.Type)
			continue
		}
	}
}

func (s *server) loopSyncUpdates() {
	var oldServer map[string]bool
	for {
		infos, err := s.rd.GetWithPrefix(s.path)
		if err != nil {
			blog.Errorf("discovery failed to sync from path: %s, err: %v", s.path, err)
			time.Sleep(3 * time.Second)
			continue
		}
		isUpdated := false
		newServer := make(map[string]bool)
		if len(infos) != len(oldServer) {
			isUpdated = true
		}
		for _, server := range infos {
			if !isUpdated && !oldServer[server.Value] {
				isUpdated = true
			}
			newServer[server.Value] = true
		}
		oldServer = newServer
		if isUpdated {
			s.resetServers(infos)
		}

		// loop sync every 3 seconds
		time.Sleep(3 * time.Second)
	}
}

// 当监听到服务节点变化时，将最新的服务节点信息放入该channel里
func (s *server) setServersChan() {
	// 即使没有其他服务消费该channel，也能保证该channel不会阻塞
	for len(s.serversChan) >= 1 {
		<-s.serversChan
	}
	s.serversChan <- s.getInstances()
}

// 获取服务发现上最新的服务节点信息channel
func (s *server) GetServersChan() chan []string {
	return s.serversChan
}

// 获取所有注册服务节点的ip:port
func (s *server) getInstances() []string {
	addrArr := []string{}
	s.RLock()
	defer s.RUnlock()
	for _, info := range s.servers {
		addrArr = append(addrArr, info.Instance())
	}
	return addrArr
}

func (s *server) updateServer(key, data string) {
	if key == "" {
		blog.Errorf("discovery received invalid event, for key is empty")
		return
	}

	server := new(types.ServerInfo)
	if err := json.Unmarshal([]byte(data), server); err != nil {
		blog.Errorf("unmarshal server info failed, key: %s, info:%s, err: %v", key, data, err)
		return
	}
	if server.Port == 0 {
		blog.Errorf("invalid port 0, with discovery key: %s", key)
		return
	}
	if len(server.RegisterIP) == 0 {
		blog.Errorf("invalid ip with discovery key: %s", key)
		return
	}
	if server.Scheme != "https" {
		server.Scheme = "http"
	}

	s.Lock()
	s.servers[key] = server
	s.Unlock()

	blog.Infof("update component with new server instance: %v, key: %s", data, key)
}

func (s *server) removeServer(key, data string) {
	if key == "" {
		blog.Errorf("discovery received invalid event, for key is empty")
		return
	}

	s.Lock()
	delete(s.servers, key)
	s.Unlock()

	blog.Infof("remove component server instance: %v, key: %s", data, key)
}

func (s *server) resetServers(kvs []registerdiscover.KeyVal) {
	servers := make(map[string]*types.ServerInfo, 0)

	for _, kv := range kvs {
		if kv.Key == "" {
			blog.Errorf("discovery sync invalid server info for key is empty")
			continue
		}
		server := new(types.ServerInfo)
		if err := json.Unmarshal([]byte(kv.Value), server); err != nil {
			blog.Errorf("unmarshal server info failed, key: %s, info:%s, err: %v", kv.Key, kv.Value, err)
			continue
		}
		if server.Port == 0 {
			blog.Errorf("invalid port 0, with discovery key: %s", kv.Key)
			return
		}
		if len(server.RegisterIP) == 0 {
			blog.Errorf("invalid ip with discovery key: %s", kv.Key)
			return
		}
		if server.Scheme != "https" {
			server.Scheme = "http"
		}
		servers[kv.Key] = server
	}

	s.Lock()
	s.servers = servers
	s.Unlock()
}