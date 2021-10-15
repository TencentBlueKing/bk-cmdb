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

package registerdiscover

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common/blog"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type EventType int32

const (
	EventPut EventType = 0
	EventDel EventType = 1
)

// DiscoverEvent event of service discovery
type DiscoverEvent struct {
	Type	EventType
	Key    string
	Value  string
}

// KeyVal storage key and value in register and discover
type KeyVal struct {
	Key   string
	Value string
}

const (
	EtcdAuthUser = "cc"
	EtcdAuthPwd  = "3.0#bkcc"
)

// Config register and discover config
type Config struct {
	Host       string           // etcd host info
	User       string           // user name for authentication
	Passwd     string           // password relative to user
	TLS        *tls.Config      // tls config for https
}

// RegDiscv data structure of register and discover
type RegDiscv struct {
	client   *clientv3.Client
	election *concurrency.Election
	cancel   context.CancelFunc
	rootCxt  context.Context
}

// NewRegDiscv creates a register and discover object
func NewRegDiscv(config *Config) (*RegDiscv, error) {
	endpoints := strings.Split(config.Host, ",")
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("create regdiscv failed for no endpoints")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:        endpoints,
		DialTimeout:      time.Second * 5,
		AutoSyncInterval: time.Minute * 5,
		TLS:              config.TLS,
		Username:         EtcdAuthUser,
		Password:         EtcdAuthPwd,
	})
	if err != nil {
		return nil, err
	}

	regDiscv := &RegDiscv{
		client: client,
	}

	// create root context
	regDiscv.rootCxt, regDiscv.cancel = context.WithCancel(context.Background())

	return regDiscv, nil
}

// Ping verifies register and discover accessibility
func (rd *RegDiscv) Ping() error {
	if len (rd.client.Endpoints()) == 0 {
		return fmt.Errorf("etcd has no endpoint");
	}
	_, err := rd.client.Dial(rd.client.Endpoints()[0])
	return err
}

// Get gets corresponding value with key from register and discover
func (rd *RegDiscv) Get(key string) (string, error) {
	rsp, err := rd.client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(rsp.Kvs) == 0 {
		return "", fmt.Errorf("etcd get nothing from %s", key)
	}
	value := string(rsp.Kvs[0].Value)
	return value, nil
}

// GetWithPrefix gets corresponding value with key prefix from register and discover
func (rd *RegDiscv) GetWithPrefix(key string) ([]KeyVal, error) {
	rsp, err := rd.client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var values []KeyVal
	for _, kv := range rsp.Kvs {
		values = append(values, KeyVal{Key: string(kv.Key), Value: string(kv.Value)})
	}
	return values, nil
}

// Put puts key and value to register and discover
func (rd *RegDiscv) Put(key, val string) error {
	if _, err := rd.client.Put(context.Background(), key, val); err != nil {
		return err
	}
	return nil
}

// Delete deletes key and value from register and discover
func (rd *RegDiscv) Delete(key string) error {
	if _, err := rd.client.Delete(context.Background(), key); err != nil {
		return err
	}
	return nil
}

// Watch watches on a key or prefix. The watched events will be returned through the returned channel
func (rd *RegDiscv) Watch(ctx context.Context, key string) (<-chan *DiscoverEvent, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("invalid empty watch key")
	}

	ch := make(chan *DiscoverEvent, 10)

	go func() {
		evChan := ch
		watcher := clientv3.NewWatcher(rd.client)
		watchChan := watcher.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithPrevKV())
		for {
			select {
			case result := <-watchChan:
				if result.Err() != nil {
					blog.Errorf("watch key: %s, get err: %v", key, result.Err())
					continue
				}
				for _, event := range result.Events {
					blog.Infof("watch key: %s, get event: %+v", key, event)
					discoverEvent := new(DiscoverEvent)
					switch event.Type {
					case clientv3.EventTypePut:
						discoverEvent.Type = EventPut
						discoverEvent.Key = string(event.Kv.Key)
						discoverEvent.Value = string(event.Kv.Value)
						evChan <- discoverEvent
					case clientv3.EventTypeDelete:
						discoverEvent.Type = EventDel
						discoverEvent.Key = string(event.Kv.Key)
						// delete event return previous value before delete operation
						discoverEvent.Value = string(event.PrevKv.Value)
						evChan <- discoverEvent
					default:
						blog.Warnf("watch key: %s, get unknown event type: %v", key, event.Type)
					}
				}
			case <-ctx.Done():
				blog.Infof("watch stopped because of context done, key: %s", key)
				return
			}
		}
	}()

	return ch, nil
}

// Campaign puts a value as eligible for the election. It blocks until
// it is elected, an error occurs, or the context is cancelled.
func (rd *RegDiscv) Campaign(key, val string) error {
	// create session for election competition
	session, err := concurrency.NewSession(rd.client)
	if err != nil {
		return err
	}
	defer session.Close()
	rd.election = concurrency.NewElection(session, key)

	if err := rd.election.Campaign(context.Background(), val); err != nil {
		return err
	}

	return nil
}

// Resign lets a leader start a new election
func (rd *RegDiscv) Resign() error {
	if rd.election == nil {
		return nil
	}
	return rd.election.Resign(context.Background())
}

// RegisterAndKeepAlive registers a kv and keeps its lease alive
func (rd *RegDiscv) RegisterAndKeepAlive(key , val string) error {
	go func() {
		lease := clientv3.NewLease(rd.client)
		var curLeaseId clientv3.LeaseID = 0
		for {
			select {
			case <-rd.rootCxt.Done():
				blog.Infof("register and keep alive done, key: %s", key)
				return
			default:
				// go through the following procedure
			}

			if curLeaseId == 0 {
				// grant the lease with 5s
				leaseResp, err := lease.Grant(context.Background(), 5)
				if err != nil {
					blog.Errorf("failed to grant lease, key: %s, err: %v", key, err)
					time.Sleep(1 * time.Second)
					continue
				}
				_, err = rd.client.Put(context.Background(), key, val, clientv3.WithLease(leaseResp.ID))
				if err != nil {
					blog.Errorf("put with lease err: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}
				curLeaseId = leaseResp.ID
			} else {
				// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
				if _, err := lease.KeepAliveOnce(context.Background(), curLeaseId); err != nil {
					blog.Errorf("keep alive lease err: %v", err)
					curLeaseId = 0
					time.Sleep(1 * time.Second)
					continue
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}

// Cancel stops register and discover
func (rd *RegDiscv) Cancel() {
	rd.cancel()
}