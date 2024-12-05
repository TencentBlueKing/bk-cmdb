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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// CreateProcessInstanceRelation create process instance relation
func (p *processOperation) CreateProcessInstanceRelation(kit *rest.Kit,
	relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {

	if err := p.validateRelation(kit, relation); err != nil {
		return nil, err
	}
	relation.TenantID = kit.TenantID
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Insert(kit.Ctx,
		relation); err != nil {
		blog.Errorf("insert process instance relation failed, table: %s, relation: %+v, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, relation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return relation, nil
}

// CreateProcessInstanceRelations create process instance relations
func (p *processOperation) CreateProcessInstanceRelations(kit *rest.Kit,
	relations []*metadata.ProcessInstanceRelation) ([]*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {

	for _, relation := range relations {
		if err := p.validateRelation(kit, relation); err != nil {
			return nil, err
		}
		relation.TenantID = kit.TenantID
	}

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Insert(kit.Ctx,
		relations); err != nil {
		blog.Errorf("insert process instance relations failed, table: %s, relations: %+v, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, relations, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return relations, nil
}

func (p *processOperation) validateRelation(kit *rest.Kit,
	relation *metadata.ProcessInstanceRelation) errors.CCErrorCoder {

	// base attribute validate
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("process instance relation validation failed, code: %d, err: %v, rid: %s",
			common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, relation.BizID); err != nil {
		blog.Errorf("bizID validation failed, code: %d, err: %v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	relation.BizID = bizID

	// validate service category id field
	_, err = p.GetServiceInstance(kit, relation.ServiceInstanceID)
	if err != nil {
		blog.Errorf("get service instance by id failed, code: %d, err: %v, rid: %s", common.CCErrCommParamsInvalid, err,
			kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_id")
	}
	// TODO: asset bizID == category.BizID

	return nil
}

// GetProcessInstanceRelation get process instance relation by process instance id
func (p *processOperation) GetProcessInstanceRelation(kit *rest.Kit,
	processInstanceID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation := metadata.ProcessInstanceRelation{}

	filter := map[string]int64{
		common.BKProcessIDField: processInstanceID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Find(filter).One(kit.Ctx,
		&relation); err != nil {
		blog.Errorf("get process instance relation failed, table: %s, filter: %+v, relation: %+v, err: %v, rid: %s",
			common.BKTableNameServiceTemplate, filter, relation, err, kit.Rid)
		if mongodb.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &relation, nil
}

// UpdateProcessInstanceRelation update process instance relation
func (p *processOperation) UpdateProcessInstanceRelation(kit *rest.Kit, processInstanceID int64,
	input metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	relation, err := p.GetProcessInstanceRelation(kit, processInstanceID)
	if err != nil {
		return nil, err
	}

	// TODO: nothing to update currently

	// update fields to local object
	if field, err := relation.Validate(); err != nil {
		blog.Errorf("process instance relation valid failed, code: %d, err: %v, rid: %s",
			common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	filter := map[string]int64{"process_id": processInstanceID}

	// do update
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Update(kit.Ctx, filter,
		relation); err != nil {
		blog.Errorf("update process instance relation failed, table: %s, filter: %+v, relation: %+v, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, filter, relation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return relation, nil
}

// ListProcessInstanceRelation list process instance relation
func (p *processOperation) ListProcessInstanceRelation(kit *rest.Kit,
	option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation,
	errors.CCErrorCoder) {

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

	total, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Find(filter).
		Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count process instance relation failed, table: %s, filter: %+v, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Find(filter).
		Sort(option.Page.Sort).Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).All(kit.Ctx,
		&relations); err != nil {
		blog.Errorf("list process instance relation failed, table: %s, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleProcessInstanceRelation{
		Count: total,
		Info:  relations,
	}
	return result, nil
}

// DeleteProcessInstanceRelation delete process instance relation
func (p *processOperation) DeleteProcessInstanceRelation(kit *rest.Kit,
	option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder {

	deleteFilter := map[string]interface{}{}
	if option.BusinessID != nil {
		deleteFilter[common.BKAppIDField] = *option.BusinessID
	}

	parameterEnough := false
	if len(option.ProcessIDs) != 0 {
		parameterEnough = true
		deleteFilter[common.BKProcIDField] = map[string]interface{}{
			common.BKDBIN: option.ProcessIDs,
		}
	}

	if len(option.ProcessTemplateIDs) != 0 {
		parameterEnough = true
		deleteFilter[common.BKProcessTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ProcessTemplateIDs,
		}
	}

	if len(option.ServiceInstanceIDs) != 0 {
		parameterEnough = true
		deleteFilter[common.BKServiceInstanceIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}

	if len(option.ModuleIDs) != 0 {
		parameterEnough = true
		deleteFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}

	if !parameterEnough {
		blog.Errorf("filter parameters not enough, filter: %+v, rid: %s", deleteFilter, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParametersCountNotEnough)
	}

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Delete(kit.Ctx,
		deleteFilter); err != nil {
		blog.Errorf("delete process instance relation failed, table: %s, filter: %+v, err: %v, rid: %s",
			common.BKTableNameProcessInstanceRelation, deleteFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
