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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// GetSystemUserConfig TODO
func (s *coreService) GetSystemUserConfig(ctx *rest.Contexts) {
	ctx.RespEntityWithError(s.core.SystemOperation().GetSystemUserConfig(ctx.Kit))
}

// UpdateGlobalConfig update global setting config.
func (s *coreService) UpdateGlobalConfig(ctx *rest.Contexts) {
	typeId := ctx.Request.PathParameter("type")
	input := make(mapstr.MapStr)
	if jsErr := ctx.DecodeInto(&input); nil != jsErr {
		ctx.RespAutoError(jsErr)
		return
	}

	err := s.core.SystemOperation().UpdatePlatformSettingConfig(ctx.Kit, input, typeId)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// SearchGlobalConfig search global setting config.
func (s *coreService) SearchGlobalConfig(ctx *rest.Contexts) {
	options := new(metadata.GlobalConfOptions)
	if err := ctx.DecodeInto(options); err != nil {
		ctx.RespAutoError(err)
		return
	}

	conf, err := s.core.SystemOperation().SearchGlobalSettingConfig(ctx.Kit, options)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(conf)
}

// GetHostSnapDataID get host snap data id
func (s *coreService) GetHostSnapDataID(ctx *rest.Contexts) {
	cond := map[string]interface{}{
		common.BKFieldDBID: "gse_data_id",
	}

	dataIDInfo := new(metadata.DataIDInfo)
	err := mongodb.Shard(ctx.Kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).
		Fields("host_snap."+ctx.Kit.TenantID).One(ctx.Kit.Ctx, &dataIDInfo)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("get host snap data id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if len(dataIDInfo.HostSnap) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "host_snap"))
		return
	}

	dataID, exists := dataIDInfo.HostSnap[ctx.Kit.TenantID]
	if !exists {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "host_snap"))
		return
	}

	ctx.RespEntity(dataID)
}
