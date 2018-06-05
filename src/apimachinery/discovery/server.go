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

	regd "configcenter/src/common/RegisterDiscover"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/framework/core/errors"
)

func newServerDiscover(disc *regd.RegDiscover, path string) (Interface, error) {
	discoverChan, eventErr := disc.DiscoverService(path)
	if nil != eventErr {
		return nil, eventErr
	}

	svr := &server{
		path:         path,
		servers:      make([]string, 0),
		discoverChan: discoverChan,
	}

	svr.run()
	return svr, nil
}

type server struct {
	sync.Mutex
	index        int
	path         string
	servers      []string
	discoverChan <-chan *regd.DiscoverEvent
}

func (s *server) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
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

		scheme := "http"
		if server.Scheme == "https" {
			scheme = "https"
		}

		if server.Port == 0 {
			blog.Errorf("invalid port 0, with zk path: %s", s.path)
			continue
		}

		if len(server.IP) == 0 {
			blog.Errorf("invalid ip with zk path: %s", s.path)
			continue
		}

		host := fmt.Sprintf("%s://%s:%d", scheme, server.IP, server.Port)
		newSvr = append(newSvr, host)
	}
	
	s.Lock()
	defer s.Unlock()

	if len(newSvr) != 0 {
		s.servers = newSvr
		blog.V(3).Infof("update component with new server instance[%s] about path: %s", strings.Join(newSvr, "; "), s.path)
	}
}
