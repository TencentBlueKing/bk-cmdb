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

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// instanceAuditLog provides methods to generate and save instance audit log.
type instanceAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of instance.
func (i *instanceAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, objID string, data []mapstr.MapStr) (
	[]metadata.AuditLog, error) {
	return i.generateAuditLog(parameter, objID, data)
}

// GenerateAuditLogByCondGetData generate audit log of instance, auto get current instance by objID and condition.
func (i *instanceAuditLog) GenerateAuditLogByCondGetData(parameter *generateAuditCommonParameter, objID string,
	condition map[string]interface{}) ([]metadata.AuditLog, error) {
	data, err := i.getInstByCond(parameter.kit, objID, condition, nil)
	if err != nil {
		blog.ErrorJSON("get instances failed, err: %s, condition: %s, rid: %s", err, condition, parameter.kit.Rid)
		return nil, err
	}
	return i.generateAuditLog(parameter, objID, data)
}

func (i *instanceAuditLog) generateAuditLog(parameter *generateAuditCommonParameter, objID string, data []mapstr.MapStr) (
	[]metadata.AuditLog, error) {
	auditLogs := make([]metadata.AuditLog, len(data))
	kit := parameter.kit
	isMainline, err := i.isMainline(parameter.kit, objID)
	if err != nil {
		blog.Errorf("check if object is mainline failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	for index, inst := range data {
		id, err := util.GetInt64ByInterface(inst[metadata.GetInstIDFieldByObjID(objID)])
		if err != nil {
			blog.ErrorJSON("failed to get inst id, error info is %s, inst: %s, rid: %s", err.Error(), inst, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, objID, metadata.GetInstIDFieldByObjID(objID), "int", err.Error())
		}

		var bizID int64
		if _, exist := inst[common.BKAppIDField]; exist {
			bizID, err = util.GetInt64ByInterface(inst[common.BKAppIDField])
		} else if _, exist := inst[metadata.BKMetadata]; exist {
			bizID, err = metadata.ParseBizIDFromData(inst)
		}
		if err != nil {
			blog.ErrorJSON("failed to get biz id from metadata, error info is %s, inst: %s, rid: %s", err.Error(), inst, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, objID, common.BKAppIDField, "int", err.Error())
		}

		action := parameter.action
		updateFields := parameter.updateFields

		var details *metadata.BasicContent
		switch action {
		case metadata.AuditCreate:
			details = &metadata.BasicContent{
				CurData: inst,
			}
		case metadata.AuditDelete:
			details = &metadata.BasicContent{
				PreData: inst,
			}
		case metadata.AuditUpdate:
			if updateFields[common.BKDataStatusField] != inst[common.BKDataStatusField] {
				switch updateFields[common.BKDataStatusField] {
				case string(common.DataStatusDisabled):
					action = metadata.AuditArchive
				case string(common.DataStatusEnable):
					action = metadata.AuditRecover
				}
			}

			details = &metadata.BasicContent{
				PreData:      inst,
				UpdateFields: updateFields,
			}
		}

		auditLog := metadata.AuditLog{
			AuditType:    metadata.GetAuditTypeByObjID(objID, isMainline),
			ResourceType: metadata.GetResourceTypeByObjID(objID, isMainline),
			Action:       action,
			BusinessID:   bizID,
			ResourceID:   id,
			OperateFrom:  parameter.operateFrom,
			ResourceName: util.GetStrByInterface(inst[metadata.GetInstNameFieldName(objID)]),
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					Details: details,
				},
				ModelID: objID,
			},
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

func (i *instanceAuditLog) isMainline(kit *rest.Kit, objID string) (bool, error) {
	cond := &metadata.QueryCondition{Condition:
	map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline},
	}

	asst, err := i.clientSet.Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("[audit] failed to find mainline association, err: %v, rid: %s", err, kit.Rid)
		return false, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !asst.Result {
		blog.Errorf("failed to find mainline association, err code: %d , err msg: %s, rid: %s", err, asst.Code, asst.ErrMsg, kit.Rid)
		return false, asst.CCError()
	}

	for _, mainline := range asst.Data.Info {
		if mainline.AsstObjID == objID {
			return true, nil
		}
	}

	return false, nil
}

func NewInstanceAudit(clientSet coreservice.CoreServiceClientInterface) *instanceAuditLog {
	return &instanceAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
