/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package modelquote

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// BatchCreateQuotedInstance batch create quoted instances.
func BatchCreateQuotedInstance(cts *rest.Contexts) {
	instances := make([]mapstr.MapStr, 0)
	if err := cts.DecodeInto(&instances); err != nil {
		cts.RespAutoError(err)
		return
	}

	if len(instances) == 0 || len(instances) > common.BKMaxLimitSize {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrArrayLengthWrong, "instances", common.BKMaxLimitSize))
		return
	}

	objID := cts.Request.PathParameter(common.BKObjIDField)
	if len(objID) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	if err := validateCreateQuotedInstances(cts.Kit, objID, instances); err != nil {
		cts.RespAutoError(err)
		return
	}

	table := common.GetInstTableName(objID, cts.Kit.TenantID)
	ids, err := mongodb.Client().NextSequences(cts.Kit.Ctx, table, len(instances))
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	now := time.Now()

	for idx := range instances {
		instances[idx].Set(common.BKFieldID, ids[idx])
		instances[idx].Set(common.TenantID, cts.Kit.TenantID)
		instances[idx].Set(common.CreateTimeField, now)
		instances[idx].Set(common.LastTimeField, now)
	}

	err = mongodb.Client().Table(table).Insert(cts.Kit.Ctx, instances)
	if err != nil {
		blog.Errorf("create quoted instances failed, err: %v, data: %+v, rid: %v", err, instances, cts.Kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err)))
			return
		}
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}

	cts.RespEntity(metadata.BatchCreateResult{IDs: ids})
}

// ListQuotedInstance list quoted instances.
func ListQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	objID := cts.Request.PathParameter(common.BKObjIDField)
	if len(objID) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}
	table := common.GetInstTableName(objID, cts.Kit.TenantID)

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(table).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count quoted instances failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.InstDataInfo{Count: int(count)})
		return
	}

	instances := make([]mapstr.MapStr, 0)
	err = mongodb.Client().Table(table).Find(filter).Start(uint64(opt.Page.Start)).Limit(uint64(opt.Page.Limit)).
		Fields(opt.Fields...).All(cts.Kit.Ctx, &instances)
	if err != nil {
		blog.Errorf("list quoted instances failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(metadata.InstDataInfo{Info: instances})
}

// BatchUpdateQuotedInstance batch update quoted instances.
func BatchUpdateQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.CommonUpdateOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	objID := cts.Request.PathParameter(common.BKObjIDField)
	if len(objID) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	if err := validateUpdateQuotedInst(cts.Kit, objID, opt.Data); err != nil {
		cts.RespAutoError(err)
		return
	}

	table := common.GetInstTableName(objID, cts.Kit.TenantID)

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	opt.Data.Set(common.LastTimeField, time.Now())
	err = mongodb.Client().Table(table).Update(cts.Kit.Ctx, filter, opt.Data)
	if err != nil {
		blog.Errorf("list quoted instances failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err)))
			return
		}
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	cts.RespEntity(nil)
}

// BatchDeleteQuotedInstance batch delete quoted instances.
func BatchDeleteQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.CommonFilterOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	objID := cts.Request.PathParameter(common.BKObjIDField)
	if len(objID) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}
	table := common.GetInstTableName(objID, cts.Kit.TenantID)

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	err = mongodb.Client().Table(table).Delete(cts.Kit.Ctx, filter)
	if err != nil {
		blog.Errorf("delete quoted instances failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	cts.RespEntity(nil)
}
