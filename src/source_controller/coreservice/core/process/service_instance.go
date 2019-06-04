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
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateServiceInstance(ctx core.ContextParams, instance metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	// base attribute validate
	if field, err := instance.Validate(); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, instance.Metadata); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// keep metadata clean
	instance.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))

	// validate service template id field
	serviceTemplate, err := p.GetServiceTemplate(ctx, instance.ServiceTemplateID)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, service_template_id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "service_template_id")
	}

	// validate module id field
	if err = p.validateModuleID(ctx, instance.ModuleID); err != nil {
		blog.Errorf("CreateServiceInstance failed, module id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "module_id")
	}

	// validate host id field
	if err = p.validateHostID(ctx, instance.HostID); err != nil {
		blog.Errorf("CreateServiceInstance failed, host id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "host_id")
	}

	// make sure biz id identical with service template
	serviceTemplateBizID, err := metadata.BizIDFromMetadata(serviceTemplate.Metadata)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, parse biz id from service template failed, code: %d, err: %+v, rid: %s", common.CCErrCommInternalServerError, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParseBizIDFromMetadataInDBFailed)
	}
	if bizID != serviceTemplateBizID {
		blog.Errorf("CreateServiceInstance failed, validation failed, input bizID:%d not equal service template bizID:%d, rid: %s", bizID, serviceTemplateBizID, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameProcessTemplate)
	if nil != err {
		blog.Errorf("CreateServiceInstance failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommGenerateRecordIDFailed)
	}
	instance.ID = int64(id)

	instance.Creator = ctx.User
	instance.Modifier = ctx.User
	instance.CreateTime = time.Now()
	instance.LastTime = time.Now()
	instance.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Insert(ctx.Context, &instance); nil != err {
		blog.Errorf("CreateServiceInstance failed, mongodb failed, table: %s, instance: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instance, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBInsertFailed)
	}
	return &instance, nil
}

func (p *processOperation) GetServiceInstance(ctx core.ContextParams, templateID int64) (*metadata.ServiceInstance, error) {
	instance := metadata.ServiceInstance{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).One(ctx.Context, &instance); nil != err {
		blog.Errorf("GetServiceInstance failed, mongodb failed, table: %s, instance: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instance, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	return &instance, nil
}

func (p *processOperation) UpdateServiceInstance(ctx core.ContextParams, instanceID int64, input metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	instance, err := p.GetServiceInstance(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	if field, err := input.Validate(); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// update fields to local object
	// TODO: fixme with update other fields than name
	instance.Name = input.Name

	// do update
	filter := map[string]int64{common.BKFieldID: instanceID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Update(ctx, filter, instance); nil != err {
		blog.Errorf("UpdateServiceTemplate failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBUpdateFailed)
	}
	return instance, nil
}

func (p *processOperation) ListServiceInstance(ctx core.ContextParams, option metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, error) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(option.BusinessID, 10))
	filter := map[string]interface{}{}
	filter[common.MetadataField] = md.ToMapStr()

	if option.ServiceTemplateID != 0 {
		filter[common.BKServiceTemplateIDField] = option.ServiceTemplateID
	}

	if option.HostID != 0 {
		filter[common.BKHostIDField] = option.HostID
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	instances := make([]metadata.ServiceInstance, 0)
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Start(
		uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).All(ctx.Context, &instances); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	if option.WithName == true {
		for idx, instance := range instances {
			instanceName, err := p.GetServiceInstanceName(ctx, instance.ID)
			if err != nil {
				blog.Errorf("ListServiceInstance failed, construct instance name failed, instanceID: %d, err: %+v, rid: %s", instance.ID, err, ctx.ReqID)
				return nil, err
			}
			instances[idx].Name = instanceName
		}
	}

	result := &metadata.MultipleServiceInstance{
		Count: total,
		Info:  instances,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceInstance(ctx core.ContextParams, serviceInstanceID int64) error {
	instance, err := p.GetServiceInstance(ctx, serviceInstanceID)
	if err != nil {
		blog.Errorf("DeleteServiceInstance failed, GetServiceInstance failed, instanceID: %d, err: %+v, rid: %s", serviceInstanceID, err, ctx.ReqID)
		return err
	}

	// service template that referenced by process template shouldn't be removed
	usageFilter := map[string]int64{common.BKServiceInstanceIDField: instance.ID}
	usageCount, err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(usageFilter).Count(ctx.Context)
	if nil != err {
		blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, usageFilter, err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceInstance failed, forbidden delete service instance be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	serviceInstanceFilter := map[string]int64{common.BKFieldID: instance.ID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Delete(ctx, serviceInstanceFilter); nil != err {
		blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, deleteFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, serviceInstanceFilter, err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}

// GetServiceInstanceName get service instance's name, format: `IP + first process name + first process port`
// 可能应用场景：1. 查询服务实例时组装名称；2. 更新进程信息时根据组装名称直接更新到 `name` 字段
// issue: https://github.com/Tencent/bk-cmdb/issues/2485
func (p *processOperation) GetServiceInstanceName(ctx core.ContextParams, instanceID int64) (string, error) {
	instanceName := ""

	// get instance
	instance := metadata.ServiceInstance{}
	instanceFilter := map[string]interface{}{
		common.BKFieldID: instanceID,
	}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(instanceFilter).One(ctx.Context, &instance); err != nil {
		blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instanceFilter, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) == true {
			return "", ctx.Error.Errorf(common.CCErrCommNotFound)
		}
		return "", ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	// get host inner ip
	host := struct {
		InnerIP string `json:"bk_host_innerip" bson:"bk_host_innerip"`
		HostID  int    `json:"bk_host_id" bson:"bk_host_id"`
	}{}

	hostFilter := map[string]interface{}{
		common.BKHostIDField: instance.HostID,
	}
	if err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(hostFilter).One(ctx.Context, &host); err != nil {
		blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameBaseHost, hostFilter, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) == true {
			return "", ctx.Error.Errorf(common.CCErrCommNotFound)
		}
		return "", ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	instanceName += host.InnerIP

	// get first process instance relation
	relation := metadata.ProcessInstanceRelation{}
	relationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: instance.ID,
	}
	order := "id"
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).Sort(order).One(ctx.Context, &relation); err != nil {
		// relation not found means no process in service instance, service instance's name will only contains ip in that case
		if p.dbProxy.IsNotFoundError(err) != true {
			blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, relationFilter, err, ctx.ReqID)
			return "", ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
		}
	}

	if relation.ProcessID != 0 {
		// get process instance
		process := metadata.Process{}
		processFilter := map[string]interface{}{
			common.BKProcIDField: relation.ProcessID,
		}
		if err := p.dbProxy.Table(common.BKTableNameBaseProcess).Find(processFilter).One(ctx.Context, &process); err != nil {
			blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameBaseProcess, processFilter, err, ctx.ReqID)
			if p.dbProxy.IsNotFoundError(err) == true {
				return "", ctx.Error.Errorf(common.CCErrCommNotFound)
			}
			return "", ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
		}

		instanceName += fmt.Sprintf("-%s-%s", process.ProcessName, process.Port)
	}
	return instanceName, nil
}
