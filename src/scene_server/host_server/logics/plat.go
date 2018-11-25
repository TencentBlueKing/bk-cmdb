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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) IsPlatExist(header http.Header, cond interface{}) (bool, error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	rid := util.GetHTTPCCRequestID(header)

	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     1,
		Sort:      common.BKCloudIDField,
		Fields:    common.BKCloudIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDPlat, header, query)
	if err != nil {
		blog.Errorf("IsPlatExist http do error, err:%s, cond:%+v,rid:%s", err.Error(), cond, rid)
		return false, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("IsPlatExist http response error, err:%s, cond:%+v,rid:%s", err.Error(), cond, rid)
		return false, defErr.New(result.Code, result.ErrMsg)
	}

	if 1 == result.Data.Count {
		return true, nil
	}

	return false, nil
}
