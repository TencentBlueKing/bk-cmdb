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
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateProcessInstanceRelation(ctx core.ContextParams, relation metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, error) {
	// base attribute validate
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.New(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, relation.Metadata); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// keep metadata clean
	relation.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))

	// validate service category id field
	_, err = p.GetServiceInstance(ctx, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}

	// validate service category id field
	_, err = p.GetServiceInstance(ctx, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, service instance id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}
	// TODO: asset bizID == category.Metadata.Label.bk_biz_id

	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Insert(ctx.Context, &relation); nil != err {
		blog.Errorf("CreateProcessInstanceRelation failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, ctx.ReqID)
		return nil, err
	}
	return &relation, nil
}

func (p *processOperation) GetProcessInstanceRelation(ctx core.ContextParams, processInstanceID int64) (*metadata.ProcessInstanceRelation, error) {
	relation := metadata.ProcessInstanceRelation{}

	filter := map[string]int64{"process_id": processInstanceID}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).One(ctx.Context, &relation); nil != err {
		blog.Errorf("GetProcessInstanceRelation failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceTemplate, err, ctx.ReqID)
		if err.Error() == "document not found" {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, err
	}

	return &relation, nil
}

func (p *processOperation) UpdateProcessInstanceRelation(ctx core.ContextParams, processInstanceID int64, input metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, error) {
	relation, err := p.GetProcessInstanceRelation(ctx, processInstanceID)
	if err != nil {
		return nil, err
	}

	// TODO: nothing to update currently

	// update fields to local object
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.New(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// do update
	filter := map[string]int64{"process_id": processInstanceID}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(ctx, filter, relation); nil != err {
		blog.Errorf("UpdateProcessInstanceRelation failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, ctx.ReqID)
		return nil, err
	}
	return relation, nil
}

func (p *processOperation) ListProcessInstanceRelation(ctx core.ContextParams, bizID int64, serviceInstanceID int64, hostID int64, limit metadata.SearchLimit) (*metadata.MultipleProcessInstanceRelation, error) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
	filter := map[string]interface{}{}
	filter["metadata"] = md.ToMapStr()

	// filter with matching any sub category
	if serviceInstanceID > 0 {
		filter["service_instance_id"] = serviceInstanceID
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceTemplate, err, ctx.ReqID)
		return nil, err
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(filter).Start(
		uint64(limit.Offset)).Limit(uint64(limit.Limit)).All(ctx.Context, &relations); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, ctx.ReqID)
		return nil, err
	}

	result := &metadata.MultipleProcessInstanceRelation{
		Count: total,
		Info:  relations,
	}
	return result, nil
}

func (p *processOperation) DeleteProcessInstanceRelation(ctx core.ContextParams, processInstanceID int64) error {
	relation, err := p.GetProcessInstanceRelation(ctx, processInstanceID)
	if err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, GetProcessInstanceRelation failed, templateID: %d, err: %+v, rid: %s", processInstanceID, err, ctx.ReqID)
		return err
	}

	deleteFilter := map[string]int64{"process_id": relation.ProcessID}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessInstanceRelation failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, err, ctx.ReqID)
		return err
	}
	return nil
}
