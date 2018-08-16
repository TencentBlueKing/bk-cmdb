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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) GetProcessbyProcID(procID string, forward http.Header) (map[string]interface{}, error) {
	condition := map[string]interface{}{
		common.BKProcessIDField: procID,
	}

	reqParam := new(metadata.QueryInput)
	reqParam.Condition = condition
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDProc, forward, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		return nil, fmt.Errorf("get process by procID(%s) failed. err: %v, errcode: %d, errmsg: %s", procID, err, ret.Code, ret.ErrMsg)
	}

	if len(ret.Data.Info) < 1 {
		return nil, fmt.Errorf("there is no process with procID(%s)", procID)
	}

	return ret.Data.Info[0], nil
}
