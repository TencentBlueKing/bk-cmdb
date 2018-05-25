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
 
package zkclient

import (
	//	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDataValue(t *testing.T) {
	fmt.Println("TEST Data Value")
	zkClient := NewZkClient([]string{"127.0.0.1:2181"})

	defer zkClient.Close()

	err := zkClient.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = zkClient.Create("/data1", []byte("DATA")); err != nil {
		fmt.Println(err)
	}

	result, err := zkClient.Get("/data1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)

	if err = zkClient.Create("/data1/ip/act", []byte("act")); err != nil {
		fmt.Println(err)
	}

	result, err = zkClient.Get("/data1/ip/act")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func TestZkLock(t *testing.T) {
	fmt.Println("TEST zk lock")
	zkLock1 := NewZkLock([]string{"127.0.0.1:2181"})

	if err := zkLock1.Lock("/lock"); err != nil {
		fmt.Printf("lock1 fail lock. err:%s \n", err.Error())
	} else {
		fmt.Println("lock1 lock")
	}

	zkLock1.UnLock()

}

func Test_WatchChildren(t *testing.T) {
	t.Log("----- start test WatchChildren -----")

	zkClient := NewZkClient([]string{"127.0.0.1:2181"})

	err := zkClient.Connect()
	if err != nil {
		t.Log("fail to connect zk, err:%s", err.Error())
		return
	}

	defer zkClient.Close()

	path := "/admin/delete_topics"

	t.Logf("start to watch child(%s)", path)

	go func() {
		for {
			childs, watchEnv, watchErr := zkClient.WatchChildren(path)
			if watchErr != nil {
				t.Logf("fail to watch children(%s), err:%s", path, watchErr.Error())
				return
			}

			//time.Sleep(2 * time.Second)

			t.Logf("childs:%+v, watchEnv:%v\n", childs, watchEnv)
			env := <-watchEnv
			t.Logf("env: %v, child:%+v\n", env, childs)
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Printf("begin to create eph node\n")
	newPath, err := zkClient.CreateEphAndSeqEx(path+"/node", []byte("test"))
	if err != nil {
		t.Logf("create eph node failed. err:%s", err.Error())
		return
	}

	t.Logf("newPath: %s", newPath)

	time.Sleep(10 * time.Second)

	t.Logf("----- end test WatchChildren -----")

}
