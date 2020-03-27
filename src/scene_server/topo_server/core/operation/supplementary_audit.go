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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
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
	CommitCreateLog(preData, currData *WrapperResult, inst inst.Inst, nonInnerAttributes []model.AttributeInterface)
	CommitDeleteLog(preData, currData *WrapperResult, inst inst.Inst)
	CommitUpdateLog(preData, currData *WrapperResult, inst inst.Inst, nonInnerAttributes []model.AttributeInterface)
}

type auditLog struct {
	client apimachinery.ClientSetInterface
	inst   InstOperationInterface
	params types.ContextParams
	obj    model.Object
}

// nonInnerAttributes 用于加速，避免不必要的数据查询(批量创建实例时，每次创建实例都会执行nonInnerAttributes)
func (a *auditLog) commitSnapshot(preData, currData *WrapperResult, action auditoplog.AuditOpType, nonInnerAttributes []model.AttributeInterface) {

	var targetData *WrapperResult
	isPreItem := false
	if nil != currData {
		targetData = currData
	} else if nil != preData {
		isPreItem = true
		targetData = preData
	} else {
		blog.Errorf("[audit] the curr data is empty, rid: %s", a.params.ReqID)
		return
	}

	if nonInnerAttributes == nil {
		var err error
		nonInnerAttributes, err = a.obj.GetNonInnerAttributes()
		if nil != err {
			blog.Errorf("[audit]failed to get the object(%s)' attribute, error info is %s, rid: %s", a.obj.Object().ObjectID, err.Error(), a.params.ReqID)
			return
		}
	}
	for _, targetItem := range targetData.datas {

		id, err := targetItem.GetInstID()
		if nil != err {
			blog.Errorf("[audit]failed to get the inst id, error info is %s, rid: %s", err.Error(), a.params.ReqID)
			return
		}

		var preDataTmp, currDataTmp mapstr.MapStr
		if !isPreItem {
			currDataTmp = targetItem.GetValues()
		} else {
			preDataTmp = targetItem.GetValues()
		}

		if nil != preData && !isPreItem {
			for _, preItem := range preData.datas {
				preID, err := preItem.GetInstID()
				if nil != err {
					blog.Errorf("[audit]failed to get the inst id, error info is %s, rid: %s", err.Error(), a.params.ReqID)
					continue
				}
				if id == preID {
					preDataTmp = preItem.GetValues()
				}
			}
		}

		desc := ""
		switch action {
		case auditoplog.AuditOpTypeAdd:
			desc = "create " + a.obj.GetObjectType()
		case auditoplog.AuditOpTypeDel:
			desc = "delete " + a.obj.GetObjectType()
		case auditoplog.AuditOpTypeModify:
			if currDataTmp[common.BKDataStatusField] != preDataTmp[common.BKDataStatusField] {
				switch currDataTmp[common.BKDataStatusField] {
				case common.DataStatusDisabled:
					desc = "disabled " + a.obj.GetObjectType()
				case common.DataStatusEnable:
					desc = "enable " + a.obj.GetObjectType()
				default:
					desc = "update " + a.obj.GetObjectType()
				}
			} else {
				desc = "update " + a.obj.GetObjectType()
			}

		}

		headers := []Header{}
		for _, attr := range nonInnerAttributes {
			headers = append(headers, Header{
				PropertyID:   attr.Attribute().PropertyID,
				PropertyName: attr.Attribute().PropertyName,
			})
		}
		var bizID int64
		if targetItem.GetValues() != nil {
			if _, exist := targetItem.GetValues()[common.BKAppIDField]; exist {
				if biz, err := targetItem.GetValues().Int64(common.BKAppIDField); nil != err {
					blog.V(3).Infof("[audit] failed to get the biz id from the data(%#v), error info is %s, rid: %s", targetItem.GetValues(), err.Error(), a.params.ReqID)
					return
				} else {
					bizID = biz
				}
			}
		}

		auditlog := metadata.SaveAuditLogParams{
			ID:    id,
			Model: a.obj.GetObjectID(),
			Content: Content{
				CurData: currDataTmp,
				PreData: preDataTmp,
				Headers: headers,
			},
			OpDesc: desc,
			OpType: action,
			BizID:  bizID,
		}

		auditresp, err := a.client.CoreService().Audit().SaveAuditLog(context.Background(), a.params.Header, auditlog)
		if err != nil {
			blog.V(3).Infof("[audit] failed to get the bizid from the data(%#v), error info is %s, rid: %s", targetItem.GetValues(), err.Error(), a.params.ReqID)
			return
		}
		if !auditresp.Result {
			blog.V(3).Infof("[audit] failed to get the bizid from the data(%#v), resp info is %v, rid: %s", targetItem.GetValues(), auditresp, a.params.ReqID)
			return
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
		blog.Errorf("[audit] failed to create the snapshot, error info is %s, rid: %s", err.Error(), a.params.ReqID)
	}

	result := &WrapperResult{}
	for _, inst := range insts {
		result.datas = append(result.datas, inst)
	}

	return result
}

func (a *auditLog) CommitCreateLog(preData, currData *WrapperResult, inst inst.Inst, nonInnerAttributes []model.AttributeInterface) {
	if nil == currData {
		currData = &WrapperResult{}
		currData.datas = append(currData.datas, inst)
	}
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeAdd, nonInnerAttributes)
}

func (a *auditLog) CommitDeleteLog(preData, currData *WrapperResult, inst inst.Inst) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeDel, nil)
}

func (a *auditLog) CommitUpdateLog(preData, currData *WrapperResult, inst inst.Inst, nonInnerAttributes []model.AttributeInterface) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeModify, nonInnerAttributes)
}
