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
	"errors"
	"fmt"
	"strconv"

	// "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) SearchMainlineModelTopo(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	return s.core.TopoOperation().SearchMainlineModelTopo()
}

func (s *coreService) SearchMainlineInstanceTopo(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bkBizID := pathParams("bk_biz_id")
	if len(bkBizID) == 0 {
		return nil, errors.New("bk_biz_id field empty")
	}
	bizID, err := strconv.ParseInt(bkBizID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bk_biz_id field invalid, %v", err)
	}

	// TODO add parse withDetail option
	return s.core.TopoOperation().SearchMainlineInstanceTopo(bizID, false)
}
