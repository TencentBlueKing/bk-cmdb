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

package operation

import (
	"context"

	"configcenter/src/common/condition"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// WrapperResult the data wrapper
type WrapperResult struct {
	datas []inst.Inst
}

// AuditInterface audit log methods
type AuditInterface interface {
	CreateSnapshot(instID int64, cond mapstr.MapStr) *WrapperResult
	CommitCreateLog(preData, currData *WrapperResult, inst inst.Inst)
	CommitDeleteLog(preData, currData *WrapperResult, inst inst.Inst)
	CommitUpdateLog(preData, currData *WrapperResult, inst inst.Inst)
}

type auditLog struct {
	client apimachinery.ClientSetInterface
	inst   InstOperationInterface
	params types.ContextParams
	obj    model.Object
}

func (a *auditLog) commitSnapshot(preData, currData *WrapperResult, action auditoplog.AuditOpType) {
	if nil == currData {
		blog.Errorf("[audit] the curr data is empty")
		return
	}

	for idx, currItem := range currData.datas {

		var preDataTmp interface{}
		if nil != preData {
			if idx < len(preData.datas) {
				preDataTmp = preData.datas[idx].GetValues()
			}
		}

		desc := ""
		switch action {
		case auditoplog.AuditOpTypeAdd:
			desc = "create " + a.obj.GetObjectType()
		case auditoplog.AuditOpTypeDel:
			desc = "delete " + a.obj.GetObjectType()
		case auditoplog.AuditOpTypeModify:
			desc = "update " + a.obj.GetObjectType()

		}

		id, err := currItem.GetInstID()
		if nil != err {
			blog.Errorf("[audit]failed to get the inst id, error info is %s", err.Error())
			return
		}

		headers := []Header{}
		attrs, err := a.obj.GetAttributes()
		if nil != err {
			blog.Errorf("[audit]failed to get the object(%s)' attribute, error info is %s", a.obj.GetID(), err.Error())
			return
		}
		for _, attr := range attrs {
			headers = append(headers, Header{
				PropertyID:   attr.GetID(),
				PropertyName: attr.GetName(),
			})
		}

		data := common.KvMap{
			common.BKContentField: Content{
				CurData: currItem.GetValues(),
				PreData: preDataTmp,
				Headers: headers,
			},
			common.BKOpDescField:   desc,
			common.BKOpTypeField:   auditoplog.AuditOpTypeAdd,
			common.BKOpTargetField: a.obj.GetID(),
			"inst_id":              id,
		}

		bizID, err := currItem.GetValues().String(common.BKAppIDField)
		if nil != err {
			blog.V(3).Infof("[audit] failed to get the bizid from the data(%#v), error info is %s", currItem.GetValues(), err.Error())
			bizID = "0"
		}
		//fmt.Println("the data:", data, "obj:", a.obj.GetID())
		switch a.obj.GetObjectType() {
		default:

			rsp, err := a.client.AuditController().AddObjectLog(context.Background(), a.params.SupplierAccount, bizID, a.params.User, a.params.Header, data)
			if nil != err {
				blog.Errorf("[audit] failed to add audit log, error info is %s", err.Error())
				return
			}
			if !rsp.Result {
				blog.Errorf("[audit] failed to add audit log, error info is %s", rsp.ErrMsg)
				return
			}
		case common.BKInnerObjIDApp, common.BKINnerObjIDObject:
			rsp, err := a.client.AuditController().AddObjectLog(context.Background(), a.params.SupplierAccount, bizID, a.params.User, a.params.Header, data)
			if nil != err {
				blog.Errorf("[audit] failed to add audit log, error info is %s", err.Error())
				return
			}
			if !rsp.Result {
				blog.Errorf("[audit] failed to add audit log, error info is %s", rsp.ErrMsg)
				return
			}
		case common.BKInnerObjIDModule:
			rsp, err := a.client.AuditController().AddModuleLog(context.Background(), a.params.SupplierAccount, bizID, a.params.User, a.params.Header, data)
			if nil != err {
				blog.Errorf("[audit] failed to add audit log, error info is %s", err.Error())
				return
			}
			if !rsp.Result {
				blog.Errorf("[audit] failed to add audit log, error info is %s", rsp.ErrMsg)
				return
			}
		case common.BKInnerObjIDSet:
			rsp, err := a.client.AuditController().AddSetLog(context.Background(), a.params.SupplierAccount, bizID, a.params.User, a.params.Header, data)
			if nil != err {
				blog.Errorf("[audit] failed to add audit log, error info is %s", err.Error())
				return
			}
			if !rsp.Result {
				blog.Errorf("[audit] failed to add audit log, error info is %s", rsp.ErrMsg)
				return
			}
		}
	}
}

func (a *auditLog) CreateSnapshot(instID int64, cond mapstr.MapStr) *WrapperResult {

	query := &metadata.QueryInput{}

	if instID >= 0 {
		innerCond := condition.CreateCondition()
		innerCond.Field(a.obj.GetInstIDFieldName()).Eq(instID)
		cond.Merge(innerCond.ToMapStr())
	}

	query.Condition = cond
	_, insts, err := a.inst.FindInst(a.params, a.obj, query, false)
	if nil != err {
		blog.Errorf("[audit] failed to create the snapshot, error info is %s", err.Error())
	}

	result := &WrapperResult{}
	for _, inst := range insts {
		result.datas = append(result.datas, inst)
	}

	return result
}

func (a *auditLog) CommitCreateLog(preData, currData *WrapperResult, inst inst.Inst) {
	if nil == currData {
		currData = &WrapperResult{}
		currData.datas = append(currData.datas, inst)
	}
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeAdd)
}

func (a *auditLog) CommitDeleteLog(preData, currData *WrapperResult, inst inst.Inst) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeDel)
}

func (a *auditLog) CommitUpdateLog(preData, currData *WrapperResult, inst inst.Inst) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeModify)
}
