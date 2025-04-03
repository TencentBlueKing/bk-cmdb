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

// Package identifier TODO
package identifier

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/task"
)

// Identity is the host identity event flow struct
type Identity struct {
	tasks []*task.Task
}

// GetWatchTasks returns the event watch tasks
func (i *Identity) GetWatchTasks() []*task.Task {
	return i.tasks
}

// NewIdentity new host identifier event watch
func NewIdentity() (*Identity, error) {
	identity := &Identity{tasks: make([]*task.Task, 0)}

	base := identityOptions{}

	host := base
	host.key = event.HostKey
	host.watchFields = needCareHostFields
	if err := identity.addWatchTask(host); err != nil {
		blog.Errorf("new host identify host event failed, err: %v", err)
		return nil, err
	}
	blog.Info("host identity events, watch host success.")

	relation := base
	relation.key = event.ModuleHostRelationKey
	relation.watchFields = []string{common.BKHostIDField}
	if err := identity.addWatchTask(relation); err != nil {
		blog.Errorf("new host identify host relation event failed, err: %v", err)
		return nil, err
	}
	blog.Info("host identity events, watch host relation success.")

	process := base
	process.key = event.ProcessKey
	process.watchFields = []string{common.BKProcessIDField}
	if err := identity.addWatchTask(process); err != nil {
		blog.Errorf("new host identify process event failed, err: %v", err)
		return nil, err
	}
	blog.Info("host identity events, watch process success.")

	procRel := base
	procRel.key = event.ProcessInstanceRelationKey
	procRel.watchFields = []string{common.BKHostIDField}
	if err := identity.addWatchTask(procRel); err != nil {
		blog.Errorf("new host identify process relation event failed, err: %v", err)
		return nil, err
	}
	blog.Info("host identity events, watch process relation success.")

	return identity, nil
}
