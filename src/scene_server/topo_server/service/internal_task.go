/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/logics/settemplate"
)

// SyncModuleTaskHandler sync module under set template by service template task handler
func (s *Service) SyncModuleTaskHandler(ctx *rest.Contexts) {
	// parse task body
	backendWorker := settemplate.BackendWorker{
		ClientSet:       s.Engine.CoreAPI,
		ModuleOperation: s.Logics.ModuleOperation(),
		InstOperation:   s.Logics.InstOperation(),
	}
	task := &metadata.SyncModuleTask{}
	if err := ctx.DecodeInto(task); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := backendWorker.DoModuleSyncTask(ctx.Kit.Header, task.Set, task.ModuleDiff); err != nil {
			blog.ErrorJSON("do module sync task failed, task: %s, err: %s, rid: %s", task, err, ctx.Kit.Rid)
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
