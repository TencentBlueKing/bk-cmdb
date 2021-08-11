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

package confregdiscover

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common/backbone/service_mange"
)

//DiscoverEvent if servers changed, will create a discover event
type DiscoverEvent struct { //
	Err  error
	Key  string
	Data []byte
}

// ConfRegDiscover is config register and discover
type ConfRegDiscover struct {
	client  service_mange.ClientInterface
	cancel  context.CancelFunc
	rootCtx context.Context
}

// NewConfRegDiscover used to create a object of ConfRegDiscover
func NewConfRegDiscover(cli service_mange.ClientInterface) *ConfRegDiscover {
	confRegDiscover := &ConfRegDiscover{
		client: cli,
	}
	confRegDiscover.rootCtx, confRegDiscover.cancel = context.WithCancel(context.Background())
	return confRegDiscover
}

// Ping to ping server
func (crd *ConfRegDiscover) Ping() error {
	return crd.client.Ping()
}

//Write the configure data
func (crd *ConfRegDiscover) Write(key string, data []byte) error {
	return crd.client.Put(key, string(data))
}

// Read read the configure data
func (crd *ConfRegDiscover) Read(path string) (string, error) {
	return crd.client.Get(path)
}

// Discover discover the config data
func (crd *ConfRegDiscover) Discover(key string) (<-chan *DiscoverEvent, error) {
	env := make(chan *DiscoverEvent, 1)
	go crd.loopDiscover(key, env)
	return env, nil
}

func (crd *ConfRegDiscover) loopDiscover(key string, env chan *DiscoverEvent) {
	for {
		discvEnv := &DiscoverEvent{
			Err: nil,
			Key: key,
		}

		data, err := crd.client.Get(key)
		if err != nil {
			fmt.Printf("fail to watch context for path(%s), err:%s\n", key, err.Error())
			discvEnv.Err = err
			env <- discvEnv
			time.Sleep(5 * time.Second)
			continue
		}

		discvEnv.Data = []byte(data)

		// write into discoverEvent channel
		env <- discvEnv

		time.Sleep(2 * time.Second)

		select {
		case <-crd.rootCtx.Done():
			fmt.Printf("discover path(%s) done\n", key)
			return
		default:
			fmt.Printf("watch found the content of path(%s) value\n", key)
		}
	}
}
