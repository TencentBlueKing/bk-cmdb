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

package service

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type EventAudit struct {
	ctx              context.Context
	header           http.Header
	clientSet        apimachinery.ClientSetInterface
	auditType        metadata.AuditType
	resourceType     metadata.ResourceType
	subscriptionID   int64
	subscriptionName string
	systemName       string
	preData          metadata.Subscription
	curData          metadata.Subscription
}

func NewEventAudit(ctx context.Context, header http.Header, clientSet apimachinery.ClientSetInterface) *EventAudit {
	return &EventAudit{
		ctx:          ctx,
		header:       header,
		clientSet:    clientSet,
		auditType:    metadata.EventPushType,
		resourceType: metadata.EventPushRes,
	}
}

func (log *EventAudit) SaveAuditLog(auditAction metadata.ActionType) errors.CCError {
	preData := mapstr.MapStr{}
	curData := mapstr.MapStr{}
	switch auditAction {
	case metadata.AuditDelete:
		preData = log.preData.ToMapStr()
	case metadata.AuditCreate:
		curData = log.curData.ToMapStr()
	case metadata.AuditUpdate:
		preData = log.preData.ToMapStr()
		curData = log.curData.ToMapStr()
	}
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.BasicOpDetail{
			ResourceID:   log.subscriptionID,
			ResourceName: log.subscriptionName,
			Details: &metadata.BasicContent{
				PreData: preData,
				CurData: curData,
			},
		},
	}
	auditResult, err := log.clientSet.CoreService().Audit().SaveAuditLog(log.ctx, log.header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %+v", auditAction, log.resourceType, err, auditResult)
		return err
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %+v", auditAction, log.resourceType, auditResult.ErrMsg, auditResult)
		return errors.New(common.CCErrAuditSaveLogFailed, auditResult.ErrMsg)
	}
	return nil
}

func (log *EventAudit) buildSnapshotForPre(sub metadata.Subscription) *EventAudit {
	log.preData = sub
	log.subscriptionName = log.preData.SubscriptionName
	log.systemName = log.preData.SystemName
	return log
}

func (log *EventAudit) buildSnapshotForCur(sub metadata.Subscription) *EventAudit {
	log.curData = sub
	log.subscriptionName = log.curData.SubscriptionName
	log.systemName = log.curData.SystemName
	return log
}
