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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateDynamicGroup creates a new dynamic group object.
func (s *coreService) CreateDynamicGroup(ctx *rest.Contexts) {
	newDynamicGroup := meta.DynamicGroup{}
	if err := ctx.DecodeInto(&newDynamicGroup); err != nil {
		blog.Errorf("create dynamic group failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	filter := common.KvMap{common.BKFieldName: newDynamicGroup.Name, common.BKAppIDField: newDynamicGroup.AppID}
	rowCount, err := mongodb.Client().Table(common.BKTableNameDynamicGroup).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("create dynamic group failed, query count err: %+v, filter: %v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if rowCount != 0 {
		blog.Errorf("create dynamic group failed, dynamic group[%s] already exist, rid: %s", newDynamicGroup.Name, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "name"))
		return
	}

	// gen new dynamic group ID.
	newDynamicGroupID, err := meta.NewDynamicGroupID()
	if err != nil {
		blog.Errorf("create dynamic group failed, gen new ID, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCSystemUnknownError))
		return
	}

	newDynamicGroup.ID = newDynamicGroupID
	newDynamicGroup.ModifyUser = util.GetUser(ctx.Kit.Header)
	newDynamicGroup.CreateTime = time.Now().UTC()
	newDynamicGroup.UpdateTime = newDynamicGroup.CreateTime

	err = mongodb.Client().Table(common.BKTableNameDynamicGroup).Insert(ctx.Kit.Ctx, newDynamicGroup)
	if err != nil {
		blog.Errorf("create dynamic group failed, group: %+v err: %+v, rid: %s", newDynamicGroup, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}
	ctx.RespEntity(meta.ID{ID: newDynamicGroupID})
}

// UpdateDynamicGroup updates target dynamic group.
func (s *coreService) UpdateDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	bizIDUint64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Error("update dynamic group failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	data["modify_user"] = util.GetUser(ctx.Kit.Header)
	data[common.LastTimeField] = time.Now().UTC()

	filter := common.KvMap{common.BKAppIDField: bizIDUint64, common.BKFieldID: targetID}
	err = mongodb.Client().Table(common.BKTableNameDynamicGroup).Update(ctx.Kit.Ctx, filter, data)
	if err != nil {
		blog.Errorf("update dynamic group failed, err: %+v, ctx: %v, rid: %s", err, ctx, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}
	ctx.RespEntity(nil)
}

// DeleteDynamicGroup deletes target dynamic group.
func (s *coreService) DeleteDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	bizIDUint64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Error("delete dynamic group failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	filter := common.KvMap{common.BKFieldID: targetID, common.BKAppIDField: bizIDUint64}
	rowCount, err := mongodb.Client().Table(common.BKTableNameDynamicGroup).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("delete dynamic group failed, err: %+v, ctx: %v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if rowCount != 1 {
		blog.Errorf("delete dynamic group failed, not permissions or not exists, ctx: %v, rid: %s", filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}

	err = mongodb.Client().Table(common.BKTableNameDynamicGroup).Delete(ctx.Kit.Ctx, filter)
	if err != nil {
		blog.Errorf("delete dynamic group failed, err: %+v, ctx: %v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}
	ctx.RespEntity(nil)
}

// GetDynamicGroup returns target dynamic group detail.
func (s *coreService) GetDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	//  target dynamic group ID.
	targetID := req.PathParameter("id")

	bizIDUint64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("get dynamic group failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	filter := common.KvMap{common.BKFieldID: targetID, common.BKAppIDField: bizIDUint64}

	result := &meta.DynamicGroup{}
	err = mongodb.Client().Table(common.BKTableNameDynamicGroup).Find(filter).One(ctx.Kit.Ctx, result)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("get dynamic group failed, ID: %s, err: %+v, rid: %s", targetID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if len(result.Name) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommNotFound))
		return
	}
	ctx.RespEntity(result)
}

// SearchDynamicGroup returns dynamic group list with target conditions.
func (s *coreService) SearchDynamicGroup(ctx *rest.Contexts) {
	input := new(meta.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		blog.Errorf("search dynamic groups failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	var condition map[string]interface{}
	if input.Condition != nil {
		condition = input.Condition
	} else {
		condition = make(map[string]interface{})
	}

	start, limit, sort := input.Page.Start, input.Page.Limit, input.Page.Sort
	if len(sort) == 0 {
		sort = common.CreateTimeField
	}

	var finalCount uint64

	if !input.DisableCounter {
		count, err := mongodb.Client().Table(common.BKTableNameDynamicGroup).Find(condition).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("search dynamic groups failed, can't open counter, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}
		finalCount = count
	}

	result := []meta.DynamicGroup{}

	if err := mongodb.Client().Table(common.BKTableNameDynamicGroup).Find(condition).Fields(input.Fields...).Sort(sort).
		Start(uint64(start)).Limit(uint64(limit)).All(ctx.Kit.Ctx, &result); err != nil {

		blog.Errorf("search dynamic groups failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(meta.DynamicGroupBatch{
		Count: finalCount,
		Info:  result,
	})
}
