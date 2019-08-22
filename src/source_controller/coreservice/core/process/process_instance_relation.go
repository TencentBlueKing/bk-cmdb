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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
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
	if bizID, err = p.validateBizID(ctx, relation.Metadata); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// keep metadata clean
	relation.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))

	// validate service category id field
	_, err = p.GetServiceInstance(ctx, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}

	// validate service category id field
	_, err = p.GetServiceInstance(ctx, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}
	// TODO: asset bizID == category.Metadata.Label.bk_biz_id

	relation.SupplierAccount = ctx.SupplierAccount
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Insert(ctx.Context, &relation); nil != err {
		blog.Errorf("CreateProcessInstanceRelation failed, mongodb failed, table: %s, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, relation, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
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

	// do update
	filter := map[string]int64{"process_id": processInstanceID}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(ctx, filter, relation); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, relation: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, filter, relation, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return relation, nil
}

func (p *processOperation) ListProcessInstanceRelation(ctx core.ContextParams, option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(option.BusinessID, 10))
	filter := map[string]interface{}{}
	filter[common.MetadataField] = md.ToMapStr()

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
		blog.Errorf("DeleteProcessInstanceRelation failed, filter parameters not enough, filter: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParametersCountNotEnough)
	}

	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, deleteFilter, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
