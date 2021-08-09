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

package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type EtcdCli struct {
	etcdCli *clientv3.Client
}

// NewEtcdClient create a object of EtcdClient
func NewEtcdClient(etcdAddress string, timeOut time.Duration) (*EtcdCli, error) {
	etcdAddresses := strings.Split(etcdAddress, ",")
	etcdConf := clientv3.Config{
		Endpoints:   etcdAddresses,
		DialTimeout: timeOut,
	}
	cli, err := clientv3.New(etcdConf)
	if err != nil {
		return nil, err
	}
	return &EtcdCli{etcdCli: cli}, nil
}

// Stop stop the etcd
func (etcd *EtcdCli) Stop() error {
	etcd.etcdCli.Close()
	return nil
}

// Ping renew the lease
func (etcd *EtcdCli) Ping(leaseid clientv3.LeaseID) error {
	if _, err := etcd.etcdCli.KeepAlive(context.Background(), leaseid); err != nil {
		return err
	}
	return nil
}

// PutAndBindLease put a kv and bind lease
func (etcd *EtcdCli) PutAndBindLease(key, val string, ttl int64) (clientv3.LeaseID, error) {
	resp, err := etcd.etcdCli.Grant(context.Background(), ttl)
	if err != nil {
		return clientv3.NoLease, err
	}
	if _, err = etcd.etcdCli.Put(context.Background(), key, val, clientv3.WithLease(resp.ID)); err != nil {
		return clientv3.NoLease, err
	}
	return resp.ID, nil
}

// GetWithPrefix path as a prefix to get the key value
func (etcd *EtcdCli) GetWithPrefix(path string) ([]string, error) {
	rangeResp, err := etcd.etcdCli.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var values []string
	for _, kv := range rangeResp.Kvs {
		values = append(values, string(kv.Value))
	}
	return values, nil
}

// Get get the key value
func (etcd *EtcdCli) Get(path string) (string, error) {
	rangeResp, err := etcd.etcdCli.Get(context.Background(), path)
	if err != nil {
		return "", err
	}
	for _, kv := range rangeResp.Kvs {
		return string(kv.Value), nil
	}
	return "", nil
}

// Delete use a key to delete kv
func (e *EtcdCli) Delete(key string) error {
	if _, err := e.etcdCli.Delete(context.Background(), key); err != nil {
		return err
	}
	return nil
}

// Put put kv
func (e *EtcdCli) Put(key, val string) error {
	if _, err := e.etcdCli.Put(context.Background(), key, val); err != nil {
		return err
	}
	return nil
}
