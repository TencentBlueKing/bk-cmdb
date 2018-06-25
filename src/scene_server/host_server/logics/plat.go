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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"context"
	"fmt"
)

func (lgc *Logics) IsPlatExist(pheader http.Header, cond interface{}) (bool, error) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     1,
		Sort:      common.BKAppIDField,
		Fields:    common.BKAppIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDPlat, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return false, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	if 1 == result.Data.Count {
		return true, nil
	}

	return false, nil
}
