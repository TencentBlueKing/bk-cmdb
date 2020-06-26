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

package client

import (
	"errors"
	"sync"

	"configcenter/src/apimachinery/util"
	"github.com/prometheus/client_golang/prometheus"
)

type IamConfig struct {
	// blueking's auth center addresses
	Address []string
	// app code is used for authorize used.
	AppCode string
	// app secret is used for authorized
	AppSecret string
	// the system id which used in auth center.
	SystemID string
	// http TLS config
	TLS util.TLSClientConfig
}

func (a IamConfig) Validate() error {
	if len(a.Address) == 0 {
		return errors.New("no iam address")
	}

	if len(a.AppCode) == 0 {
		return errors.New("no iam app code")
	}

	if len(a.AppSecret) == 0 {
		return errors.New("no iam app secret")
	}
	return nil
}

type Options struct {
	Metric prometheus.Registerer
}

type acDiscovery struct {
	// auth's servers address, must prefixed with http:// or https://
	servers []string
	index   int
	sync.Mutex
}

func (s *acDiscovery) GetServers() ([]string, error) {
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

func (s *acDiscovery) GetServersChan() chan []string {
	return nil
}
