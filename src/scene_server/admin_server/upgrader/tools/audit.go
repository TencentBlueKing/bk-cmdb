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
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

// AddCreateAuditLog add create data audit log
func AddCreateAuditLog(kit *rest.Kit, db dal.RDB, data map[string]interface{}, fields *AuditField) error {

	action := metadata.AuditCreate
	id, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.Errorf("get next audit log ID failed, err: %s", err)
		return err
	}
	log := auditLog{
		ID:          int64(id),
		Action:      action,
		OperateFrom: metadata.FromCCSystem,
		OperationDetail: &metadata.BasicOpDetail{
			Details: NewBasicContent(data, action),
		},
		OperationTime: NewTime().CreateTime,
		AuditField:    fields,
	}

	if err = db.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, log); err != nil {
		blog.Errorf("add audit log %s error %v", log, err)
		return err
	}
	return nil
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
	ID                 int64                    `bson:"ID"`
	User               string                   `bson:"user"`
	Action             metadata.ActionType      `bson:"action"`
	OperateFrom        metadata.OperateFromType `bson:"operate_from"`
	OperationDetail    metadata.DetailFactory   `bson:"operation_detail"`
	OperationTime      time.Time                `bson:"operation_time"`
	AppCode            string                   `bson:"code,omitempty"`
	RequestID          string                   `bson:"rid,omitempty"`
	ExtendResourceName string                   `bson:"extend_resource_name"`
	AuditField         *AuditField              `bson:",inline"`
}

// AuditField audit field
type AuditField struct {
	AuditType    metadata.AuditType    `bson:"audit_type"`
	BusinessID   int64                 `bson:"bk_biz_id,omitempty"`
	ResourceID   interface{}           `bson:"resource_id"`
	ResourceName string                `bson:"resource_name"`
	ResourceType metadata.ResourceType `bson:"resource_type"`
}
