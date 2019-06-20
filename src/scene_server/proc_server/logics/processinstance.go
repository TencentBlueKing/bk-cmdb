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
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logic) ListProcessInstanceWithIDs(kit *rest.Kit, procIDs []int64) ([]metadata.Process, error) {
	cond := condition.CreateCondition()
	cond.AddConditionItem(condition.ConditionItem{Field: common.BKProcessIDField, Operator: common.BKDBIN, Value: procIDs})

	reqParam := new(metadata.QueryCondition)
	reqParam.Condition = cond.ToMapStr()
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("rid: %s list process instance with procID: %d failed, err: %v", kit.Rid, procIDs, err)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !ret.Result {
		blog.Errorf("rid: %s list process instance with procID: %d failed, err: %v", kit.Rid, procIDs, ret.ErrMsg)
		return nil, kit.CCError.New(ret.Code, ret.ErrMsg)

	}

	processes := make([]metadata.Process, 0)
	for _, p := range ret.Data.Info {
		process := new(metadata.Process)
		if err := p.MarshalJSONInto(process); err != nil {
			return nil, kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
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
		blog.Errorf("rid: %s get process instance with procID: %d failed, err: %v", kit.Rid, procID, err)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !ret.Result {
		blog.Errorf("rid: %s get process instance with procID: %d failed, err: %v", kit.Rid, procID, ret.ErrMsg)
		return nil, errors.New(ret.Code, ret.ErrMsg)

	}

	process := new(metadata.Process)
	if len(ret.Data.Info) == 0 {
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if err := ret.Data.Info[0].MarshalJSONInto(process); err != nil {
		blog.Errorf("GetProcessInstanceWithID fai", err)
		return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return process, nil
}

func (lgc *Logic) UpdateProcessInstance(kit *rest.Kit, procID int64, info mapstr.MapStr) error {
	option := metadata.UpdateOption{
		Data: info,
		Condition: map[string]interface{}{
			common.BKProcessIDField: procID,
		},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		return err
	}

	if !result.Result {
		blog.Errorf("rid: %s, update process instance: %d failed, err: %s", kit.Rid, procID, result.ErrMsg)
		return kit.CCError.New(result.Code, result.ErrMsg)
	}
	return nil
}

func (lgc *Logic) DeleteProcessInstance(kit *rest.Kit, procID int64) error {
	option := metadata.DeleteOption{
		Condition: map[string]interface{}{
			common.BKProcessIDField: procID,
		},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		return err
	}

	if !result.Result {
		blog.Errorf("rid: %s, delete process instance: %d failed, err: %s", kit.Rid, procID, result.ErrMsg)
		return kit.CCError.Error(result.Code)
	}

	return nil
}

func (lgc *Logic) DeleteProcessInstanceBatch(kit *rest.Kit, procIDs []int64) error {
	cond := condition.CreateCondition()
	cond.AddConditionItem(condition.ConditionItem{Field: common.BKProcessIDField, Operator: condition.BKDBIN, Value: procIDs})
	option := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}
	result, err := lgc.CoreAPI.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &option)
	if err != nil {
		return err
	}

	if !result.Result {
		blog.Errorf("rid: %s, delete process instance: %d failed, err: %s", kit.Rid, procIDs, result.ErrMsg)
		return kit.CCError.Error(result.Code)
	}

	return nil
}

func (lgc *Logic) CreateProcessInstance(kit *rest.Kit, proc *metadata.Process) (int64, error) {
	inst := metadata.CreateModelInstance{
		Data: mapstr.NewFromStruct(proc, "field"),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &inst)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		blog.Errorf("rid: %s, create process instance: %+v failed, err: %s", kit.Rid, proc, result.ErrMsg)
		return 0, errors.New(result.Code, result.ErrMsg)
	}

	return int64(result.Data.Created.ID), nil
}

// it works to find the different attribute value between the process instance and it's bounded process template.
// return with the changed attribute's details.
func (lgc *Logic) DiffWithProcessTemplate(t *metadata.ProcessProperty, i *metadata.Process, attrMap map[string]metadata.Attribute) []metadata.ProcessChangedAttribute {
	changes := make([]metadata.ProcessChangedAttribute, 0)
	if t == nil || i == nil {
		return changes
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
		if (t.StopCmd.Value == nil && len(i.StopCmd) != 0) || (t.StopCmd.Value != nil && *t.StopCmd.Value != i.StopCmd) {
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
		if (t.RestartCmd.Value == nil && len(i.RestartCmd) != 0) || (t.RestartCmd.Value != nil && *t.RestartCmd.Value != i.RestartCmd) {
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
		if (t.ForceStopCmd.Value == nil && len(i.ForceStopCmd) != 0) || (t.ForceStopCmd.Value != nil && *t.ForceStopCmd.Value != i.ForceStopCmd) {
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
		if (t.FuncName.Value == nil && len(i.FuncName) != 0) || (t.FuncName.Value != nil && *t.FuncName.Value != i.FuncName) {
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
		if (t.WorkPath.Value == nil && len(i.WorkPath) != 0) || (t.WorkPath.Value != nil && *t.WorkPath.Value != i.WorkPath) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["work_path"].ID,
				PropertyID:            "work_path",
				PropertyName:          attrMap["work_path"].PropertyName,
				PropertyValue:         i.WorkPath,
				TemplatePropertyValue: t.WorkPath,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.BindIP.AsDefaultValue) {
		if (t.BindIP.Value == nil && i.BindIP != nil) || (t.BindIP.Value != nil && i.BindIP == nil) || (t.BindIP.Value != nil && i.BindIP != nil && *t.BindIP.Value != *i.BindIP) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bind_ip"].ID,
				PropertyID:            "bind_ip",
				PropertyName:          attrMap["bind_ip"].PropertyName,
				PropertyValue:         i.BindIP,
				TemplatePropertyValue: t.BindIP,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.Priority.AsDefaultValue) {
		if (t.Priority.Value == nil && i.Priority != nil) || (t.Priority.Value != nil && i.Priority == nil) || (t.Priority.Value != nil && *t.Priority.Value != *i.Priority) {
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
		if (t.ReloadCmd.Value == nil && len(i.ReloadCmd) != 0) || (t.ReloadCmd.Value != nil && *t.ReloadCmd.Value != i.ReloadCmd) {
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
		if (t.ProcessName.Value == nil && len(i.ProcessName) != 0) || (t.ProcessName.Value != nil && *t.ProcessName.Value != i.ProcessName) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_process_name"].ID,
				PropertyID:            "bk_process_name",
				PropertyName:          attrMap["bk_process_name"].PropertyName,
				PropertyValue:         i.ProcessName,
				TemplatePropertyValue: t.ProcessName,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.Port.AsDefaultValue) {
		if (t.Port.Value == nil && len(i.Port) != 0) || (t.Port.Value != nil && *t.Port.Value != i.Port) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["port"].ID,
				PropertyID:            "port",
				PropertyName:          attrMap["port"].PropertyName,
				PropertyValue:         i.Port,
				TemplatePropertyValue: t.Port,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.PidFile.AsDefaultValue) {
		if (t.PidFile.Value == nil && len(i.PidFile) != 0) || (t.PidFile.Value != nil && *t.PidFile.Value != i.PidFile) {
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
		if (t.AutoStart.Value == nil && i.AutoStart != false) || (t.AutoStart.Value != nil && *t.AutoStart.Value != i.AutoStart) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_start"].ID,
				PropertyID:            "auto_start",
				PropertyName:          attrMap["auto_start"].PropertyName,
				PropertyValue:         i.AutoStart,
				TemplatePropertyValue: t.AutoStart,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.AutoTimeGapSeconds.AsDefaultValue) {
		if (t.AutoTimeGapSeconds.Value == nil && i.AutoTimeGap != nil) || (t.AutoTimeGapSeconds.Value != nil && i.AutoTimeGap == nil) || (t.AutoTimeGapSeconds.Value != nil && *t.AutoTimeGapSeconds.Value != *i.AutoTimeGap) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_time_gap"].ID,
				PropertyID:            "auto_time_gap",
				PropertyName:          attrMap["auto_time_gap"].PropertyName,
				PropertyValue:         i.AutoTimeGap,
				TemplatePropertyValue: t.AutoTimeGapSeconds,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.StartCmd.AsDefaultValue) {
		if (t.StartCmd.Value == nil && len(i.StartCmd) != 0) || (t.StartCmd.Value != nil && *t.StartCmd.Value != i.StartCmd) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["start_cmd"].ID,
				PropertyID:            "start_cmd",
				PropertyName:          attrMap["start_cmd"].PropertyName,
				PropertyValue:         i.StartCmd,
				TemplatePropertyValue: t.StartCmd,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.FuncID.AsDefaultValue) {
		if (t.FuncID.Value == nil && len(i.FuncID) != 0) || (t.FuncID.Value != nil && *t.FuncID.Value != i.FuncID) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_func_id"].ID,
				PropertyID:            "bk_func_id",
				PropertyName:          attrMap["bk_func_id"].PropertyName,
				PropertyValue:         i.FuncID,
				TemplatePropertyValue: t.FuncID,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.User.AsDefaultValue) {
		if (t.User.Value == nil && len(i.User) != 0) || (t.User.Value != nil && *t.User.Value != i.User) {
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
		if (t.TimeoutSeconds.Value == nil && i.TimeoutSeconds != nil) || (t.TimeoutSeconds.Value != nil && i.TimeoutSeconds == nil) || (t.TimeoutSeconds.Value != nil && *t.TimeoutSeconds.Value != *i.TimeoutSeconds) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["timeout"].ID,
				PropertyID:            "timeout",
				PropertyName:          attrMap["timeout"].PropertyName,
				PropertyValue:         i.TimeoutSeconds,
				TemplatePropertyValue: t.TimeoutSeconds,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.Protocol.AsDefaultValue) {
		if (t.Protocol.Value == nil && len(i.Protocol) != 0) || (t.Protocol.Value != nil && *t.Protocol.Value != i.Protocol) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["protocol"].ID,
				PropertyID:            "protocol",
				PropertyName:          attrMap["protocol"].PropertyName,
				PropertyValue:         i.Protocol,
				TemplatePropertyValue: t.Protocol,
			})
		}
	}

	if metadata.IsAsDefaultValue(t.Description.AsDefaultValue) {
		if (t.Description.Value == nil && len(i.Description) != 0) || (t.Description.Value != nil && *t.Description.Value != i.Description) {
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
		if (t.StartParamRegex.Value == nil && len(i.StartParamRegex) != 0) || (t.StartParamRegex.Value != nil && *t.StartParamRegex.Value != i.StartParamRegex) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_start_param_regex"].ID,
				PropertyID:            "bk_start_param_regex",
				PropertyName:          attrMap["bk_start_param_regex"].PropertyName,
				PropertyValue:         i.StartParamRegex,
				TemplatePropertyValue: t.StartParamRegex,
			})
		}
	}

	return changes
}

// if process instance is not same with the process template, then update the process instance's value,
// and return the updated process, with a true bool value.
