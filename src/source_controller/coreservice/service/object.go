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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// SearchUUIDByObj search object uuid by object id
func (s *coreService) SearchUUIDByObj(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter("bk_obj_id")
	if len(objID) == 0 {
		blog.Errorf("object id is empty, rid: %s")
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	result := metadata.Object{}
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(mapstr.MapStr{"bk_obj_id": objID}).
		Fields(metadata.ModelFieldObjUUID).One(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get object uuid failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.UUID)
	return
}

// SearchObjByUUID search object by uuid
func (s *coreService) SearchObjByUUID(ctx *rest.Contexts) {

	uuid := ctx.Request.PathParameter("uuid")
	if len(uuid) == 0 {
		blog.Errorf("uuid is empty, rid: %s")
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjUUID))
		return
	}

	result := metadata.Object{}
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(mapstr.MapStr{
		metadata.ModelFieldObjUUID: uuid}).Fields("bk_obj_id").One(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get object failed, err: %v， rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.ObjectID)
	return
}
