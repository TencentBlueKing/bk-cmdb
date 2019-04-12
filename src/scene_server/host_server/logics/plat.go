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

package logics

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) IsPlatExist(ctx context.Context, cond mapstr.MapStr) (bool, errors.CCError) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		Fields:    []string{common.BKCloudIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDPlat, query)
	if err != nil {
		blog.Errorf("IsPlatExist http do error, err:%s, cond:%#v,rid:%s", err.Error(), cond, lgc.rid)
		return false, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("IsPlatExist http response error, err code:%d, err msg:%s, cond:%#v,rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return false, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if 1 == result.Data.Count {
		return true, nil
	}

	return false, nil
}
