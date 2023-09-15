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
	"configcenter/src/common/util"
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

// ListFieldTemplateSyncStatus whether there is a difference between the real-time calculation template and the model,
// the following scenarios are considered to be different:
// 1. attribute conflict. 2. new attribute field 3. unmanagement attribute available 4. update attribute field
// 5. unique check conflict. 6. add a unique check. 7.update unique checksum. 8.unmanage unique checks
func (s *service) ListFieldTemplateSyncStatus(ctx *rest.Contexts) {
	input := new(metadata.ListFieldTmpltSyncStatusOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	res, err := s.logics.FieldTemplateOperation().ListFieldTemplateSyncStatus(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// ListFieldTemplateTasksStatus get the execution status of the
// asynchronous task created by the field combination template.
func (s *service) ListFieldTemplateTasksStatus(ctx *rest.Contexts) {

	input := new(metadata.ListFieldTmplTaskStatusOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := input.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	taskIDs := util.StrArrayUnique(input.TaskIDs)
	cond := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKTaskTypeField: common.SyncFieldTemplateTaskFlag,
			common.BKTaskIDField:   map[string]interface{}{common.BKDBIN: taskIDs},
			common.BKOwnerIDField:  ctx.Kit.SupplierAccount,
		},
		Fields:         []string{common.BKStatusField, common.BKTaskIDField},
		DisableCounter: true,
	}

	taskRes, err := s.clientSet.TaskServer().Task().ListSyncStatusHistory(ctx.Kit.Ctx, ctx.NewHeader(), cond)
	if err != nil {
		blog.Errorf("get task status failed, task ids: %v, err: %v, rid: %s", taskIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(taskRes.Info) != len(taskIDs) {
		blog.Errorf("there is an illegal taskID, task ids:(%v), rid: %s", taskIDs, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "task_ids"))
		return
	}

	result := make([]metadata.ListFieldTmplTaskStatusResult, len(taskIDs))
	for id := range taskRes.Info {
		result[id] = metadata.ListFieldTmplTaskStatusResult{
			TaskID: taskRes.Info[id].TaskID,
			Status: string(taskRes.Info[id].Status),
		}
	}
	ctx.RespEntity(result)
}

// SyncFieldTemplateToObjectTask synchronize field template information to model tasks.
func (s *service) SyncFieldTemplateToObjectTask(ctx *rest.Contexts) {

	syncOption := new(metadata.SyncObjectTask)
	if err := ctx.DecodeInto(syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := syncOption.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objectID, err := s.preCheckExecuteTaskAndGetObjID(ctx.Kit, syncOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.doSyncFieldTemplateTask(ctx.Kit, syncOption, objectID); err != nil {
			blog.Errorf("do sync field template task(%#v) failed, err: %v, rid: %s", syncOption, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// SyncFieldTemplateInfoToObjects synchronize the field combination template information to the corresponding model
func (s *service) SyncFieldTemplateInfoToObjects(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateSyncOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	idPausedMap, err := s.getObjectPausedAttr(ctx.Kit, opt.ObjectIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	tasks := make([]metadata.CreateTaskRequest, 0)
	for _, objID := range opt.ObjectIDs {
		if idPausedMap[objID] {
			blog.Warnf("the model has been paused, no attr synchronization, object id: %d, rid: %s", objID, ctx.Kit.Rid)
			continue
		}

		tasks = append(tasks, metadata.CreateTaskRequest{
			TaskType: common.SyncFieldTemplateTaskFlag,
			InstID:   opt.TemplateID,
			Extra:    objID,
			Data: []interface{}{metadata.SyncObjectTask{
				TemplateID: opt.TemplateID,
				ObjectID:   objID,
			}},
		})
	}

	taskIDs := make([]string, 0)
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		taskRes, err := s.clientSet.TaskServer().Task().CreateFieldTemplateBatch(ctx.Kit.Ctx, ctx.Kit.Header, tasks)
		if err != nil {
			blog.Errorf("create field template sync task(%#v) failed, err: %v, rid: %s", tasks, err, ctx.Kit.Rid)
			return err
		}
		for id := range taskRes {
			taskIDs = append(taskIDs, taskRes[id].TaskID)
		}
		blog.V(4).Infof("successfully created field template sync task: %#v, rid: %s", taskRes, ctx.Kit.Rid)
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(taskIDs)
}
