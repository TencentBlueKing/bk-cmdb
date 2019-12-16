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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/settemplate"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *Service) SyncModuleTaskHandler(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// parse task body
	backendWorker := settemplate.BackendWorker{
		ClientSet:       s.Engine.CoreAPI,
		Engine:          s.Engine,
		ObjectOperation: s.Core.ObjectOperation(),
		ModuleOperation: s.Core.ModuleOperation(),
	}
	task := &metadata.SyncModuleTask{}
	if err := data.MarshalJSONInto(task); err != nil {
		blog.ErrorJSON("unmarshal body into task data failed, body: %s, err: %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	err := backendWorker.DoModuleSyncTask(task.Header, task.Set, task.ModuleDiff)
	if err != nil {
		blog.ErrorJSON("DoModuleSyncTask failed, task: %s, err: %s, rid: %s", task, err, params.ReqID)
		return nil, params.Err.CCError(common.CCErrorTopoSyncModuleTaskFailed)
	}
	return nil, nil
}
