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

package flow

import (
	"context"

	"configcenter/pkg/tenant"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/task"
)

// NewEvent new event flow
func NewEvent() (*Event, error) {
	e := &Event{
		tasks: make([]*task.Task, 0),
	}

	if err := e.addHostTask(); err != nil {
		blog.Errorf("add host event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addModuleHostRelationTask(); err != nil {
		blog.Errorf("add module host config event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addBizSetTask(); err != nil {
		blog.Errorf("add biz set event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addBizTask(); err != nil {
		blog.Errorf("add biz event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addSetTask(); err != nil {
		blog.Errorf("add set event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addModuleTask(); err != nil {
		blog.Errorf("add module event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addObjectBaseTask(context.Background()); err != nil {
		blog.Errorf("add object base event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addProcessTask(); err != nil {
		blog.Errorf("add process event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addProcessInstanceRelationTask(); err != nil {
		blog.Errorf("add process instance relation event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addInstAsstTask(); err != nil {
		blog.Errorf("add instance association event flow task failed, err: %v", err)
		return nil, err
	}

	if err := e.addPlatTask(); err != nil {
		blog.Errorf("add plat event flow task failed, err: %v", err)
	}

	if err := e.addProjectTask(); err != nil {
		blog.Errorf("add project event flow task failed, err: %v", err)
	}

	return e, nil
}

// Event is the event flow struct
type Event struct {
	tasks []*task.Task
}

// GetWatchTasks returns the event flow tasks
func (e *Event) GetWatchTasks() []*task.Task {
	return e.tasks
}

func (e *Event) addHostTask() error {
	opts := flowOptions{
		key:         event.HostKey,
		EventStruct: new(metadata.HostMapStr),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addModuleHostRelationTask() error {
	opts := flowOptions{
		key:         event.ModuleHostRelationKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addBizTask() error {
	opts := flowOptions{
		key:         event.BizKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addSetTask() error {
	opts := flowOptions{
		key:         event.SetKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addModuleTask() error {
	opts := flowOptions{
		key:         event.ModuleKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addObjectBaseTask(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ObjectBaseKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addInstanceFlowTask(ctx, opts, parseEvent)
}

func (e *Event) addProcessTask() error {
	opts := flowOptions{
		key:         event.ProcessKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addProcessInstanceRelationTask() error {
	opts := flowOptions{
		key:         event.ProcessInstanceRelationKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addInstAsstTask() error {
	opts := flowOptions{
		key:         event.InstAsstKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addInstAsstFlowTask(opts, parseInstAsstEvent)
}

func (e *Event) addBizSetTask() error {
	opts := flowOptions{
		key:         event.BizSetKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addPlatTask() error {
	opts := flowOptions{
		key:         event.PlatKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addProjectTask() error {
	opts := flowOptions{
		key:         event.ProjectKey,
		EventStruct: new(map[string]interface{}),
	}

	return e.addFlowTask(opts, parseEvent)
}

func (e *Event) addFlowTask(opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}

	flowTask, err := flow.GenWatchTask()
	if err != nil {
		return err
	}

	e.tasks = append(e.tasks, flowTask)
	return nil
}

func (e *Event) addInstanceFlowTask(ctx context.Context, opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}
	instFlow := InstanceFlow{
		Flow: flow,
		mainlineObjectMap: &mainlineObjectMap{
			data: make(map[string]map[string]struct{}),
		},
	}

	err = tenant.ExecForAllTenants(func(tenantID string) error {
		mainlineObjMap, err := instFlow.getMainlineObjectMap(ctx, tenantID)
		if err != nil {
			blog.Errorf("run object instance watch, but get tenant %s mainline objects failed, err: %v", tenantID, err)
			return err
		}
		instFlow.mainlineObjectMap.Set(tenantID, mainlineObjMap)

		go instFlow.syncMainlineObjectMap(tenantID)
		return nil
	})
	if err != nil {
		return err
	}

	flowTask, err := instFlow.GenWatchTask()
	if err != nil {
		return err
	}

	e.tasks = append(e.tasks, flowTask)
	return nil
}

func (e *Event) addInstAsstFlowTask(opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}
	instAsstFlow := InstAsstFlow{
		Flow: flow,
	}

	flowTask, err := instAsstFlow.GenWatchTask()
	if err != nil {
		return err
	}

	e.tasks = append(e.tasks, flowTask)
	return nil
}
