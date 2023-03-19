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

// UpdateProjectID 更新bk_project_id, 此接口为BCS进行项目数据迁移时的专用接口，其他平台不可使用
func (s *coreService) UpdateProjectID(ctx *rest.Contexts) {
	opt := new(metadata.UpdateProjectIDOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := mapstr.MapStr{
		common.BKFieldID: opt.ID,
	}
	updateData := mapstr.MapStr{
		common.BKProjectIDField: opt.ProjectID,
	}
	err := mongodb.Client().Table(common.BKTableNameBaseProject).Update(ctx.Kit.Ctx, filter, updateData)
	if err != nil {
		blog.Errorf("update project bk_project_id failed, err: %v, filter: %v, updateData: %v, rid: %s", err, filter,
			updateData, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}
