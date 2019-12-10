/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package process

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateProcessInstanceRelation(ctx core.ContextParams, relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, relation.BizID); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	relation.BizID = bizID

	// validate service category id field
	_, err = p.GetServiceInstance(ctx, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}
	// TODO: asset bizID == category.BizID

	relation.SupplierAccount = ctx.SupplierAccount
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Insert(ctx.Context, &relation); nil != err {
		blog.Errorf("CreateProcessInstanceRelation failed, mongodb failed, table: %s, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, relation, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	// push process instance relation create event
	event := eventclient.NewEventWithHeader(ctx.Header)
	event.EventType = metadata.EventTypeRelation
	event.ObjType = metadata.EventObjTypeProcModule
	event.Action = metadata.EventActionCreate
	event.Data = []metadata.EventData{
		{CurData: relation},
	}
	err = p.eventCli.Push(ctx, event)
	if err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return relation, nil
}

func (p *processOperation) GetProcessInstanceRelation(ctx core.ContextParams, processInstanceID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation := metadata.ProcessInstanceRelation{}

	filter := map[string]int64{
		common.BKProcessIDField: processInstanceID,
	}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).One(ctx.Context, &relation); nil != err {
		blog.Errorf("GetProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, relation: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, relation, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &relation, nil
}

func (p *processOperation) UpdateProcessInstanceRelation(ctx core.ContextParams, processInstanceID int64, input metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation, err := p.GetProcessInstanceRelation(ctx, processInstanceID)
	if err != nil {
		return nil, err
	}

	// TODO: nothing to update currently

	// update fields to local object
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	filter := map[string]int64{"process_id": processInstanceID}

	preData := make([]map[string]interface{}, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).All(ctx.Context, &preData); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, find relation in mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// do update
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(ctx, filter, relation); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, relation, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
	}

	// push process instance relation update event
	var eventArr []*metadata.EventInst
	for _, data := range preData {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeProcModule
		event.Action = metadata.EventActionUpdate
		event.Data = []metadata.EventData{
			{PreData: data, CurData: relation},
		}
		eventArr = append(eventArr, event)
	}
	if err := p.eventCli.Push(ctx, eventArr...); err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return relation, nil
}

func (p *processOperation) ListProcessInstanceRelation(ctx core.ContextParams, option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}

	// filter with matching any sub category
	if option.ServiceInstanceIDs != nil && len(option.ServiceInstanceIDs) > 0 {
		filter[common.BKServiceInstanceIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}

	if option.ProcessTemplateID > 0 {
		filter[common.BKProcessTemplateIDField] = option.ProcessTemplateID
	}

	if option.HostID > 0 {
		filter[common.BKHostIDField] = option.HostID
	}

	if option.ProcessIDs != nil && len(option.ProcessIDs) > 0 {
		processIDFilter := map[string]interface{}{
			common.BKDBIN: option.ProcessIDs,
		}
		filter[common.BKProcIDField] = processIDFilter
	}

	blog.Debug("filter: %v", filter)
	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Start(
		uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).All(ctx.Context, &relations); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleProcessInstanceRelation{
		Count: total,
		Info:  relations,
	}
	return result, nil
}

func (p *processOperation) DeleteProcessInstanceRelation(ctx core.ContextParams, option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder {
	deleteFilter := map[string]interface{}{}
	if option.BusinessID != nil {
		deleteFilter[common.BKAppIDField] = option.BusinessID
	}
	parameterEnough := false
	if option.ProcessIDs != nil {
		parameterEnough = true
		deleteFilter[common.BKProcIDField] = map[string]interface{}{
			common.BKDBIN: option.ProcessIDs,
		}
	}
	if option.ProcessTemplateIDs != nil {
		parameterEnough = true
		deleteFilter[common.BKProcessTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ProcessTemplateIDs,
		}
	}
	if option.ServiceInstanceIDs != nil {
		parameterEnough = true
		deleteFilter[common.BKServiceInstanceIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}
	if option.ModuleIDs != nil {
		parameterEnough = true
		deleteFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}

	if parameterEnough == false {
		blog.Errorf("DeleteProcessInstanceRelation failed, filter parameters not enough, filter: %+v, rid: %s", deleteFilter, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParametersCountNotEnough)
	}

	preData := make([]map[string]interface{}, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(deleteFilter).All(ctx.Context, &preData); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, find relation in mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	// push process instance relation delete event
	var eventArr []*metadata.EventInst
	for _, data := range preData {
		event := eventclient.NewEventWithHeader(ctx.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeProcModule
		event.Action = metadata.EventActionDelete
		event.Data = []metadata.EventData{
			{PreData: data},
		}
		eventArr = append(eventArr, event)
	}
	err := p.eventCli.Push(ctx, eventArr...)
	if err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return nil
}
