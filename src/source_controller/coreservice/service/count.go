/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// GetCountByFilter TODO
// get counts in table based on filters, returns in the same order
func (s *coreService) GetCountByFilter(ctx *rest.Contexts) {
	req := struct {
		Table   string                   `json:"table"`
		Filters []map[string]interface{} `json:"filters"`
	}{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}
	filters := req.Filters
	table := req.Table

	// to speed up, multi goroutine to get count
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	results := make([]int64, len(filters))
	for idx, filter := range filters {
		pipeline <- true
		wg.Add(1)
		go func(idx int, filter map[string]interface{}) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			filter = util.SetQueryOwner(filter, ctx.Kit.SupplierAccount)
			count, err := mongodb.Client().Table(table).Find(filter).Count(ctx.Kit.Ctx)
			if err != nil {
				blog.ErrorJSON("GetCountByFilter failed, error: %s, table: %s, filter: %s, rid: %s", err.Error(), table, filter, ctx.Kit.Rid)
				if firstErr == nil {
					firstErr = ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
				}
				return
			}

			lock.Lock()
			results[idx] = int64(count)
			lock.Unlock()

		}(idx, filter)
	}

	wg.Wait()

	if firstErr != nil {
		ctx.RespAutoError(firstErr)
		return
	}

	ctx.RespEntity(results)
}
