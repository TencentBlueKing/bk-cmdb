/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package tools

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

// AddCreateAuditLog add create data audit log
func AddCreateAuditLog(kit *rest.Kit, db dal.RDB, auditData []map[string]interface{}, auditField *AuditStruct) error {
	if auditField.AuditTypeData == nil {
		return nil
	}

	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameAuditLog,
		len(auditData))
	if err != nil {
		blog.Errorf("get next %d audit log IDs failed, err: %v", len(auditData), err)
		return err
	}

	auditLogs := make([]*auditLog, len(auditData))
	for i, data := range auditData {
		var bizID int64
		if auditField.AuditDataField.BizIDField != "" {
			bizID, err = util.GetInt64ByInterface(data[auditField.AuditDataField.BizIDField])
			if err != nil {
				blog.Errorf("get businessID failed, err: %v, auditData: %+v", err, data)
				return err
			}
		}
		resName := util.GetStrByInterface(data[auditField.AuditDataField.ResNameField])
		auditLogs[i] = getAuditLog(int64(ids[i]), data, resName, auditField.AuditTypeData,
			data[auditField.AuditDataField.ResIDField], bizID)
	}

	if err = db.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, auditLogs); err != nil {
		blog.Errorf("add audit log %+v error %v", auditLogs, err)
		return err
	}
	return nil
}

// AddTmpAuditLog add tmp data audit log
func AddTmpAuditLog(kit *rest.Kit, db dal.RDB, auditData []map[string]interface{}, auditField *AuditStruct) error {
	if auditField.AuditTypeData == nil {
		return nil
	}

	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameAuditLog,
		len(auditData))
	if err != nil {
		blog.Errorf("get next %d audit log IDs failed, err: %v", len(auditData), err)
		return err
	}

	auditLogs := make([]*auditLog, len(auditData))
	for i, data := range auditData {
		var bizID int64
		if auditField.AuditDataField.BizIDField != "" {
			bizID, err = util.GetInt64ByInterface(data[auditField.AuditDataField.BizIDField])
			if err != nil {
				blog.Errorf("get businessID failed, err: %v, auditData: %+v", err, data)
				return err
			}
		}
		auditLogs[i] = getAuditLog(int64(ids[i]), data, auditField.ResNames[i], auditField.AuditTypeData,
			data[auditField.AuditDataField.ResIDField], bizID)
	}

	if err = db.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, auditLogs); err != nil {
		blog.Errorf("add audit log %+v error %v", auditLogs, err)
		return err
	}
	return nil
}

func getAuditLog(id int64, data mapstr.MapStr, resName string, auditType *AuditResType,
	resID interface{}, bizID int64) *auditLog {

	auditData := &auditLog{
		ID:          id,
		BusinessID:  bizID,
		Action:      metadata.AuditCreate,
		OperateFrom: metadata.FromCCSystem,
		OperationDetail: &metadata.BasicOpDetail{
			Details: NewBasicContent(data, metadata.AuditCreate),
		},
		OperationTime: time.Now(),
		ResourceName:  resName,
		AuditField:    auditType,
		ResourceID:    resID,
	}

	return auditData
}

// NewBasicContent get basicContent by data and self.
func NewBasicContent(data map[string]interface{}, action metadata.ActionType) *metadata.BasicContent {
	var basicDetail *metadata.BasicContent
	switch action {
	case metadata.AuditCreate:
		basicDetail = &metadata.BasicContent{
			CurData: data,
		}
	default:
		blog.Errorf("audit log action type %s not support", action)
		return nil
	}
	return basicDetail
}

type auditLog struct {
	ID                 int64                    `bson:"id"`
	User               string                   `bson:"user"`
	Action             metadata.ActionType      `bson:"action"`
	OperateFrom        metadata.OperateFromType `bson:"operate_from"`
	OperationDetail    metadata.DetailFactory   `bson:"operation_detail"`
	OperationTime      time.Time                `bson:"operation_time"`
	AppCode            string                   `bson:"code,omitempty"`
	RequestID          string                   `bson:"rid,omitempty"`
	ExtendResourceName string                   `bson:"extend_resource_name"`
	AuditField         *AuditResType            `bson:",inline"`
	BusinessID         int64                    `bson:"bk_biz_id,omitempty"`
	ResourceID         interface{}              `bson:"resource_id"`
	ResourceName       string                   `bson:"resource_name"`
}

// AuditResType audit type field
type AuditResType struct {
	AuditType    metadata.AuditType    `bson:"audit_type"`
	ResourceType metadata.ResourceType `bson:"resource_type"`
}

// AuditDataField audit data field
type AuditDataField struct {
	BizIDField   string
	ResIDField   string
	ResNameField string
}

type AuditStruct struct {
	AuditDataField *AuditDataField `bson:",inline"`
	AuditTypeData  *AuditResType   `bson:",inline"`
	ResNames       []string
}
