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
