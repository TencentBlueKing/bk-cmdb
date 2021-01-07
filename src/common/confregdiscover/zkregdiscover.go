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

	"configcenter/src/common/zkclient"
	"configcenter/src/common/backbone/service_mange/zk"
)

// ZkRegDiscover config register and discover by zookeeper
type ZkRegDiscover struct {
	zkcli          *zkclient.ZkClient
	cancel         context.CancelFunc
	rootCtx        context.Context
	sessionTimeOut time.Duration
}

// NewZkRegDiscover create a object of ZkRegDiscover
func NewZkRegDiscover(client *zk.ZkClient) *ZkRegDiscover {
	ctx, ctxCancel := client.WithCancel()
	return &ZkRegDiscover{
		zkcli:          client.Client(),
		sessionTimeOut: client.SessionTimeOut(),
		rootCtx:        ctx,
		cancel:         ctxCancel,
	}
}

// Ping to ping server
func (zkRD *ZkRegDiscover) Ping() error {
	return zkRD.zkcli.Ping()
}

//Write to save config data into zookeeper
func (zkRD *ZkRegDiscover) Write(path string, data []byte) error {
	return zkRD.zkcli.Update(path, string(data))
}

func (zkRD *ZkRegDiscover) Read(path string) (string, error) {
	return zkRD.zkcli.Get(path)
}

func (zkRD *ZkRegDiscover) Discover(key string) (<-chan *DiscoverEvent, error) {

	env := make(chan *DiscoverEvent, 1)

	go zkRD.loopDiscover(zkRD.rootCtx, key, env)

	return env, nil
}

func (zkRD *ZkRegDiscover) loopDiscover(discvCtx context.Context, path string, env chan *DiscoverEvent) {
	for {
		discvEnv := &DiscoverEvent{
			Err: nil,
			Key: path,
		}

		data, _, watchEnv, err := zkRD.zkcli.GetW(path)
		if err != nil {
			fmt.Printf("fail to watch context for path(%s), err:%s\n", path, err.Error())
			if zkclient.ErrNoNode == err {
				fmt.Printf("node(%s) is not exist, will watch after 5s\n", path)
				time.Sleep(5 * time.Second)
				continue
			}

			discvEnv.Err = err
			env <- discvEnv
			time.Sleep(5 * time.Second)
			continue
		}

		discvEnv.Data = data

		// write into discoverEvent channel
		env <- discvEnv

		select {
		case <-discvCtx.Done():
			fmt.Printf("discover path(%s) done\n", path)
			return
		case <-watchEnv:
			fmt.Printf("watch found the content of path(%s) changed\n", path)
		}
	}
}
