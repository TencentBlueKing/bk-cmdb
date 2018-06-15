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
	"errors"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) GetResoulePoolModuleID(pheader http.Header, condition interface{}) (int64, error) {
	query := metadata.QueryInput{
		Start:     0,
		Limit:     1,
		Sort:      common.BKModuleIDField,
		Fields:    common.BKModuleIDField,
		Condition: condition,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKModuleIDField, pheader, &query)
	if err != nil || (err == nil && !result.Result) {
		return -1, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return -1, errors.New("search resource pool, but get nil data")
	}

	return result.Data.Info[0].Int64(common.BKModuleIDField)
}
