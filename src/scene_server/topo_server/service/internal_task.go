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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/settemplate"
)

func (s *Service) SyncModuleTaskHandler(ctx *rest.Contexts) {
	// parse task body
	backendWorker := settemplate.BackendWorker{
		ClientSet:       s.Engine.CoreAPI,
		Engine:          s.Engine,
		ObjectOperation: s.Core.ObjectOperation(),
		ModuleOperation: s.Core.ModuleOperation(),
	}
	task := &metadata.SyncModuleTask{}
	if err := ctx.DecodeInto(task); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := backendWorker.DoModuleSyncTask(ctx.Kit.Header, task.Set, task.ModuleDiff); err != nil {
			blog.ErrorJSON("DoModuleSyncTask failed, task: %s, err: %s, rid: %s", task, err, ctx.Kit.Rid)
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
