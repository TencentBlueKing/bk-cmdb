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

package rpc

import (
	"fmt"
	"net/rpc"
)

// Client rpc client methods
type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
	Close() error
}

// NewClient create a new client instance
func NewClient(cfg Config) Client {
	return &client{
		conf: cfg,
	}
}

type client struct {
	rpcClient *rpc.Client
	conf      Config
}

func (c *client) Open() error {
	cli, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", c.conf.IPAddr, c.conf.Port))
	if nil != err {
		return err
	}
	c.rpcClient = cli
	return nil
}

func (c *client) Call(serviceMethod string, args interface{}, reply interface{}) error {

	if nil == c.rpcClient {
		if err := c.Open(); nil != err {
			return err
		}
	}

	err := c.rpcClient.Call(serviceMethod, args, reply)
	if nil != err {
		c.rpcClient.Close()
		c.rpcClient = nil
		return err
	}

	return nil
}

func (c *client) Close() error {
	if nil != c.rpcClient {
		return c.rpcClient.Close()
	}

	return nil
}
