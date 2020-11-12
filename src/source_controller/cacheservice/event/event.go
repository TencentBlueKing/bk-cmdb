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

package event

import (
	"context"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/stream"
)

func NewEvent(watch stream.Interface, isMaster discovery.ServiceManageInterface) error {
	e := Event{
		watch:    watch,
		isMaster: isMaster,
	}

	if err := e.runHost(context.Background()); err != nil {
		blog.Errorf("run host event flow failed, err: %v", err)
		return err
	}

	if err := e.runModuleHostRelation(context.Background()); err != nil {
		blog.Errorf("run module host config event flow failed, err: %v", err)
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

	if err := e.runSetTemplate(context.Background()); err != nil {
		blog.Errorf("run set_template event flow failed, err: %v", err)
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

	return nil
}

type Event struct {
	watch    stream.Interface
	isMaster discovery.ServiceManageInterface
}

func (e *Event) runHost(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseHost,
		key:        HostKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runModuleHostRelation(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameModuleHostConfig,
		key:        ModuleHostRelationKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runBiz(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseApp,
		key:        BizKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runSet(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseSet,
		key:        SetKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runModule(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseModule,
		key:        ModuleKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runSetTemplate(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameSetTemplate,
		key:        SetTemplateKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runObjectBase(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseInst,
		key:        ObjectBaseKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runProcess(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameBaseProcess,
		key:        ProcessKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}

func (e *Event) runProcessInstanceRelation(ctx context.Context) error {
	opts := FlowOptions{
		Collection: common.BKTableNameProcessInstanceRelation,
		key:        ProcessInstanceRelationKey,
		watch:      e.watch,
		isMaster:   e.isMaster,
	}

	return newFlow(ctx, opts)
}
