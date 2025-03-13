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

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/task"
)

// NewEvent new event flow
func NewEvent(watchTask *task.Task) error {
	e := Event{
		task: watchTask,
	}

	if err := e.runHost(context.Background()); err != nil {
		blog.Errorf("run host event flow failed, err: %v", err)
		return err
	}

	if err := e.runModuleHostRelation(context.Background()); err != nil {
		blog.Errorf("run module host config event flow failed, err: %v", err)
		return err
	}

	if err := e.runBizSet(context.Background()); err != nil {
		blog.Errorf("run biz set event flow failed, err: %v", err)
		return err
	}

	if err := e.runBiz(context.Background()); err != nil {
		blog.Errorf("run biz event flow failed, err: %v", err)
		return err
	}

	if err := e.runSet(context.Background()); err != nil {
		blog.Errorf("run set event flow failed, err: %v", err)
		return err
	}

	if err := e.runModule(context.Background()); err != nil {
		blog.Errorf("run module event flow failed, err: %v", err)
		return err
	}

	if err := e.runObjectBase(context.Background()); err != nil {
		blog.Errorf("run object base event flow failed, err: %v", err)
		return err
	}

	if err := e.runProcess(context.Background()); err != nil {
		blog.Errorf("run process event flow failed, err: %v", err)
		return err
	}

	if err := e.runProcessInstanceRelation(context.Background()); err != nil {
		blog.Errorf("run process instance relation event flow failed, err: %v", err)
		return err
	}

	if err := e.runInstAsst(context.Background()); err != nil {
		blog.Errorf("run instance association event flow failed, err: %v", err)
		return err
	}

	if err := e.runPlat(context.Background()); err != nil {
		blog.Errorf("run plat event flow failed, err: %v", err)
	}

	if err := e.runProject(context.Background()); err != nil {
		blog.Errorf("run project event flow failed, err: %v", err)
	}

	return nil
}

// Event is the event flow struct
type Event struct {
	task *task.Task
}

func (e *Event) runHost(ctx context.Context) error {
	opts := flowOptions{
		key:         event.HostKey,
		task:        e.task,
		EventStruct: new(metadata.HostMapStr),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runModuleHostRelation(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ModuleHostRelationKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runBiz(ctx context.Context) error {
	opts := flowOptions{
		key:         event.BizKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runSet(ctx context.Context) error {
	opts := flowOptions{
		key:         event.SetKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runModule(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ModuleKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runObjectBase(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ObjectBaseKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newInstanceFlow(ctx, opts, parseEvent)
}

func (e *Event) runProcess(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ProcessKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runProcessInstanceRelation(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ProcessInstanceRelationKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runInstAsst(ctx context.Context) error {
	opts := flowOptions{
		key:         event.InstAsstKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newInstAsstFlow(ctx, opts, parseInstAsstEvent)
}

func (e *Event) runBizSet(ctx context.Context) error {
	opts := flowOptions{
		key:         event.BizSetKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runPlat(ctx context.Context) error {
	opts := flowOptions{
		key:         event.PlatKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}

func (e *Event) runProject(ctx context.Context) error {
	opts := flowOptions{
		key:         event.ProjectKey,
		task:        e.task,
		EventStruct: new(map[string]interface{}),
	}

	return newFlow(ctx, opts, parseEvent)
}
