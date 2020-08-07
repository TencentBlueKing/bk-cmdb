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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (p *processOperation) CreateProcessInstanceRelation(kit *rest.Kit, relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, relation.BizID); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	relation.BizID = bizID

	// validate service category id field
	_, err = p.GetServiceInstance(kit, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}
	// TODO: asset bizID == category.BizID

	relation.SupplierAccount = kit.SupplierAccount
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Insert(kit.Ctx, &relation); nil != err {
		blog.Errorf("CreateProcessInstanceRelation failed, mongodb failed, table: %s, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, relation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	// push process instance relation create event
	event := eventclient.NewEventWithHeader(kit.Header)
	event.EventType = metadata.EventTypeRelation
	event.ObjType = metadata.EventObjTypeProcModule
	event.Action = metadata.EventActionCreate
	event.Data = []metadata.EventData{
		{CurData: relation},
	}
	err = p.eventCli.Push(kit.Ctx, event)
	if err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return relation, nil
}

func (p *processOperation) GetProcessInstanceRelation(kit *rest.Kit, processInstanceID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation := metadata.ProcessInstanceRelation{}

	filter := map[string]int64{
		common.BKProcessIDField: processInstanceID,
	}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).One(kit.Ctx, &relation); nil != err {
		blog.Errorf("GetProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, relation: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, relation, err, kit.Rid)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &relation, nil
}

func (p *processOperation) UpdateProcessInstanceRelation(kit *rest.Kit, processInstanceID int64, input metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation, err := p.GetProcessInstanceRelation(kit, processInstanceID)
	if err != nil {
		return nil, err
	}

	// TODO: nothing to update currently

	// update fields to local object
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	filter := map[string]int64{"process_id": processInstanceID}

	preData := make([]map[string]interface{}, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).All(kit.Ctx, &preData); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, find relation in mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// do update
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(kit.Ctx, filter, relation); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, relation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}

	// push process instance relation update event
	var eventArr []*metadata.EventInst
	for _, data := range preData {
		event := eventclient.NewEventWithHeader(kit.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeProcModule
		event.Action = metadata.EventActionUpdate
		event.Data = []metadata.EventData{
			{PreData: data, CurData: relation},
		}
		eventArr = append(eventArr, event)
	}
	if err := p.eventCli.Push(kit.Ctx, eventArr...); err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return relation, nil
}

func (p *processOperation) ListProcessInstanceRelation(kit *rest.Kit, option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {
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
	if total, err = p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Count(kit.Ctx); nil != err {
		blog.Errorf("ListProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Sort(option.Page.Sort).Start(
		uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).All(kit.Ctx, &relations); nil != err {
		blog.Errorf("ListProcessInstanceRelation failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleProcessInstanceRelation{
		Count: total,
		Info:  relations,
	}
	return result, nil
}

func (p *processOperation) ListHostProcessRelation(kit *rest.Kit, option *metadata.ListProcessInstancesWithHostOption) (*metadata.MultipleHostProcessRelation, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: option.BizID,
	}
	if option.HostIDs != nil && len(option.HostIDs) > 0 {
		filter[common.BKHostIDField] = map[string]interface{}{
			common.BKDBIN: option.HostIDs,
		}
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)
	count, err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("ListHostProcessRelation count mongodb failed, table: %s, filter: %s, err: %s, rid: %s", common.BKTableNameProcessInstanceRelation, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	relations := make([]metadata.HostProcessRelation, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Start(
		uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(option.Page.Sort).All(kit.Ctx, &relations); err != nil {
		blog.ErrorJSON("ListHostProcessRelation select mongodb failed, table: %s, filter: %s, err: %s, rid: %s", common.BKTableNameProcessInstanceRelation, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return &metadata.MultipleHostProcessRelation{
		Count: count,
		Info:  relations,
	}, nil
}

func (p *processOperation) DeleteProcessInstanceRelation(kit *rest.Kit, option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder {
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

	if !parameterEnough {
		blog.Errorf("DeleteProcessInstanceRelation failed, filter parameters not enough, filter: %+v, rid: %s", deleteFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParametersCountNotEnough)
	}

	preData := make([]map[string]interface{}, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(deleteFilter).All(kit.Ctx, &preData); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, find relation in mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Delete(kit.Ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	// push process instance relation delete event
	var eventArr []*metadata.EventInst
	for _, data := range preData {
		event := eventclient.NewEventWithHeader(kit.Header)
		event.EventType = metadata.EventTypeRelation
		event.ObjType = metadata.EventObjTypeProcModule
		event.Action = metadata.EventActionDelete
		event.Data = []metadata.EventData{
			{PreData: data},
		}
		eventArr = append(eventArr, event)
	}
	err := p.eventCli.Push(kit.Ctx, eventArr...)
	if err != nil {
		blog.Errorf("process instance relation event push failed, error:%v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)
	}
	return nil
}
