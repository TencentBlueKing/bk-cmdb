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
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (s *coreService) GetHostModulesIDs(ctx *rest.Contexts) {
	dat := &meta.ModuleHostConfigParams{}
	if err := ctx.DecodeInto(dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	condition := map[string]interface{}{common.BKAppIDField: dat.ApplicationID, common.BKHostIDField: dat.HostID}
	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	moduleIDs, err := s.getModuleIDsByHostID(ctx.Kit, condition)
	if nil != err {
		blog.Errorf("get host module id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrGetModule))
		return
	}
	ctx.RespEntity(moduleIDs)
}

func (s *coreService) getModuleIDsByHostID(kit *rest.Kit, moduleCond interface{}) ([]int64, error) {
	result := make([]meta.ModuleHost, 0)
	var ret []int64

	err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(moduleCond).Fields(common.BKModuleIDField).All(kit.Ctx, &result)
	if nil != err {
		blog.Errorf("get module id by host id failed, error: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	for _, r := range result {
		ret = append(ret, r.ModuleID)
	}
	return ret, nil
}
