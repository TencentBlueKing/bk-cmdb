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

package hostsnap

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func newFilter() *filter {
	f := &filter{
		pool: make(map[string]int64),
		// 30 minutes, with range minutes [30, 90]
		ttlSeconds:      60 * 60,
		ttlRangeSeconds: [2]int{-30 * 60, 30 * 60},
	}
	go f.gc()
	return f
}

// filter is used to filter out the ip which is not register to cmdb,
// but still report snapshot data to cmdb, this will forced us to search it
// from mongodb again and again, because we always can not find it in the cache.
type filter struct {
	lock sync.Mutex
	// key: ip:cloud_id
	// v: unix time for key expire ttl seconds.
	pool       map[string]int64
	ttlSeconds int64
	// min:[0], max:[1]
	ttlRangeSeconds [2]int
}

// only set the invalid ip which can not found in db.
func (f *filter) Set(ip string, cloudID int64) {
	key := fmt.Sprintf("%s:%d", ip, cloudID)
	f.lock.Lock()
	// give a random ttl to avoid refresh the key concurrently.
	f.pool[key] = f.randomTTL()
	f.lock.Unlock()
}

// if exist, then do not use this ip to search in mongodb.
func (f *filter) Exist(ip string, cloudID int64) bool {
	key := fmt.Sprintf("%s:%d", ip, cloudID)
	f.lock.Lock()
	ttl, exist := f.pool[key]
	f.lock.Unlock()

	if !exist {
		return false
	}

	if time.Now().Unix()-ttl > f.ttlSeconds {
		// force to update the cache in the next round
		return false
	}

	return true
}

func (f *filter) gc() {
	for {
		time.Sleep(time.Duration(f.ttlSeconds) * time.Second)
		f.lock.Lock()
		for host, ttl := range f.pool {
			if time.Now().Unix()-ttl > f.ttlSeconds {
				// remove from the cache.
				delete(f.pool, host)
			}
		}
		f.lock.Unlock()
	}
}

func (f *filter) randomTTL() int64 {
	rand.Seed(time.Now().UnixNano())
	seconds := rand.Intn(f.ttlRangeSeconds[1]-f.ttlRangeSeconds[0]) + f.ttlRangeSeconds[0]
	return time.Now().Unix() + int64(seconds)

}
