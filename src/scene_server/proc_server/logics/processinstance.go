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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logic) GetProcessInstanceWithID(kit *rest.Kit, procID int64) (*metadata.Process, error) {
	condition := map[string]interface{}{
		common.BKProcessIDField: procID,
	}

	reqParam := new(metadata.QueryCondition)
	reqParam.Condition = condition
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("rid: %s get process instance with procID: %d failed, err: %v", kit.Rid, procID, err)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !ret.Result {
		blog.Errorf("rid: %s get process instance with procID: %d failed, err: %v", kit.Rid, procID, ret.ErrMsg)
		return nil, kit.CCError.New(ret.Code, ret.ErrMsg)

	}

	process := new(metadata.Process)
	if len(ret.Data.Info) != 0 {
		if err := ret.Data.Info[0].ToStructByTag(process, "field"); err != nil {
			return nil, kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
		}
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
		return kit.CCError.Error(result.Code)
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

func (lgc *Logic) DeleteProcessInstanceBatch(kit *rest.Kit, procID []int64) error {
	cond := condition.CreateCondition()
	cond.AddContionItem(condition.ConditionItem{Field: "id", Operator: "$in", Value: procID})
	option := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
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

func (lgc *Logic) CreateProcessInstance(kit *rest.Kit, proc *metadata.Process) (uint64, error) {
	inst := metadata.CreateModelInstance{
		Data: mapstr.NewFromStruct(proc, "field"),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, &inst)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		blog.Errorf("rid: %s, create process instance: %+v failed, err: %s", kit.Rid, proc, result.ErrMsg)
		return 0, kit.CCError.Error(result.Code)
	}

	return result.Data.Created.ID, nil
}

// it works to find the different attribute value between the process instance and it's bounded process template.
// return with the changed attribute's details.
func (lgc *Logic) GetDifferenceInProcessTemplateAndInstance(t *metadata.ProcessProperty, i *metadata.Process, attrMap map[string]metadata.Attribute) []metadata.ProcessChangedAttribute {
	changes := make([]metadata.ProcessChangedAttribute, 0)
	if t == nil || i == nil {
		return changes
	}

	if t.ProcNum.Value != nil {
		if *t.ProcNum.Value != i.ProcNum {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["proc_num"].ID,
				PropertyID:            "proc_num",
				PropertyName:          attrMap["proc_num"].PropertyName,
				PropertyValue:         i.ProcNum,
				TemplatePropertyValue: t.ProcNum,
			})
		}
	}

	if t.StopCmd.Value != nil {
		if *t.StopCmd.Value != i.StopCmd {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["stop_cmd"].ID,
				PropertyID:            "stop_cmd",
				PropertyName:          attrMap["stop_cmd"].PropertyName,
				PropertyValue:         i.StopCmd,
				TemplatePropertyValue: t.StopCmd,
			})
		}
	}

	if t.RestartCmd.Value != nil {
		if *t.RestartCmd.Value != i.RestartCmd {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["restart_cmd"].ID,
				PropertyID:            "restart_cmd",
				PropertyName:          attrMap["restart_cmd"].PropertyName,
				PropertyValue:         i.RestartCmd,
				TemplatePropertyValue: t.RestartCmd,
			})
		}
	}

	if t.ForceStopCmd.Value != nil {
		if *t.ForceStopCmd.Value != i.ForceStopCmd {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["face_stop_cmd"].ID,
				PropertyID:            "face_stop_cmd",
				PropertyName:          attrMap["face_stop_cmd"].PropertyName,
				PropertyValue:         i.ForceStopCmd,
				TemplatePropertyValue: t.ForceStopCmd,
			})
		}
	}

	if t.FuncName.Value != nil {
		if *t.FuncName.Value != i.FuncName {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_func_name"].ID,
				PropertyID:            "bk_func_name",
				PropertyName:          attrMap["bk_func_name"].PropertyName,
				PropertyValue:         i.FuncName,
				TemplatePropertyValue: t.FuncName,
			})
		}
	}

	if t.WorkPath.Value != nil {
		if *t.WorkPath.Value != i.WorkPath {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["work_path"].ID,
				PropertyID:            "work_path",
				PropertyName:          attrMap["work_path"].PropertyName,
				PropertyValue:         i.WorkPath,
				TemplatePropertyValue: t.WorkPath,
			})
		}
	}

	if t.BindIP.Value != nil {
		if *t.BindIP.Value != i.BindIP {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bind_ip"].ID,
				PropertyID:            "bind_ip",
				PropertyName:          attrMap["bind_ip"].PropertyName,
				PropertyValue:         i.BindIP,
				TemplatePropertyValue: t.BindIP,
			})
		}
	}

	if t.Priority.Value != nil {
		if *t.Priority.Value != i.Priority {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["priority"].ID,
				PropertyID:            "priority",
				PropertyName:          attrMap["priority"].PropertyName,
				PropertyValue:         i.Priority,
				TemplatePropertyValue: t.Priority,
			})
		}
	}

	if t.ReloadCmd.Value != nil {
		if *t.ReloadCmd.Value != i.ReloadCmd {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["reload_cmd"].ID,
				PropertyID:            "reload_cmd",
				PropertyName:          attrMap["reload_cmd"].PropertyName,
				PropertyValue:         i.ReloadCmd,
				TemplatePropertyValue: t.ReloadCmd,
			})
		}
	}

	if t.ProcessName.Value != nil {
		if *t.ProcessName.Value != i.ProcessName {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_process_name"].ID,
				PropertyID:            "bk_process_name",
				PropertyName:          attrMap["bk_process_name"].PropertyName,
				PropertyValue:         i.ProcessName,
				TemplatePropertyValue: t.ProcessName,
			})
		}
	}

	if t.Port.Value != nil {
		if *t.Port.Value != i.Port {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["port"].ID,
				PropertyID:            "port",
				PropertyName:          attrMap["port"].PropertyName,
				PropertyValue:         i.Port,
				TemplatePropertyValue: t.Port,
			})
		}
	}

	if t.PidFile.Value != nil {
		if *t.PidFile.Value != i.PidFile {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["pid_file"].ID,
				PropertyID:            "pid_file",
				PropertyName:          attrMap["pid_file"].PropertyName,
				PropertyValue:         i.PidFile,
				TemplatePropertyValue: t.PidFile,
			})
		}
	}

	if t.AutoStart.Value != nil {
		if *t.AutoStart.Value != i.AutoStart {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_start"].ID,
				PropertyID:            "auto_start",
				PropertyName:          attrMap["auto_start"].PropertyName,
				PropertyValue:         i.AutoStart,
				TemplatePropertyValue: t.AutoStart,
			})
		}
	}

	if t.AutoTimeGapSeconds.Value != nil {
		if *t.AutoTimeGapSeconds.Value != i.AutoTimeGap {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["auto_time_gap"].ID,
				PropertyID:            "auto_time_gap",
				PropertyName:          attrMap["auto_time_gap"].PropertyName,
				PropertyValue:         i.AutoTimeGap,
				TemplatePropertyValue: t.AutoTimeGapSeconds,
			})
		}
	}

	if t.StartCmd.Value != nil {
		if *t.StartCmd.Value != i.StartCmd {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["start_cmd"].ID,
				PropertyID:            "start_cmd",
				PropertyName:          attrMap["start_cmd"].PropertyName,
				PropertyValue:         i.StartCmd,
				TemplatePropertyValue: t.StartCmd,
			})
		}
	}

	if t.FuncID.Value != nil {
		if *t.FuncID.Value != i.FuncID {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["bk_func_id"].ID,
				PropertyID:            "bk_func_id",
				PropertyName:          attrMap["bk_func_id"].PropertyName,
				PropertyValue:         i.FuncID,
				TemplatePropertyValue: t.FuncID,
			})
		}
	}

	if t.User.Value != nil {
		if *t.User.Value != i.User {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["user"].ID,
				PropertyID:            "user",
				PropertyName:          attrMap["user"].PropertyName,
				PropertyValue:         i.User,
				TemplatePropertyValue: t.User,
			})
		}
	}

	if t.TimeoutSeconds.Value != nil {
		if *t.TimeoutSeconds.Value != i.TimeoutSeconds {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["timeout"].ID,
				PropertyID:            "timeout",
				PropertyName:          attrMap["timeout"].PropertyName,
				PropertyValue:         i.TimeoutSeconds,
				TemplatePropertyValue: t.TimeoutSeconds,
			})
		}
	}

	if t.Protocol.Value != nil {
		if *t.Protocol.Value != i.Protocol {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["protocol"].ID,
				PropertyID:            "protocol",
				PropertyName:          attrMap["protocol"].PropertyName,
				PropertyValue:         i.Protocol,
				TemplatePropertyValue: t.Protocol,
			})
		}
	}

	if t.Description.Value != nil {
		if *t.Description.Value != i.Description {
			changes = append(changes, metadata.ProcessChangedAttribute{
				ID:                    attrMap["description"].ID,
				PropertyID:            "description",
				PropertyName:          attrMap["description"].PropertyName,
				PropertyValue:         i.Description,
				TemplatePropertyValue: t.Description,
			})
		}
	}

	return changes
}

// this function works to create a new process instance from a process template.
func (lgc *Logic) CreateProcessInstanceFromProcessTemplate(t *metadata.ProcessProperty) *metadata.Process {
	p := new(metadata.Process)
	if t.ProcNum.Value != nil {
		p.ProcNum = *t.ProcNum.Value
	}

	if t.ProcNum.Value != nil {
		p.ProcNum = *t.ProcNum.Value
	}

	if t.StopCmd.Value != nil {
		p.StopCmd = *t.StopCmd.Value
	}

	if t.RestartCmd.Value != nil {
		p.RestartCmd = *t.RestartCmd.Value
	}

	if t.ForceStopCmd.Value != nil {
		p.ForceStopCmd = *t.ForceStopCmd.Value
	}

	if t.ProcessID.Value != nil {
		p.ProcessID = *t.ProcessID.Value
	}

	if t.FuncName.Value != nil {
		p.FuncName = *t.FuncName.Value
	}

	if t.WorkPath.Value != nil {
		p.WorkPath = *t.WorkPath.Value
	}

	if t.BindIP.Value != nil {
		p.BindIP = *t.BindIP.Value
	}

	if t.Priority.Value != nil {
		p.Priority = *t.Priority.Value
	}

	if t.ReloadCmd.Value != nil {
		p.ReloadCmd = *t.ReloadCmd.Value
	}

	if t.ProcessName.Value != nil {
		p.ProcessName = *t.ProcessName.Value
	}

	if t.Port.Value != nil {
		p.Port = *t.Port.Value
	}

	if t.PidFile.Value != nil {
		p.PidFile = *t.PidFile.Value
	}

	if t.AutoStart.Value != nil {
		p.AutoStart = *t.AutoStart.Value
	}

	if t.AutoTimeGapSeconds.Value != nil {
		p.AutoTimeGap = *t.AutoTimeGapSeconds.Value
	}

	if t.StartCmd.Value != nil {
		p.StartCmd = *t.StartCmd.Value
	}

	if t.FuncID.Value != nil {
		p.FuncID = *t.FuncID.Value
	}

	if t.User.Value != nil {
		p.User = *t.User.Value
	}

	if t.TimeoutSeconds.Value != nil {
		p.TimeoutSeconds = *t.TimeoutSeconds.Value
	}

	if t.Protocol.Value != nil {
		p.Protocol = *t.Protocol.Value
	}

	if t.Description.Value != nil {
		p.Description = *t.Description.Value
	}

	return p
}
