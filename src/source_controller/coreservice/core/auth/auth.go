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

package auth

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.AuthOperation = (*authOperation)(nil)

type authOperation struct {
	dbProxy dal.DB
}

// New create a new instance manager instance
func New(dbProxy dal.DB) core.AuthOperation {
	return &authOperation{
		dbProxy: dbProxy,
	}
}

func (a *authOperation) SearchAuthResource(kit *rest.Kit, param metadata.PullResourceParam) (int64, []map[string]interface{}, errors.CCErrorCoder) {
	if param.Collection == "" {
		blog.ErrorJSON("search auth resource in empty mongo collection, param: %s, rid: %s", param, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "collection")
	}
	limit := param.Limit
	if limit > common.BKMaxPageSize && limit != common.BKNoLimit {
		blog.Errorf("search auth resource page limit %d exceeds max page size, rid: %s", limit, kit.Rid)
		return 0, nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}
	if limit == 0 {
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "page.limit")
	}
	f := a.dbProxy.Table(param.Collection).Find(param.Condition)
	count, err := f.Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("count auth resource failed, error: %s, input param: %s, rid: %s", err.Error(), param, kit.Rid)
		return 0, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(param.Fields) != 0 {
		f = f.Fields(param.Fields...)
	}

	var info []map[string]interface{}
	if param.Collection == common.BKTableNameBaseHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = f.Start(uint64(param.Offset)).Limit(uint64(limit)).All(kit.Ctx, &hosts)
		if err != nil {
			blog.ErrorJSON("search auth resource failed, error: %s, input param: %s, rid: %s", err.Error(), param, kit.Rid)
			return 0, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		info = make([]map[string]interface{}, len(hosts))
		for index, host := range hosts {
			info[index] = host
		}
	} else {
		info = make([]map[string]interface{}, 0)
		err = f.Start(uint64(param.Offset)).Limit(uint64(limit)).All(kit.Ctx, &info)
		if err != nil {
			blog.ErrorJSON("search auth resource failed, error: %s, input param: %s, rid: %s", err.Error(), param, kit.Rid)
			return 0, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
	}
	return int64(count), info, nil
}
