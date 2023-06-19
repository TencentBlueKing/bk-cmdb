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

package fieldtmpl

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// ListFieldTemplateModelStatus query the status of the model in the template
func (s *service) ListFieldTemplateModelStatus(ctx *rest.Contexts) {
	input := new(metadata.ListFieldTmplModelStatusOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// 1. 查询模型和模版的差异，得到「是否需要进行同步」的结果
	syncStatusOpt := &metadata.ListFieldTmpltSyncStatusOption{ID: input.ID, ObjectIDs: input.ObjectIDs}
	syncStatusRes, err := s.logics.FieldTemplateOperation().ListFieldTemplateSyncStatus(ctx.Kit, syncStatusOpt)
	if err != nil {
		blog.Errorf("list field template sync status failed, data: %v, err: %v, rid: %s", syncStatusOpt, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(syncStatusRes) != len(input.ObjectIDs) {
		blog.Errorf("can not find all object sync status, opt: %v, resp: %v, rid: %s", syncStatusOpt, syncStatusRes,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "object_ids"))
		return
	}

	// 2. 查询模版绑定的模型的「同步任务状态结果」
	taskSyncOpt := &metadata.ListFieldTmplSyncTaskStatusOption{ID: input.ID, ObjectIDs: input.ObjectIDs}
	taskSyncRes, err := s.clientSet.TaskServer().Task().ListFieldTemplateTaskSyncResult(ctx.Kit.Ctx, ctx.Kit.Header,
		taskSyncOpt)
	if err != nil {
		blog.Errorf("list field template task sync result failed, opt: %+v, err: %v, rid: %s", taskSyncOpt, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	taskMap := make(map[int64]metadata.ListFieldTmplTaskSyncResult)
	for _, task := range taskSyncRes {
		taskMap[task.ObjectID] = task
	}

	// 3. 根据1和2的结果，得到模型最终返回的状态
	res := make([]metadata.ListFieldTmplTaskSyncResult, len(input.ObjectIDs))
	for idx, status := range syncStatusRes {
		task, ok := taskMap[status.ObjectID]
		if !ok {
			task = metadata.ListFieldTmplTaskSyncResult{ObjectID: status.ObjectID}
		}

		// 模版和模型没有差异 == 同步完成
		if !status.NeedSync {
			task.Status = metadata.APITaskStatusSuccess
			task.FailMsg = ""
			res[idx] = task
			continue
		}

		// 模版和模型有差异 + 查不到模版绑定的模型的任务状态 = 需要同步
		// 模版和模型有差异 + 能查到模版绑定的模型的任务状态 + 状态为同步完成 = 需要同步
		if !ok || task.Status.IsSuccessful() {
			task.Status = metadata.APITAskStatusNeedSync
			res[idx] = task
			continue
		}

		// 模版和模型有差异 + 能查到模版绑定的模型的任务状态 + 任务状态为非同步完成的其他状态 = 当前任务的状态
		res[idx] = task
	}

	ctx.RespEntity(res)
}
