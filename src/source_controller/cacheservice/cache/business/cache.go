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

package business

import (
	"context"
	"fmt"
	"sync"

	"configcenter/src/common"
	"configcenter/src/source_controller/cacheservice/cache/tools"
	"configcenter/src/storage/reflector"
)

var client *Client
var clientOnce sync.Once
var cache *cacheCollection

func NewClient() *Client {

	if client != nil {
		return client
	}

	clientOnce.Do(func() {
		client = &Client{
			lock: tools.NewRefreshingLock(),
		}
	})

	return client
}

// Attention, it can only be called for once.
func NewCache(event reflector.Interface) error {

	if cache != nil {
		return nil
	}

	// cache has not been initialized.
	biz := &business{
		key:   bizKey,
		event: event,
	}

	if err := biz.Run(); err != nil {
		return fmt.Errorf("run biz cache failed, err: %v", err)
	}

	module := &moduleSet{
		key:        moduleKey,
		collection: common.BKTableNameBaseModule,
		event:      event,
		lock:       tools.NewRefreshingLock(),
	}
	if err := module.Run(); err != nil {
		return fmt.Errorf("run module cache failed, err: %v", err)
	}

	set := &moduleSet{
		key:        setKey,
		collection: common.BKTableNameBaseSet,
		event:      event,
		lock:       tools.NewRefreshingLock(),
	}
	if err := set.Run(); err != nil {
		return fmt.Errorf("run set cache failed, err: %v", err)
	}

	custom := &customLevel{
		key:         customKey,
		event:       event,
		lock:        tools.NewRefreshingLock(),
		customWatch: make(map[string]context.CancelFunc),
	}
	if err := custom.Run(); err != nil {
		return fmt.Errorf("run biz custom level cache failed, err: %v", err)
	}

	cache = &cacheCollection{
		business:    biz,
		set:         set,
		module:      module,
		customLevel: custom,
	}
	return nil
}
