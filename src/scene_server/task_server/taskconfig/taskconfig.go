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

// Package taskconfig TODO
package taskconfig

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
)

// CodeTaskConfig task queue config defined in code
type CodeTaskConfig struct {
	// Name name of the task queue
	Name string
	// SvrType service type, api, host, topo, proc etc.
	SvrType string
	// Path url path
	Path string
	// Retry max retry attempts if request the callback interface failed
	Retry int64
	// LockTTL the expire time of task lock in minutes
	LockTTL int64
}

var (
	// 在代码中配置任务的任务
	codeTaskConfigArr = make([]CodeTaskConfig, 0)
)

// init task queue
func init() {
	AddCodeTaskConfig(common.SyncSetTaskFlag, types.CC_MODULE_TOPO, "/topo/v3/internal/sync/module/task", 1, 2)
	AddCodeTaskConfig(common.SyncModuleTaskFlag, types.CC_MODULE_PROC, "/process/v3/sync/service_instance/task", 1, 2)
	AddCodeTaskConfig(common.SyncModuleHostApplyTaskFlag, types.CC_MODULE_HOST,
		"/host/v3/updatemany/module/host_apply_plan/task", 1, 2)
	AddCodeTaskConfig(common.SyncServiceTemplateHostApplyTaskFlag, types.CC_MODULE_PROC,
		"/process/v3/updatemany/service_template/host_apply_plan/task", 1, 2)
}

// AddCodeTaskConfig add task
func AddCodeTaskConfig(name, srvType, path string, retry, lockTTL int64) {
	blog.Infof("add task. name: %s, service type: %s, path: %s", name, srvType, path)
	codeTaskConfigArr = append(codeTaskConfigArr, CodeTaskConfig{
		Name:    name,
		SvrType: srvType,
		Path:    path,
		Retry:   retry,
		LockTTL: lockTTL,
	})
}

// GetCodeTaskConfig return code  task config
func GetCodeTaskConfig() []CodeTaskConfig {
	return codeTaskConfigArr
}
