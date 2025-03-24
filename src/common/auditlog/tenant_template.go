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

package auditlog

import (
	"time"

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

type tenantTemplateAuditLog struct{}

// GenerateAuditLog generate audit log
func (t *tenantTemplateAuditLog) GenerateAuditLog(data []*TenantTmpAuditOpts) []metadata.AuditLog {

	auditLogs := make([]metadata.AuditLog, 0)

	for _, item := range data {
		auditLogs = append(auditLogs, metadata.AuditLog{
			Action:      metadata.AuditTenantInit,
			OperateFrom: metadata.FromCCSystem,
			OperationDetail: &TenantTmpDetail{
				Type:    item.Type,
				CurData: item,
			},
			OperationTime: metadata.Time{Time: time.Now()},
			ResourceName:  item.ResourceName,
			AuditType:     metadata.PlatformSetting,
			ResourceType:  metadata.TenantTemplateRes,
			ResourceID:    item.ResourceID,
		})
	}
	return auditLogs
}

// SaveAuditLog save audit log
func (t *tenantTemplateAuditLog) SaveAuditLog(kit *rest.Kit, db local.DB, auditLogs ...metadata.AuditLog) error {
	logRows := make([]metadata.AuditLog, 0)

	ids, err := mongodb.Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameAuditLog, len(auditLogs))
	if err != nil {
		blog.Errorf("get next audit log id failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	for index, log := range auditLogs {
		if log.OperationDetail == nil {
			continue
		}
		if log.OperateFrom == "" {
			log.OperateFrom = metadata.FromUser
		}
		// ResourceName is assigned index, length must be less than 1024, so resourceName only save NameFieldMaxLength.
		if len(log.ResourceName) > common.NameFieldMaxLength {
			log.ResourceName = log.ResourceName[:common.NameFieldMaxLength]
		}
		log.User = kit.User
		if appCode := httpheader.GetAppCode(kit.Header); len(appCode) > 0 {
			log.AppCode = appCode
		}

		log.RequestID = kit.Rid
		log.OperationTime = metadata.Now()
		log.ID = int64(ids[index])

		logRows = append(logRows, log)
	}

	if len(logRows) == 0 {
		return nil
	}
	return db.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, logRows)
}

// NewTenantTemplateAuditLog new tenant template audit log
func NewTenantTemplateAuditLog() *tenantTemplateAuditLog {
	return &tenantTemplateAuditLog{}
}

// TenantTmpDetail tenant template audit log detail
type TenantTmpDetail struct {
	PreData      interface{}                  `json:"pre_data" bson:"pre_data"`
	CurData      interface{}                  `json:"cur_data" bson:"cur_data"`
	UpdateFields interface{}                  `json:"update_fields" bson:"update_fields"`
	Type         tenanttmp.TenantTemplateType `json:"type" bson:"type"`
}

// WithName tenant template with name
func (op *TenantTmpDetail) WithName() string {
	return "TenantTemplateDetail"
}

// TenantTmpAuditOpts tenant template audit log option
type TenantTmpAuditOpts struct {
	Data         interface{}
	ResourceID   int64
	ResourceName string
	Type         tenanttmp.TenantTemplateType
}
