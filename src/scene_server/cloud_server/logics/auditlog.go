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
	"net/http"
	
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type CloudAuditLog interface {
	WithPrevious(*rest.Kit, int64) errors.CCError
	WithCurrent(*rest.Kit, int64) errors.CCError
	SaveAuditLog(*rest.Kit, metadata.ActionType) errors.CCError
}

type SyncTaskAuditLog struct {
	logic    *Logics
	header   http.Header
	ownerID  string
	taskName string
	taskID   int64
	content  metadata.CloudSyncTaskOpContent
}

func (lgc *Logics) NewSyncTaskAuditLog(kit *rest.Kit, ownerID string) *SyncTaskAuditLog {
	return &SyncTaskAuditLog{
		logic:   lgc,
		header:  kit.Header,
		ownerID: ownerID,
	}
}

func (log *SyncTaskAuditLog) WithPrevious(kit *rest.Kit, taskID int64) errors.CCError {
	preData, err := log.buildLogData(kit, taskID)
	if err != nil {
		return err
	}
	log.content.PreData = &preData

	return nil
}

func (log *SyncTaskAuditLog) WithCurrent(kit *rest.Kit, taskID int64) errors.CCError {
	curData, err := log.buildLogData(kit, taskID)
	if err != nil {
		return err
	}
	log.content.CurData = &curData

	return nil
}

func (log *SyncTaskAuditLog) buildLogData(kit *rest.Kit, taskID int64) (metadata.CloudSyncTask, errors.CCError) {
	option := metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
	}
	res, err := log.logic.CoreAPI.CoreService().Cloud().SearchSyncTask(kit.Ctx, kit.Header, &option)
	if err != nil {
		return metadata.CloudSyncTask{}, err
	}
	if len(res.Info) <= 0 {
		return metadata.CloudSyncTask{}, kit.CCError.CCErrorf(common.CCErrCloudValidSyncTaskParamFail, common.BKCloudAccountID)
	}

	log.taskID = taskID
	log.taskName = res.Info[0].TaskName
	log.content.Properties = syncTaskAuditLogProperty

	return res.Info[0], nil
}

func (log *SyncTaskAuditLog) SaveAuditLog(kit *rest.Kit, action metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudSyncTaskRes,
		Action:       action,
		OperationDetail: &metadata.CloudSyncTaskOpDetail{
			TaskName: log.taskName,
			TaskID:   log.taskID,
			Details:  &log.content,
		},
	}

	auditResult, err := log.logic.CoreAPI.CoreService().Audit().SaveAuditLog(kit.Ctx, log.header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud sync task audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud sync task audit log failed, err: %s, result: %s,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

type AccountAuditLog struct {
	logic       *Logics
	header      http.Header
	ownerID     string
	accountName string
	accountID   int64
	Content     *metadata.BasicContent
}

func (lgc *Logics) NewAccountAuditLog(kit *rest.Kit, ownerID string) *AccountAuditLog {
	return &AccountAuditLog{
		logic:   lgc,
		header:  kit.Header,
		ownerID: ownerID,
		Content: new(metadata.BasicContent),
	}
}

func (log *AccountAuditLog) WithPrevious(kit *rest.Kit, accountID int64) errors.CCError {
	data, err := log.buildLogData(kit, accountID)
	if err != nil {
		return err
	}
	log.Content.PreData = data

	return nil
}

func (log *AccountAuditLog) WithCurrent(kit *rest.Kit, accountID int64) errors.CCError {
	data, err := log.buildLogData(kit, accountID)
	if err != nil {
		return err
	}
	log.Content.CurData = data

	return nil
}

func (log *AccountAuditLog) SaveAuditLog(kit *rest.Kit, action metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       action,
		OperationDetail: &metadata.BasicOpDetail{
			ResourceID:   log.accountID,
			ResourceName: log.accountName,
			Details:      log.Content,
		},
	}

	auditResult, err := log.logic.CoreAPI.CoreService().Audit().SaveAuditLog(kit.Ctx, log.header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud account audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud account audit log failed, err: %s, result: %s,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (log *AccountAuditLog) buildLogData(kit *rest.Kit, accountID int64) (map[string]interface{}, errors.CCError) {
	log.Content.Properties = accountAuditLogProperty

	cond := metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudAccountID: accountID},
	}
	res, err := log.logic.CoreAPI.CoreService().Cloud().SearchAccount(kit.Ctx, kit.Header, &cond)
	if err != nil {
		return nil, err
	}
	if len(res.Info) <= 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail)
	}
	data := map[string]interface{}{
		common.BKCloudAccountName: res.Info[0].AccountName,
		common.BKCloudVendor:      res.Info[0].CloudVendor,
		common.BKDescriptionField: res.Info[0].Description,
	}

	log.accountName = res.Info[0].AccountName
	log.accountID = accountID

	return data, nil
}

var accountAuditLogProperty = []metadata.Property{
	{"bk_account_name", "账户名称"},
	{"bk_cloud_vendor", "账户类型"},
	{"bk_description", "备注"},
}

var syncTaskAuditLogProperty = []metadata.Property{
	{"bk_task_id", "任务ID"},
	{"bk_task_name", "任务名称"},
	{"bk_account_id", "云账户ID"},
	{"bk_resource_type", "同步资源类型"},
	{"bk_last_sync_time", "上次同步时间"},
}
