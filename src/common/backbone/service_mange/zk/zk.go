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

package zk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common/zkclient"
)

// ZkClient do service register and discover by zookeeper
type ZkClient struct {
	zkCli          *zkclient.ZkClient
	cancel         context.CancelFunc
	rootCxt        context.Context
	sessionTimeOut time.Duration
}

// NewZkClient create a object of ZkClient
func NewZkClient(zkAddress string, timeOut time.Duration) *ZkClient {
	zkAddresses := strings.Split(zkAddress, ",")
	return &ZkClient{
		zkCli:          zkclient.NewZkClient(zkAddresses),
		sessionTimeOut: timeOut,
	}
}

// Ping to ping server
func (zk *ZkClient) Ping() error {
	return zk.zkCli.Ping()
}

// Start used to run register and discover server
func (zk *ZkClient) Start() error {
	// connect zookeeper
	if err := zk.zkCli.ConnectEx(zk.sessionTimeOut); err != nil {
		return fmt.Errorf("fail to connect zookeeper, err: %+v", err)
	}

	// create root context
	zk.rootCxt, zk.cancel = context.WithCancel(context.Background())

	return nil
}

// Stop used to stop register and discover server
func (zk *ZkClient) Stop() error {
	// close the connection of zookeeper
	zk.zkCli.Close()

	zk.cancel()

	return nil
}

// Client return zk client
func (zk *ZkClient) Client() *zkclient.ZkClient {
	return zk.zkCli
}

// SessionTimeOut client session time out
func (zk *ZkClient) SessionTimeOut() time.Duration {
	return zk.sessionTimeOut
}

// WithCancel context with cancel
func (zk *ZkClient) WithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(zk.rootCxt)
}
