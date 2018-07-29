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
	"context"
	"fmt"
	"log"
	"testing"
)

type MethodInterface interface {
	Test(a int, b *int) error
}

type MethodTest struct {
}

func NewInstance() MethodInterface {
	return &Method2Test{}
}

func (m *MethodTest) Test(a int, b *int) error {
	*b = a + 1
	return nil
}

type Method2Test struct {
}

func (m *Method2Test) Test(a int, b *int) error {
	*b = a * 2
	return nil
}

func TestRPCServer(t *testing.T) {

	cfg := Config{IPAddr: "127.0.0.1", Port: 2323}
	svr := NewServer(cfg, []interface{}{NewInstance()})

	// start server
	go svr.Run(context.Background())

	// start client
	cli := NewClient(cfg)
	//连接远程rpc服务
	var ret int
	err := cli.Call("Method2Test.Test", 100, &ret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", ret)

}
