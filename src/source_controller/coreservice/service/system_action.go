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

package service

import (
	"crypto/md5"
	"encoding/hex"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/coreservice/core"
)

// GetSystemFlag get the system define flag
func (s *coreService) GetSystemFlag(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	ownerID := pathParams(common.BKOwnerIDField)
	flag := pathParams("flag")
	cond := make(map[string]interface{})

	h := md5.New()
	h.Write([]byte(flag))
	cipherStr := h.Sum(nil)
	cond[flag] = hex.EncodeToString(cipherStr) + ownerID

	var result interface{}
	err := s.db.Table(common.BKTableNameSystem).Find(cond).One(params.Context, &result)
	if nil != err {
		blog.Errorf("get system config error :%v, cond:%#v, rid: %s", err, cond, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectSelectInstFailed)
	}

	return result, nil
}
