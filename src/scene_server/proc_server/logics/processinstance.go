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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logic) ListProcessInstances(kit *rest.Kit, bizID int64, serviceInstanceID int64, fields []string) (
	[]metadata.ProcessInstance, errors.CCErrorCoder) {

	if serviceInstanceID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}
	// list process instance relation
	relationOption := metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: []int64{serviceInstanceID},
	}
	relationsResult, err := lgc.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, &relationOption)
	if err != nil {
		return nil, kit.CCError.CCErrorf(common.CCErrProcGetServiceInstancesFailed, "list process instance "+
			"relation failed, bizID: %d, serviceInstanceID: %d, err: %+v", bizID, serviceInstanceID, err)
	}

	processIDs := make([]int64, 0)
	for _, relation := range relationsResult.Info {
		processIDs = append(processIDs, relation.ProcessID)
	}
	filter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: processIDs,
		},
	}
	reqParam := &metadata.QueryCondition{
		Condition: filter,
		Fields:    fields,
	}
	processResult, ccErr := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != ccErr {
		return nil, kit.CCError.CCErrorf(common.CCErrProcGetServiceInstancesFailed, "list process instance "+
			"property failed, bizID: %d, processIDs: %+v, err: %+v", bizID, processIDs, ccErr)
	}

	processIDPropertyMap := map[int64]mapstr.MapStr{}
	for _, process := range processResult.Data.Info {
		processIDVal, exist := process.Get(common.BKProcessIDField)
		if !exist {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParseDataFailed, "list process instance failed, parse "+
				"bk_process_id from process property failed, field not exist, bizID: %d, processIDs: %+v", bizID, processIDs)
		}
		processID, err := util.GetInt64ByInterface(processIDVal)
		if err != nil {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParseDataFailed, "list process instance failed, "+
				"parse bk_process_id from process property failed, parse field to int64 failed, bizID: %d, "+
				"processIDs: %+v, process: %+v, err: %+v", bizID, processIDs, process, err)
		}
		processIDPropertyMap[processID] = process
	}

	processInstanceList := make([]metadata.ProcessInstance, 0)
	for _, relation := range relationsResult.Info {
		processInstance := metadata.ProcessInstance{
			Property: nil,
			Relation: relation,
		}
		process, exist := processIDPropertyMap[relation.ProcessID]
		if exist {
			processInstance.Property = process
		}
		processInstanceList = append(processInstanceList, processInstance)
	}

	return processInstanceList, nil
}

func (lgc *Logic) ListProcessInstanceWithIDs(kit *rest.Kit, procIDs []int64) ([]metadata.Process, errors.CCErrorCoder) {
	reqParam := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				common.BKDBIN: procIDs,
			},
		}),
	}
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("rid: %s list process instance with procID: %d failed, err: %v", kit.Rid, procIDs, err)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !ret.Result {
		blog.Errorf("rid: %s list process instance with procID: %d failed, err: %v", kit.Rid, procIDs, ret.ErrMsg)
		return nil, ret.CCError()
	}

	processes := make([]metadata.Process, 0)
	for _, p := range ret.Data.Info {
		process := new(metadata.Process)
		if err := p.MarshalJSONInto(process); err != nil {
			return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		processes = append(processes, *process)
	}

	return processes, nil
}

func (lgc *Logic) GetProcessInstanceWithID(kit *rest.Kit, procID int64) (*metadata.Process, errors.CCErrorCoder) {
	condition := map[string]interface{}{
		common.BKProcessIDField: procID,
	}

	reqParam := new(metadata.QueryCondition)
	reqParam.Condition = condition
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("GetProcessInstanceWithID failed, get process instance with procID: %d failed, err: %v, rid: %s", procID, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !ret.Result {
		blog.Errorf("GetProcessInstanceWithID failed, get process instance with procID: %d failed, err: %v, rid: %s", procID, ret.ErrMsg, kit.Rid)
		return nil, ret.CCError()
	}

	process := new(metadata.Process)
	if len(ret.Data.Info) == 0 {
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if err := ret.Data.Info[0].MarshalJSONInto(process); err != nil {
		blog.Errorf("GetProcessInstanceWithID failed err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return process, nil
}

func (lgc *Logic) UpdateProcessInstance(kit *rest.Kit, procID int64, info mapstr.MapStr) errors.CCErrorCoder {
	delete(info, common.BkSupplierAccount)
	option := metadata.UpdateOption{
		Data: info,
		Condition: map[string]interface{}{
			common.BKProcessIDField: procID,
		},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		blog.ErrorJSON("UpdateProcessInstance failed, UpdateInstance http request failed, option: %s, err: %s, rid: %s", option, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.ErrorJSON("UpdateProcessInstance failed, UpdateInstance failed, option: %s, response: %s, rid: %s", option, result, kit.Rid)
		return errors.New(result.Code, result.ErrMsg)
	}
	return nil
}

func (lgc *Logic) DeleteProcessInstance(kit *rest.Kit, procID int64) errors.CCErrorCoder {
	rid := kit.Rid
	option := metadata.DeleteOption{
		Condition: map[string]interface{}{
			common.BKProcessIDField: procID,
		},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		blog.ErrorJSON("DeleteProcessInstance failed, DeleteInstance failed, option: %s, err: %s, rid: %s", option, err.Error(), rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("rid: %s, delete process instance: %d failed, err: %s", kit.Rid, procID, result.ErrMsg)
		return errors.New(result.Code, result.ErrMsg)
	}

	return nil
}

func (lgc *Logic) DeleteProcessInstanceBatch(kit *rest.Kit, procIDs []int64) errors.CCErrorCoder {
	if procIDs == nil {
		return nil
	}
	option := metadata.DeleteOption{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				common.BKDBIN: procIDs,
			},
		}),
	}
	result, err := lgc.CoreAPI.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("rid: %s, delete process instance: %d failed, err: %s", kit.Rid, procIDs, result.ErrMsg)
		return result.CCError()
	}

	return nil
}

func (lgc *Logic) CreateProcessInstance(kit *rest.Kit, processData map[string]interface{}) (int64, errors.CCErrorCoder) {
	inputParam := metadata.CreateModelInstance{
		Data: processData,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstance failed, http request failed, err: %+v, rid: %s", err, kit.Rid)
		return 0, errors.CCHttpError
	}

	if !result.Result {
		blog.Errorf("rid: %s, create process instance: %+v failed, err: %s", kit.Rid, processData, result.ErrMsg)
		return 0, errors.New(result.Code, result.ErrMsg)
	}

	return int64(result.Data.Created.ID), nil
}

func (lgc *Logic) CreateProcessInstances(kit *rest.Kit, processDatas []map[string]interface{}) ([]int64, errors.CCErrorCoder) {

	data := make([]mapstr.MapStr, len(processDatas))
	for idx := range processDatas {
		data[idx] = processDatas[idx]
	}

	inputParam := metadata.CreateManyModelInstance{
		Datas: data,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstances failed, http request failed, err: %+v, inputParam:%#v, rid: %s", err, inputParam, kit.Rid)
		return nil, errors.CCHttpError
	}
	if !result.Result {
		blog.Errorf("CreateProcessInstances failed, http request failed, err: %+v, inputParam:%#v, rid: %s", result.ErrMsg, inputParam, kit.Rid)
		return nil, errors.New(result.Code, result.ErrMsg)
	}

	if len(processDatas) != len(result.Data.Created) {
		blog.Errorf("CreateProcessInstances failed, len(processes) != len(result.Created), inputParam: %#v, rid: %s", inputParam, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrProcCreateProcessFailed)
	}

	processIDs := make([]int64, len(processDatas))
	for idx, created := range result.Data.Created {
		processIDs[idx] = int64(created.ID)
	}

	return processIDs, nil
}

// it works to find the different attribute value between the process instance and it's bounded process template.
// return with the changed attribute's details.
func (lgc *Logic) DiffWithProcessTemplate(t *metadata.ProcessProperty, i *metadata.Process, host map[string]interface{},
	attrMap map[string]metadata.Attribute) ([]metadata.ProcessChangedAttribute, error) {

	changes := make([]metadata.ProcessChangedAttribute, 0)
	if t == nil || i == nil {
		return changes, nil
	}

	if metadata.IsAsDefaultValue(t.ProcNum.AsDefaultValue) {
		if (t.ProcNum.Value == nil && i.ProcNum != nil) || (t.ProcNum.Value != nil && i.ProcNum == nil) || (t.ProcNum.Value != nil && *t.ProcNum.Value != *i.ProcNum) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["proc_num"].ID,
				PropertyID:            "proc_num",
				PropertyName:          attrMap["proc_num"].PropertyName,
				PropertyValue:         i.ProcNum,
				TemplatePropertyValue: t.ProcNum,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.StopCmd.AsDefaultValue) {
		if (t.StopCmd.Value == nil && i.StopCmd != nil) ||
			(t.StopCmd.Value != nil && i.StopCmd == nil) ||
			(t.StopCmd.Value != nil && i.StopCmd != nil && *t.StopCmd.Value != *i.StopCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["stop_cmd"].ID,
				PropertyID:            "stop_cmd",
				PropertyName:          attrMap["stop_cmd"].PropertyName,
				PropertyValue:         i.StopCmd,
				TemplatePropertyValue: t.StopCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.RestartCmd.AsDefaultValue) {
		if (t.RestartCmd.Value == nil && i.RestartCmd != nil) ||
			(t.RestartCmd.Value != nil && i.RestartCmd == nil) ||
			(t.RestartCmd.Value != nil && i.RestartCmd != nil && *t.RestartCmd.Value != *i.RestartCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["restart_cmd"].ID,
				PropertyID:            "restart_cmd",
				PropertyName:          attrMap["restart_cmd"].PropertyName,
				PropertyValue:         i.RestartCmd,
				TemplatePropertyValue: t.RestartCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.ForceStopCmd.AsDefaultValue) {
		if (t.ForceStopCmd.Value == nil && i.ForceStopCmd != nil) ||
			(t.ForceStopCmd.Value != nil && i.ForceStopCmd == nil) ||
			(t.ForceStopCmd.Value != nil && i.ForceStopCmd != nil && *t.ForceStopCmd.Value != *i.ForceStopCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["face_stop_cmd"].ID,
				PropertyID:            "face_stop_cmd",
				PropertyName:          attrMap["face_stop_cmd"].PropertyName,
				PropertyValue:         i.ForceStopCmd,
				TemplatePropertyValue: t.ForceStopCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.FuncName.AsDefaultValue) {
		if (t.FuncName.Value == nil && i.FuncName != nil) ||
			(t.FuncName.Value != nil && i.FuncName == nil) ||
			(t.FuncName.Value != nil && i.FuncName != nil && *t.FuncName.Value != *i.FuncName) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_func_name"].ID,
				PropertyID:            "bk_func_name",
				PropertyName:          attrMap["bk_func_name"].PropertyName,
				PropertyValue:         i.FuncName,
				TemplatePropertyValue: t.FuncName,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.WorkPath.AsDefaultValue) {
		if (t.WorkPath.Value == nil && i.WorkPath != nil) ||
			(t.WorkPath.Value != nil && i.WorkPath == nil) ||
			(t.WorkPath.Value != nil && i.WorkPath != nil && *t.WorkPath.Value != *i.WorkPath) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["work_path"].ID,
				PropertyID:            "work_path",
				PropertyName:          attrMap["work_path"].PropertyName,
				PropertyValue:         i.WorkPath,
				TemplatePropertyValue: t.WorkPath,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.BindInfo.AsDefaultValue) {
		newBindInfo, change, err := t.BindInfo.DiffWithProcessTemplate(i.BindInfo, host)
		if err != nil {
			return nil, err
		}
		if change {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap[common.BKProcBindInfo].ID,
				PropertyID:            common.BKProcBindInfo,
				PropertyName:          attrMap[common.BKProcBindInfo].PropertyName,
				PropertyValue:         i.BindInfo,
				TemplatePropertyValue: newBindInfo,
			})
		}

	}

	if metadata.IsAsDefaultValue(t.Priority.AsDefaultValue) {
		if (t.Priority.Value == nil && i.Priority != nil) ||
			(t.Priority.Value != nil && i.Priority == nil) ||
			(t.Priority.Value != nil && i.Priority != nil && *t.Priority.Value != *i.Priority) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["priority"].ID,
				PropertyID:            "priority",
				PropertyName:          attrMap["priority"].PropertyName,
				PropertyValue:         i.Priority,
				TemplatePropertyValue: t.Priority,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.ReloadCmd.AsDefaultValue) {
		if (t.ReloadCmd.Value == nil && i.ReloadCmd != nil) ||
			(t.ReloadCmd.Value != nil && i.ReloadCmd == nil) ||
			(t.ReloadCmd.Value != nil && i.ReloadCmd != nil && *t.ReloadCmd.Value != *i.ReloadCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["reload_cmd"].ID,
				PropertyID:            "reload_cmd",
				PropertyName:          attrMap["reload_cmd"].PropertyName,
				PropertyValue:         i.ReloadCmd,
				TemplatePropertyValue: t.ReloadCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.ProcessName.AsDefaultValue) {
		if (t.ProcessName.Value == nil && i.ProcessName != nil) ||
			(t.ProcessName.Value != nil && i.ProcessName == nil) ||
			(t.ProcessName.Value != nil && i.ProcessName != nil && *t.ProcessName.Value != *i.ProcessName) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_process_name"].ID,
				PropertyID:            "bk_process_name",
				PropertyName:          attrMap["bk_process_name"].PropertyName,
				PropertyValue:         i.ProcessName,
				TemplatePropertyValue: t.ProcessName,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.PidFile.AsDefaultValue) {
		if (t.PidFile.Value == nil && i.PidFile != nil) ||
			(t.PidFile.Value != nil && i.PidFile == nil) ||
			(t.PidFile.Value != nil && i.PidFile != nil && *t.PidFile.Value != *i.PidFile) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["pid_file"].ID,
				PropertyID:            "pid_file",
				PropertyName:          attrMap["pid_file"].PropertyName,
				PropertyValue:         i.PidFile,
				TemplatePropertyValue: t.PidFile,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.AutoStart.AsDefaultValue) {
		if (t.AutoStart.Value == nil && i.AutoStart != nil) ||
			(t.AutoStart.Value != nil && i.AutoStart == nil) ||
			(t.AutoStart.Value != nil && i.AutoStart != nil && *t.AutoStart.Value != *i.AutoStart) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_start"].ID,
				PropertyID:            "auto_start",
				PropertyName:          attrMap["auto_start"].PropertyName,
				PropertyValue:         i.AutoStart,
				TemplatePropertyValue: t.AutoStart,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.StartCheckSecs.AsDefaultValue) {
		if (t.StartCheckSecs.Value == nil && i.StartCheckSecs != nil) ||
			(t.StartCheckSecs.Value != nil && i.StartCheckSecs == nil) ||
			(t.StartCheckSecs.Value != nil && i.StartCheckSecs != nil && *t.StartCheckSecs.Value != *i.StartCheckSecs) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_start_check_secs"].ID,
				PropertyID:            "bk_start_check_secs",
				PropertyName:          attrMap["bk_start_check_secs"].PropertyName,
				PropertyValue:         i.StartCheckSecs,
				TemplatePropertyValue: t.StartCheckSecs,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.StartCmd.AsDefaultValue) {
		if (t.StartCmd.Value == nil && i.StartCmd != nil) ||
			(t.StartCmd.Value != nil && i.StartCmd == nil) ||
			(t.StartCmd.Value != nil && i.StartCmd != nil && *t.StartCmd.Value != *i.StartCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["start_cmd"].ID,
				PropertyID:            "start_cmd",
				PropertyName:          attrMap["start_cmd"].PropertyName,
				PropertyValue:         i.StartCmd,
				TemplatePropertyValue: t.StartCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.User.AsDefaultValue) {
		if (t.User.Value == nil && i.User != nil) ||
			(t.User.Value != nil && i.User == nil) ||
			(t.User.Value != nil && i.User != nil && *t.User.Value != *i.User) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["user"].ID,
				PropertyID:            "user",
				PropertyName:          attrMap["user"].PropertyName,
				PropertyValue:         i.User,
				TemplatePropertyValue: t.User,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.TimeoutSeconds.AsDefaultValue) {
		if (t.TimeoutSeconds.Value == nil && i.TimeoutSeconds != nil) ||
			(t.TimeoutSeconds.Value != nil && i.TimeoutSeconds == nil) ||
			(t.TimeoutSeconds.Value != nil && i.TimeoutSeconds != nil && *t.TimeoutSeconds.Value != *i.TimeoutSeconds) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["timeout"].ID,
				PropertyID:            "timeout",
				PropertyName:          attrMap["timeout"].PropertyName,
				PropertyValue:         i.TimeoutSeconds,
				TemplatePropertyValue: t.TimeoutSeconds,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.Description.AsDefaultValue) {
		if (t.Description.Value == nil && i.Description != nil) ||
			(t.Description.Value != nil && i.Description == nil) ||
			(t.Description.Value != nil && i.Description != nil && *t.Description.Value != *i.Description) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["description"].ID,
				PropertyID:            "description",
				PropertyName:          attrMap["description"].PropertyName,
				PropertyValue:         i.Description,
				TemplatePropertyValue: t.Description,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.StartParamRegex.AsDefaultValue) {
		if (t.StartParamRegex.Value == nil && i.StartParamRegex != nil) ||
			(t.StartParamRegex.Value != nil && i.StartParamRegex == nil) ||
			(t.StartParamRegex.Value != nil && i.StartParamRegex != nil && *t.StartParamRegex.Value != *i.StartParamRegex) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_start_param_regex"].ID,
				PropertyID:            "bk_start_param_regex",
				PropertyName:          attrMap["bk_start_param_regex"].PropertyName,
				PropertyValue:         i.StartParamRegex,
				TemplatePropertyValue: t.StartParamRegex,
			})
		}
	}

	return changes, nil
}
