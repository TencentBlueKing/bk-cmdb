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
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
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
	kit    *rest.Kit
	obj    model.Object
}

// nonInnerAttributes 用于加速，避免不必要的数据查询(批量创建实例时，每次创建实例都会执行nonInnerAttributes)
func (a *auditLog) commitSnapshot(preData, currData *WrapperResult, action metadata.ActionType, nonInnerAttributes []model.AttributeInterface) {
	var targetData *WrapperResult
	isPreItem := false
	if nil != currData {
		targetData = currData
	} else if nil != preData {
		isPreItem = true
		targetData = preData
	} else {
		blog.Errorf("[audit] the curr data is empty, rid: %s", a.kit.Rid)
		return
	}

	if nonInnerAttributes == nil {
		var err error
		nonInnerAttributes, err = a.obj.GetNonInnerAttributes()
		if nil != err {
			blog.Errorf("[audit]failed to get the object(%s)' attribute, error info is %s, rid: %s", a.obj.Object().ObjectID, err.Error(), a.kit.Rid)
			return
		}
	}
	for _, targetItem := range targetData.datas {

		id, err := targetItem.GetInstID()
		if nil != err {
			blog.Errorf("[audit]failed to get the inst id, error info is %s, rid: %s", err.Error(), a.kit.Rid)
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
					blog.Errorf("[audit]failed to get the inst id, error info is %s, rid: %s", err.Error(), a.kit.Rid)
					continue
				}
				if id == preID {
					preDataTmp = preItem.GetValues()
				}
			}
		}

		properties := make([]metadata.Property, 0)
		for _, attr := range nonInnerAttributes {
			properties = append(properties, metadata.Property{
				PropertyID:   attr.Attribute().PropertyID,
				PropertyName: attr.Attribute().PropertyName,
			})
		}
		var bizID int64
		if targetItem.GetValues() != nil {
			if _, exist := targetItem.GetValues()[common.BKAppIDField]; exist {
				if biz, err := targetItem.GetValues().Int64(common.BKAppIDField); nil != err {
					blog.V(3).Infof("[audit] failed to get the biz id from the data(%#v), error info is %s, rid: %s", targetItem.GetValues(), err.Error(), a.kit.Rid)
					return
				} else {
					bizID = biz
				}
			}
		}
		bizName := ""
		if bizID > 0 {
			bizName, err = auditlog.NewAudit(a.client, a.kit.Header).GetInstNameByID(a.kit.Ctx, common.BKInnerObjIDApp, bizID)
			if err != nil {
				return
			}
		}

		objID := targetItem.GetObject().GetObjectID()
		instName, err := targetItem.GetInstName()
		if err != nil {
			blog.V(3).Infof("[audit] failed to get the inst name from the data(%#v), error info is %s, rid: %s", targetItem.GetValues(), err.Error(), a.kit.Rid)
			return
		}
		if action == metadata.AuditUpdate {
			if currDataTmp[common.BKDataStatusField] != preDataTmp[common.BKDataStatusField] {
				switch currDataTmp[common.BKDataStatusField] {
				case string(common.DataStatusDisabled):
					action = metadata.AuditArchive
				case string(common.DataStatusEnable):
					action = metadata.AuditRecover
				}
			}
		}
		auditLog := metadata.AuditLog{
			AuditType:    metadata.GetAuditTypeByObjID(objID),
			ResourceType: metadata.GetResourceTypeByObjID(objID),
			Action:       action,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					BusinessID:   bizID,
					BusinessName: bizName,
					ResourceID:   id,
					ResourceName: instName,
					Details: &metadata.BasicContent{
						PreData:    preDataTmp,
						CurData:    currDataTmp,
						Properties: properties,
					},
				},
				ModelID: objID,
			},
		}

		// add biz topology label for mainline instance
		asst, err := a.client.CoreService().Association().ReadModelAssociation(context.Background(), a.kit.Header, &metadata.QueryCondition{Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
		if err != nil || !asst.Result {
			blog.V(3).Infof("[audit] failed to find mainline association, err: %v, resp: %v, rid: %s", err, asst, a.kit.Rid)
			return
		}
		if objID != common.BKInnerObjIDApp {
			for _, mainline := range asst.Data.Info {
				if mainline.ObjectID == objID || mainline.AsstObjID == objID {
					auditLog.Label = map[string]string{
						metadata.LabelBizTopology: "",
					}
					break
				}
			}
		}

		auditResp, err := a.client.CoreService().Audit().SaveAuditLog(context.Background(), a.kit.Header, auditLog)
		if err != nil {
			blog.V(3).Infof("[audit] failed to save audit log(%#v), err: %s, rid: %s", auditLog, err.Error(), a.kit.Rid)
			return
		}
		if !auditResp.Result {
			blog.V(3).Infof("[audit] failed to save audit log(%#v), err: %s, rid: %s", auditLog, auditResp.ErrMsg, a.kit.Rid)
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
	_, insts, err := a.inst.FindInst(a.kit, a.obj, query, false)
	if nil != err {
		blog.Errorf("[audit] failed to create the snapshot, error info is %s, rid: %s", err.Error(), a.kit.Rid)
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
	a.commitSnapshot(preData, currData, metadata.AuditCreate, nonInnerAttributes)
}

func (a *auditLog) CommitDeleteLog(preData, currData *WrapperResult, inst inst.Inst) {
	a.commitSnapshot(preData, currData, metadata.AuditDelete, nil)
}

func (a *auditLog) CommitUpdateLog(preData, currData *WrapperResult, inst inst.Inst, nonInnerAttributes []model.AttributeInterface) {
	a.commitSnapshot(preData, currData, metadata.AuditUpdate, nonInnerAttributes)
}
