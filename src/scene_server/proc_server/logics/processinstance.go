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
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

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
		return nil, errors.New(ret.Code, ret.ErrMsg)
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
		return nil, errors.New(ret.Code, ret.ErrMsg)
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

func (lgc *Logic) UpdateProcessInstance(kit *rest.Kit, procID int64, info mapstr.MapStr) error {
	delete(info, common.BkSupplierAccount)
	option := metadata.UpdateOption{
		Data: info,
		Condition: map[string]interface{}{
			common.BKProcessIDField:  procID,
			common.BkSupplierAccount: kit.SupplierAccount,
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
		return err
	}

	if !result.Result {
		blog.Errorf("rid: %s, delete process instance: %d failed, err: %s", kit.Rid, procIDs, result.ErrMsg)
		return kit.CCError.Error(result.Code)
	}

	return nil
}

func (lgc *Logic) CreateProcessInstance(kit *rest.Kit, process *metadata.Process) (int64, errors.CCErrorCoder) {
	processBytes, err := json.Marshal(process)
	if err != nil {
		return 0, kit.CCError.CCError(common.CCErrCommJsonEncode)
	}
	mData := mapstr.MapStr{}
	if err := json.Unmarshal(processBytes, &mData); nil != err && 0 != len(processBytes) {
		return 0, kit.CCError.CCError(common.CCErrCommJsonDecode)
	}
	inputParam := metadata.CreateModelInstance{
		Data: mData,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstance failed, http request failed, err: %+v, rid: %s", err, kit.Rid)
		return 0, errors.CCHttpError
	}

	if !result.Result {
		blog.Errorf("rid: %s, create process instance: %+v failed, err: %s", kit.Rid, process, result.ErrMsg)
		return 0, errors.New(result.Code, result.ErrMsg)
	}

	return int64(result.Data.Created.ID), nil
}

// it works to find the different attribute value between the process instance and it's bounded process template.
// return with the changed attribute's details.
func (lgc *Logic) DiffWithProcessTemplate(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute) []metadata.ProcessChangedAttribute {
	changes := make([]metadata.ProcessChangedAttribute, 0)
	if t == nil || i == nil {
		return changes
	}

	lgc.checkProcNumAsDefaultValue(t, i, attrMap, changes)
	lgc.checkStopCmdAsDefaultValue(t, i, attrMap, changes)
	lgc.checkRestartCmdAsDefaultValue(t, i, attrMap, changes)
	lgc.checkForceStopCmdAsDefaultValue(t, i, attrMap, changes)
	lgc.checkFuncNameAsDefaultValue(t, i, attrMap, changes)
	lgc.checkWorkPathAsDefaultValue(t, i, attrMap, changes)
	lgc.checkBindIPAsDefaultValue(t, i, attrMap, changes)
	lgc.checkPriorityAsDefaultValue(t, i, attrMap, changes)
	lgc.checkReloadCmdAsDefaultValue(t, i, attrMap, changes)
	lgc.checkProcessNameAsDefaultValue(t, i, attrMap, changes)
	lgc.checkPortAsDefaultValue(t, i, attrMap, changes)
	lgc.checkPidFileAsDefaultValue(t, i, attrMap, changes)
	lgc.checkAutoStartAsDefaultValue(t, i, attrMap, changes)
	lgc.checkAutoTimeGapSecondsAsDefaultValue(t, i, attrMap, changes)
	lgc.checkStartCmdAsDefaultValue(t, i, attrMap, changes)
	lgc.checkFuncIDAsDefaultValue(t, i, attrMap, changes)
	lgc.checkUserAsDefaultValue(t, i, attrMap, changes)
	lgc.checkTimeoutSecondsAsDefaultValue(t, i, attrMap, changes)
	lgc.checkProtocolAsDefaultValue(t, i, attrMap, changes)
	lgc.checkDescriptionAsDefaultValue(t, i, attrMap, changes)
	lgc.checkStartParamRegexAsDefaultValue(t, i, attrMap, changes)

	return changes
}

func (lgc *Logic) checkProcNumAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.ProcNum.AsDefaultValue) {
		if (t.ProcNum.Value == nil && i.ProcNum != nil) || (t.ProcNum.Value != nil && i.ProcNum == nil) ||
			(t.ProcNum.Value != nil && *t.ProcNum.Value != *i.ProcNum) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["proc_num"].ID,
				PropertyID:            "proc_num",
				PropertyName:          attrMap["proc_num"].PropertyName,
				PropertyValue:         i.ProcNum,
				TemplatePropertyValue: t.ProcNum,
			})
		}
	}
}

func (lgc *Logic) checkStopCmdAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkRestartCmdAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkForceStopCmdAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkFuncNameAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkWorkPathAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkBindIPAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.BindIP.AsDefaultValue) {
		if (t.BindIP.Value == nil && i.BindIP != nil) ||
			(t.BindIP.Value != nil && i.BindIP == nil) ||
			(t.BindIP.Value != nil && i.BindIP != nil && t.BindIP.Value.IP() != *i.BindIP) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bind_ip"].ID,
				PropertyID:            "bind_ip",
				PropertyName:          attrMap["bind_ip"].PropertyName,
				PropertyValue:         i.BindIP,
				TemplatePropertyValue: t.BindIP.Value.IP(),
			})
		}
	}
}

func (lgc *Logic) checkPriorityAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkReloadCmdAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkProcessNameAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkPortAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.Port.AsDefaultValue) {
		if (t.Port.Value == nil && i.Port != nil) || (t.Port.Value != nil && i.Port == nil) ||
			(t.Port.Value != nil && i.Port != nil && *t.Port.Value != *i.Port) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["port"].ID,
				PropertyID:            "port",
				PropertyName:          attrMap["port"].PropertyName,
				PropertyValue:         i.Port,
				TemplatePropertyValue: t.Port,
			})
		}
	}
}

func (lgc *Logic) checkPidFileAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.PidFile.AsDefaultValue) {
		if (t.PidFile.Value == nil && i.PidFile != nil) || (t.PidFile.Value != nil && i.PidFile == nil) ||
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
}

func (lgc *Logic) checkAutoStartAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.AutoStart.AsDefaultValue) {
		if (t.AutoStart.Value == nil && i.AutoStart != nil) || (t.AutoStart.Value != nil && i.AutoStart == nil) ||
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
}

func (lgc *Logic) checkAutoTimeGapSecondsAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.AutoTimeGapSeconds.AsDefaultValue) {
		if (t.AutoTimeGapSeconds.Value == nil && i.AutoTimeGap != nil) ||
			(t.AutoTimeGapSeconds.Value != nil && i.AutoTimeGap == nil) ||
			(t.AutoTimeGapSeconds.Value != nil && i.AutoTimeGap != nil && *t.AutoTimeGapSeconds.Value != *i.AutoTimeGap) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_time_gap"].ID,
				PropertyID:            "auto_time_gap",
				PropertyName:          attrMap["auto_time_gap"].PropertyName,
				PropertyValue:         i.AutoTimeGap,
				TemplatePropertyValue: t.AutoTimeGapSeconds,
			})
		}
	}
}

func (lgc *Logic) checkStartCmdAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.StartCmd.AsDefaultValue) {
		if (t.StartCmd.Value == nil && i.StartCmd != nil) || (t.StartCmd.Value != nil && i.StartCmd == nil) ||
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
}

func (lgc *Logic) checkFuncIDAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.FuncID.AsDefaultValue) {
		if (t.FuncID.Value == nil && i.FuncID != nil) || (t.FuncID.Value != nil && i.FuncID == nil) ||
			(t.FuncID.Value != nil && i.FuncID != nil && *t.FuncID.Value != *i.FuncID) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_func_id"].ID,
				PropertyID:            "bk_func_id",
				PropertyName:          attrMap["bk_func_id"].PropertyName,
				PropertyValue:         i.FuncID,
				TemplatePropertyValue: t.FuncID,
			})
		}
	}
}

func (lgc *Logic) checkUserAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.User.AsDefaultValue) {
		if (t.User.Value == nil && i.User != nil) || (t.User.Value != nil && i.User == nil) ||
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
}

func (lgc *Logic) checkTimeoutSecondsAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkProtocolAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.Protocol.AsDefaultValue) {
		if (t.Protocol.Value == nil && i.Protocol != nil) || (t.Protocol.Value != nil && i.Protocol == nil) ||
			(t.Protocol.Value != nil && i.Protocol != nil && *t.Protocol.Value != *i.Protocol) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["protocol"].ID,
				PropertyID:            "protocol",
				PropertyName:          attrMap["protocol"].PropertyName,
				PropertyValue:         i.Protocol,
				TemplatePropertyValue: t.Protocol,
			})
		}
	}
}

func (lgc *Logic) checkDescriptionAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
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
}

func (lgc *Logic) checkStartParamRegexAsDefaultValue(t *metadata.ProcessProperty, i *metadata.Process,
	attrMap map[string]metadata.Attribute, changes []metadata.ProcessChangedAttribute) {
	if metadata.IsAsDefaultValue(t.StartParamRegex.AsDefaultValue) {
		if (t.StartParamRegex.Value == nil && i.StartParamRegex != nil) ||
			(t.StartParamRegex.Value != nil && i.StartParamRegex == nil) ||
			(t.StartParamRegex.Value != nil && i.StartParamRegex != nil &&
				*t.StartParamRegex.Value != *i.StartParamRegex) {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_start_param_regex"].ID,
				PropertyID:            "bk_start_param_regex",
				PropertyName:          attrMap["bk_start_param_regex"].PropertyName,
				PropertyValue:         i.StartParamRegex,
				TemplatePropertyValue: t.StartParamRegex,
			})
		}
	}
}
