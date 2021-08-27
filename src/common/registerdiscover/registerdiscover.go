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
	EVENT_PUT		EventType = 0
	EVENT_DEL		EventType = 1
)

// DiscoverEvent if servers chenged, will create a discover event
type DiscoverEvent struct {
	Err    error
	Type	EventType
	Key    string
	Value  string
}

type KeyVal struct {
	Key   string
	Value string
}

const (
	ETCD_AUTH_USER = "cc"
	ETCD_AUTH_PWD  = "3.0#bkcc"
)

// Config etcd register and discover config
type Config struct {
	Host       string           // etcd host info
	User       string           // user name for authentication
	Passwd     string           // password relative to user
	TLS        *tls.Config      // tls config for https
}

// RegDiscv is data struct of register-discover
type RegDiscv struct {
	client   *clientv3.Client
	election *concurrency.Election
	cancel   context.CancelFunc
	rootCxt  context.Context
}

// NewRegDiscv create an object of RegDiscv
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
		Username:         ETCD_AUTH_USER,
		Password:         ETCD_AUTH_PWD,
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

func (rd *RegDiscv) Ping() error {
	if len (rd.client.Endpoints()) == 0 {
		return fmt.Errorf("etcd has no endpoint");
	}
	_, err := rd.client.Dial(rd.client.Endpoints()[0])
	return err
}

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

func (rd *RegDiscv) Put(key, val string) error {
	if _, err := rd.client.Put(context.Background(), key, val); err != nil {
		return err
	}
	return nil
}

// Delete use a key to delete kv
func (rd *RegDiscv) Delete(key string) error {
	if _, err := rd.client.Delete(context.Background(), key); err != nil {
		return err
	}
	return nil
}

func (rd *RegDiscv) Watch(ctx context.Context, key string) (<-chan *DiscoverEvent, error) {
	ch := make(chan *DiscoverEvent, 1)
	if ch == nil {
		return nil, fmt.Errorf("expected watcher channel, got nil")
	}

	go func() {
		evChan := ch
		watcher := clientv3.NewWatcher(rd.client)
		watchChan := watcher.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithPrevKV())
		for {
			select {
			case result := <-watchChan:
				if result.Err() != nil {
					blog.Errorf("watch get err: %v", result.Err())
					continue
				}
				for _, event := range result.Events {
					blog.Infof("watch key(%s) event(%v)\n", key, event)
					discvEvent := new(DiscoverEvent)
					switch event.Type {
					case clientv3.EventTypePut:
						discvEvent.Err = nil
						discvEvent.Type = EVENT_PUT
						discvEvent.Key = string(event.Kv.Key)
						discvEvent.Value = string(event.Kv.Value)
					case clientv3.EventTypeDelete:
						discvEvent.Err = nil
						discvEvent.Type = EVENT_DEL
						discvEvent.Key = string(event.Kv.Key)
						// delete event return previous value before delete operation
						discvEvent.Value = string(event.PrevKv.Value)
					default:
						discvEvent.Err = fmt.Errorf("discover unknown event type (%v)", event.Type)
					}
					evChan <- discvEvent
				}
			case <-ctx.Done():
				blog.Infof("watch stopped because of context done")
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

// RegisterAndKeepAlive register a kv and keep its lease alive
func (rd *RegDiscv) RegisterAndKeepAlive(key , val string) error {
	go func() {
		lease := clientv3.NewLease(rd.client)
		var curLeaseId clientv3.LeaseID = 0
		for {
			if curLeaseId == 0 {
				// grant the lease with 10s
				leaseResp, err := lease.Grant(context.Background(), 10)
				if err != nil {
					blog.Errorf("grant lease err: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}
				_, err = rd.client.Put(context.Background(), key, val, clientv3.WithLease(leaseResp.ID))
				if err != nil {
					blog.Errorf("put with lease err: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}
				curLeaseId = leaseResp.ID
			} else {
				// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
				if _, err := lease.KeepAliveOnce(context.Background(), curLeaseId); err != nil {
					blog.Errorf("keep alive lease err: %v", err)
					curLeaseId = 0
					time.Sleep(3 * time.Second)
					continue
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return nil
}

func (rd *RegDiscv) Cancel() {
	rd.cancel()
}