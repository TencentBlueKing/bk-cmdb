/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

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

	wg := sync.WaitGroup{}
	var err error
	results := make([]int64, len(filters))
	for idx, filter := range filters {
		wg.Add(1)
		go func(idx int, filter map[string]interface{}, ctx *rest.Contexts) {
			filter = util.SetQueryOwner(filter, ctx.Kit.SupplierAccount)
			var count uint64
			count, err = mongodb.Client().Table(table).Find(filter).Count(ctx.Kit.Ctx)
			if err != nil {
				blog.ErrorJSON("GetCountByFilter failed, error: %s, table: %s, filter: %s, rid: %s", err.Error(), table, filter, ctx.Kit.Rid)
				return
			}
			results[idx] = int64(count)
			wg.Done()
		}(idx, filter, ctx)
	}

	wg.Wait()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	ctx.RespEntity(results)
}
