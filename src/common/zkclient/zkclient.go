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
	"errors"
	//"bcs/bcs-common/common/blog"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	ErrNoNode        = zk.ErrNoNode
	EventNodeDeleted = zk.EventNodeDeleted
	ErrNodeExists    = zk.ErrNodeExists
)

const (
	AUTH_USER = "cc"
	AUTH_PWD  = "3.0#bkcc"
)

type ZkLock struct {
	zkHost []string
	zkConn *zk.Conn
	zkAcl  []zk.ACL
	zkLock *zk.Lock
}

func NewZkLock(host []string) *ZkLock {
	zlock := ZkLock{
		zkHost: host[:],
		zkConn: nil,
		zkLock: nil,
		zkAcl:  zk.DigestACL(zk.PermAll, AUTH_USER, AUTH_PWD),
	}
	return &zlock
}

func (zlock *ZkLock) Lock(path string) error {
	return zlock.LockEx(path, time.Second*5)
}

func (zlock *ZkLock) LockEx(path string, sessionTimeOut time.Duration) error {
	if zlock.zkConn == nil {
		conn, _, connErr := zk.Connect(zlock.zkHost, sessionTimeOut)
		if connErr != nil {
			return connErr
		}

		//auth
		auth := AUTH_USER + ":" + AUTH_PWD
		if err := conn.AddAuth("digest", []byte(auth)); err != nil {
			conn.Close()
			return err
		}

		zlock.zkConn = conn
	}

	lock := zk.NewLock(zlock.zkConn, path, zlock.zkAcl)
	if lock == nil {
		return fmt.Errorf("fail to new lock for path(%s)", path)
	}

	zlock.zkLock = lock

	return zlock.zkLock.Lock()
}

func (zlock *ZkLock) UnLock() error {
	if zlock.zkLock != nil {
		if err := zlock.zkLock.Unlock(); err != nil {
			zlock.zkConn.Close()
			return err
		}
	}

	if zlock.zkConn != nil {
		zlock.zkConn.Close()
	}

	return nil
}

type ZkClient struct {
	ZkHost []string
	ZkConn *zk.Conn
	zkAcl  []zk.ACL
}

func NewZkClient(host []string) *ZkClient {
	c := ZkClient{
		ZkHost: host[:],
		ZkConn: nil,
		zkAcl:  zk.DigestACL(zk.PermAll, AUTH_USER, AUTH_PWD),
	}

	return &c
}

func (z *ZkClient) Connect() error {
	return z.ConnectEx(time.Second * 60)
}

func (z *ZkClient) ConnectEx(sessionTimeOut time.Duration) error {
	if z.ZkConn != nil {
		z.Close()
	}

	c, _, err := zk.Connect(z.ZkHost, sessionTimeOut)
	if err != nil {
		return err
	}

	// AddAuth
	auth := AUTH_USER + ":" + AUTH_PWD
	if err := c.AddAuth("digest", []byte(auth)); err != nil {
		c.Close()
		return err
	}

	z.ZkConn = c
	return nil
}

func (z *ZkClient) Close() {
	if nil != z.ZkConn {
		z.ZkConn.Close()
	}
}

// Ping to ping server
func (z *ZkClient) Ping() error {
	if z.ZkConn.State() != zk.StateHasSession {
		return errors.New("not health")
	}
	return nil
}

func (z *ZkClient) Get(path string) (string, error) {
	data, _, err := z.ZkConn.Get(path)
	if err == zk.ErrNoAuth {
		if err = z.AddAuth(); err != nil {
			return string(data), err
		}
		data, _, err = z.ZkConn.Get(path)
	}
	return string(data), err
}

func (z *ZkClient) AddAuth() error {
	auth := AUTH_USER + ":" + AUTH_PWD
	return z.ZkConn.AddAuth("digest", []byte(auth))
}

func (z *ZkClient) GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	data, statm, ch, err := z.ZkConn.GetW(path)
	if err == zk.ErrNoAuth {
		if err = z.AddAuth(); err != nil {
			return data, statm, ch, err
		}
		data, statm, ch, err = z.ZkConn.GetW(path)
	}
	return data, statm, ch, err
}

func (z *ZkClient) GetEx(path string) ([]byte, *zk.Stat, error) {
	return z.ZkConn.Get(path)
}

func (z *ZkClient) GetChildren(path string) ([]string, error) {
	data, _, err := z.ZkConn.Children(path)
	return data, err
}

func (z *ZkClient) GetChildrenEx(path string) ([]string, *zk.Stat, error) {
	return z.ZkConn.Children(path)
}

func (z *ZkClient) WatchChildren(path string) ([]string, <-chan zk.Event, error) {
	data, _, env, err := z.ZkConn.ChildrenW(path)
	if err == zk.ErrNoAuth {
		if err = z.AddAuth(); err != nil {
			return data, env, err
		}
		data, _, env, err = z.ZkConn.ChildrenW(path)
	}
	return data, env, err
}

func (z *ZkClient) ChildrenW(path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	return z.ZkConn.ChildrenW(path)
}

func (z *ZkClient) Set(path, data string, version int32) error {
	_, err := z.ZkConn.Set(path, []byte(data), version)
	return err
}

func (z *ZkClient) Del(path string, version int32) error {
	err := z.ZkConn.Delete(path, version)
	return err
}

func (z *ZkClient) Exist(path string) (bool, error) {
	b, _, err := z.ZkConn.Exists(path)
	return b, err
}

func (z *ZkClient) ExistEx(path string) (bool, *zk.Stat, error) {
	return z.ZkConn.Exists(path)
}

func (z *ZkClient) ExistW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
	return z.ZkConn.ExistsW(path)
}

func (z *ZkClient) Create(path string, data []byte) error {
	_, err := z.ZkConn.Create(path, data, 0, z.zkAcl)
	return err
}

func (z *ZkClient) State() zk.State {
	return z.ZkConn.State()
}

//CreateEphemeral create ephemeral node
func (z *ZkClient) CreateEphAndSeq(path string, data []byte) error {
	tmpPath := strings.Split(path, "/")
	if len(tmpPath) > 2 {
		rootPath := strings.Join(tmpPath[0:len(tmpPath)-1], "/")
		b, _ := z.Exist(rootPath)
		if !b {
			if err := z.CreateDeepNode(rootPath, []byte("")); err != nil {
				return err
			}
		}
	}

	_, err := z.ZkConn.CreateProtectedEphemeralSequential(path, data, z.zkAcl)
	if err == zk.ErrNoAuth {
		if err = z.AddAuth(); err != nil {
			return err
		}
		_, err = z.ZkConn.CreateProtectedEphemeralSequential(path, data, z.zkAcl)
	}
	return err
}

func (z *ZkClient) CreateEphAndSeqEx(path string, data []byte) (string, error) {
	tmpPath := strings.Split(path, "/")
	if len(tmpPath) > 2 {
		rootPath := strings.Join(tmpPath[0:len(tmpPath)-1], "/")
		b, _ := z.Exist(rootPath)
		if !b {
			if err := z.CreateDeepNode(rootPath, []byte("")); err != nil {
				return "", err
			}
		}
	}

	return z.ZkConn.CreateProtectedEphemeralSequential(path, data, z.zkAcl)
}

func (z *ZkClient) Update(path, data string) error {
	b, _ := z.Exist(path)

	if b {
		err := z.Set(path, data, -1)
		if err != nil {
			return fmt.Errorf("fail to set value to path(%s) b(%t), err(%s)", path, b, err.Error())
		}
	} else {
		err := z.CreateDeepNode(path, []byte(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *ZkClient) CreateDeepNode(path string, data []byte) error {
	nodes := strings.Split(path, "/")
	tmpPath := ""
	ctx := []byte("")
	for index, nd := range nodes {
		if nd == "" {
			continue
		}

		if index+1 == len(nodes) {
			ctx = data
		}

		tmpPath += "/" + nd
		err := z.CreateNode(tmpPath, ctx)
		if err != nil {
			return fmt.Errorf("fail to create node(%s), err(%s)", tmpPath, err.Error())
		}
	}
	/*
		if err := z.Set(path, string(data), -1); err != nil {
			return fmt.Errorf("fail to set value to path(%s), err(%s)", path, err.Error())
		}
	*/
	return nil
}

func (z *ZkClient) CreateNode(path string, data []byte) error {
	fmt.Printf("creating path %v\n", path)
	bExist, err := z.Exist(path)
	if err != nil {
		return err
	}

	fmt.Printf("path %v not exist,will create now\n", path)
	if !bExist {
		err := z.Create(path, data)
		if err == zk.ErrNodeExists {
			err = z.Set(path, string(data), -1)
		}
		if err != nil {
			fmt.Printf("path %v create failed :%v\n", path, err)
			return err
		}
	}

	fmt.Printf("path %v created\n", path)

	return nil
}

func (z *ZkClient) CheckNode(path string, data []byte) error {
	exist, _ := z.Exist(path)
	if exist == false {
		err := z.Create(path, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *ZkClient) CheckMulNode(path string, data []byte) error {
	var tempPath = ""
	temp := strings.Split(path, "/")
	for i := 1; i < len(temp); i++ {
		if temp[i] == "" {
			continue
		}
		tempPath += "/" + temp[i]
		err := z.CheckNode(tempPath, []byte(temp[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *ZkClient) GetAll2Json(path string) (string, error) {
	childs, err := z.GetChildren(path)
	if err != nil {
		//blog.Warn("fail to get children from path(%s). err:%s", path, err.Error())
		return "", err
	}

	if len(childs) <= 0 {
		ctx, getErr := z.Get(path)
		if getErr != nil {
			//blog.Warn("fail to get value from path(%s), err:%s", path, err.Error())
			return "", getErr
		}

		return ctx, nil
	}

	mpChilds := make(map[string]string)

	for _, child := range childs {
		chPath := path + "/" + child
		val, _ := z.GetAll2Json(chPath)
		mpChilds[child] = val
		//blog.Info("children path(%s), value(%s)", chPath, val)
	}

	data, err := json.Marshal(mpChilds)

	//blog.Info("data:%s", string(data))
	return string(data), err
}
