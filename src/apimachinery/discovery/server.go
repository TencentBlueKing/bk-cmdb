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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

func newServerDiscover(disc *registerdiscover.RegDiscover, path, name string) (*server, error) {
	discoverChan, eventErr := disc.DiscoverService(path)
	if nil != eventErr {
		return nil, eventErr
	}

	svr := &server{
		path:         path,
		name:         name,
		servers:      make([]string, 0),
		discoverChan: discoverChan,
	}

	svr.run()
	return svr, nil
}

type server struct {
	sync.RWMutex
	index int
	// server's name
	name         string
	path         string
	servers      []string
	discoverChan <-chan *registerdiscover.DiscoverEvent
}

func (s *server) GetServers() ([]string, error) {
	if s == nil {
		return []string{}, nil
	}
	s.RLock()
	defer s.RUnlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, fmt.Errorf("oops, there is no %s can be used", s.name)
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}

// IsMaster 判断当前进程是否为master 进程， 服务注册节点的第一个节点
func (s *server) IsMaster(strAddrs string) bool {
	if s == nil {
		return false
	}
	s.RLock()
	defer s.RUnlock()
	if 0 < len(s.servers) {
		return s.servers[0] == strAddrs
	}
	return false

}

func (s *server) run() {
	blog.Infof("start to discover cc component from zk, path:[%s].", s.path)
	go func() {
		for svr := range s.discoverChan {
			blog.Warnf("received one zk event from path %s.", s.path)
			if svr.Err != nil {
				blog.Errorf("get zk event with error about path[%s]. err: %v", s.path, svr.Err)
				continue
			}

			if len(svr.Server) <= 0 {
				blog.Warnf("get zk event with 0 instance with path[%s], reset its servers", s.path)
				s.resetServer()
				continue
			}

			s.updateServer(svr.Server)
		}
	}()
}

func (s *server) resetServer() {
	s.Lock()
	defer s.Unlock()
	s.servers = make([]string, 0)
}

func (s *server) updateServer(svrs []string) {
	newSvr := make([]string, 0)

	for _, svr := range svrs {
		server := new(types.ServerInfo)
		if err := json.Unmarshal([]byte(svr), server); err != nil {
			blog.Errorf("unmarshal server info failed, zk path[%s], err: %v", s.path, err)
			continue
		}

		if server.Scheme != "https" {
			server.Scheme = "http"
		}

		if server.Port == 0 {
			blog.Errorf("invalid port 0, with zk path: %s", s.path)
			continue
		}

		if len(server.IP) == 0 {
			blog.Errorf("invalid ip with zk path: %s", s.path)
			continue
		}

		host := server.Address()
		newSvr = append(newSvr, host)
	}

	s.Lock()
	defer s.Unlock()

	if len(newSvr) != 0 {
		s.servers = newSvr
		blog.V(5).Infof("update component with new server instance[%s] about path: %s", strings.Join(newSvr, "; "), s.path)
	}
}
