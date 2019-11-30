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

package taskconfig

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
)

type CodeTaskConfig struct {
	//  task name,
	Name string
	// service name, apiserver, host, topo, proc etc
	SvrType string
	// url path
	Path string
	// http request error. max retry
	Retry int64
}

var (
	// 在代码中配置任务的任务
	codeTaskConfigArr = []CodeTaskConfig{}
)

// init for auto task
func init() {
	AddCodeTaskConfig("sync-settemplate2set", types.CC_MODULE_TOPO, "/topo/v3/internal/task", 1)
}

// AddCodeTaskConfig add task
func AddCodeTaskConfig(name, srvType, path string, retry int64) {
	blog.Infof("add task. name:%s, service type:%s, path:%s", name, srvType, path)
	codeTaskConfigArr = append(codeTaskConfigArr, CodeTaskConfig{
		Name:    name,
		SvrType: srvType,
		Path:    path,
		Retry:   retry,
	})
}

// GetCodeTaskConfig return code  task config
func GetCodeTaskConfig() []CodeTaskConfig {
	return codeTaskConfigArr
}
