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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

type tenantTemplateAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log
func (t *tenantTemplateAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter,
	data []interface{}, auditOpt *AuditOpts) []metadata.AuditLog {

	auditLogs := make([]metadata.AuditLog, 0)

	for index, item := range data {
		resourceName := ""
		if len(auditOpt.ResourceName) != 0 {
			resourceName = auditOpt.ResourceName[index]
		}
		auditLogs = append(auditLogs, metadata.AuditLog{
			Action:      metadata.AuditTenantInit,
			OperateFrom: metadata.FromCCSystem,
			OperationDetail: &metadata.TenantTmpDetail{
				Type: auditOpt.Type,
				TenantTmpBasicOpDetail: metadata.TenantTmpBasicOpDetail{
					Details: parameter.NewTmpBasicContent(item),
				},
			},
			OperationTime: metadata.Time{Time: time.Now()},
			ResourceName:  resourceName,
			AuditType:     metadata.PlatformSetting,
			ResourceType:  metadata.TenantTemplateRes,
			ResourceID:    auditOpt.ResourceID[index],
		})
	}
	return auditLogs
}

// SaveAuditLog save audit log
func (t *tenantTemplateAuditLog) SaveAuditLog(kit *rest.Kit, db local.DB, auditLogs ...metadata.AuditLog) error {
	logRows := make([]metadata.AuditLog, 0)

	ids, err := mongodb.Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameAuditLog, len(auditLogs))
	if err != nil {
		blog.Errorf("get next audit log id failed, err: %s", err.Error())
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
		if rid := kit.Rid; len(rid) > 0 {
			log.RequestID = kit.Rid
		}
		log.OperationTime = metadata.Now()
		log.ID = int64(ids[index])

		logRows = append(logRows, log)
	}

	if len(logRows) == 0 {
		return nil
	}
	return db.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, logRows)
}

// AuditOpts audit options
type AuditOpts struct {
	Type         string
	ResourceID   []int64
	ResourceName []string
}

// NewTenantTemplateAuditLog new tenant template audit log
func NewTenantTemplateAuditLog() *tenantTemplateAuditLog {
	return &tenantTemplateAuditLog{}
}
