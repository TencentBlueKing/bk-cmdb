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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) GetHostModulesIDs(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	dat := &meta.ModuleHostConfigParams{}
	if err := data.MarshalJSONInto(dat); err != nil {
		blog.Errorf("get host module id failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	condition := map[string]interface{}{common.BKAppIDField: dat.ApplicationID, common.BKHostIDField: dat.HostID}
	condition = util.SetModOwner(condition, params.SupplierAccount)
	moduleIDs, err := s.getModuleIDsByHostID(params, condition)
	if nil != err {
		blog.Errorf("get host module id failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrGetModule)
	}

	return moduleIDs, nil
}

func (s *coreService) getModuleIDsByHostID(params core.ContextParams, moduleCond interface{}) ([]int64, error) {
	result := make([]meta.ModuleHost, 0)
	var ret []int64

	err := s.db.Table(common.BKTableNameModuleHostConfig).Find(moduleCond).Fields(common.BKModuleIDField).All(params.Context, &result)
	if nil != err {
		blog.Errorf("get module id by host id failed, error: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	for _, r := range result {
		ret = append(ret, r.ModuleID)
	}
	return ret, nil
}
