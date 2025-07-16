/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	collectorlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

// watchAudit watch audit log event and push to audit center
func (a *Audit) watchAudit() error {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	name := "audit"
	tokenHandler := tokenhandler.NewSingleTokenHandler(name, mongodb.Client())

	startAtTime, err := tokenHandler.GetStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get %s start watch time failed, err: %v", name, err)
		return err
	}

	operationType := types.Insert
	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name: name,
			WatchOpt: &types.WatchOptions{
				Options: types.Options{
					OperationType:           &operationType,
					Filter:                  make(mapstr.MapStr),
					EventStruct:             new(metadata.AuditLog),
					Collection:              common.BKTableNameAuditLog,
					StartAtTime:             startAtTime,
					WatchFatalErrorCallback: tokenHandler.ResetWatchToken,
				},
			},
			TokenHandler: tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 2,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: a.doBatch,
		},
		BatchSize: 100,
	}

	if err = a.loopW.WithBatch(loopOptions); err != nil {
		blog.Errorf("watch %s failed, err: %v", name, err)
		return err
	}

	return nil
}

// doBatch batch handle audit event
func (a *Audit) doBatch(es []*types.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	rid := es[0].ID()

	logs := make([]*logsv1.ResourceLogs, 0)

	for idx := range es {
		event := es[idx]

		if event.OperationType != types.Insert {
			blog.V(4).Infof("received ignored audit event operation type: %s, doc: %s, rid: %s", event.OperationType,
				event.DocBytes, rid)
			continue
		}

		resLog, isValid := a.convertToResourceLog(event, rid)
		if !isValid {
			blog.V(4).Infof("received invalid audit event: %+v, rid: %s", *event, rid)
			continue
		}
		logs = append(logs, resLog)

		blog.V(5).Infof("watch audit, received oid: %s, op-time: %s event, rid: %s", event.Oid,
			event.ClusterTime.String(), rid)
	}

	if len(logs) == 0 {
		return false
	}

	req := &collectorlogs.ExportLogsServiceRequest{
		ResourceLogs: logs,
	}
	err := a.client.ReportAuditData(context.Background(), make(http.Header), req)
	if err != nil {
		blog.ErrorJSON("report audit log data failed, err: %s, req: %s, rid: %s", err, req, rid)
		return true
	}

	return false
}

func (a *Audit) convertToResourceLog(event *types.Event, rid string) (*logsv1.ResourceLogs, bool) {
	auditLog, ok := event.Document.(*metadata.AuditLog)
	if !ok {
		blog.Errorf("received invalid audit event doc(%#v), oid: %s, rid: %s", event.Document, event.Oid, rid)
		return nil, false
	}

	detail, err := json.Marshal(auditLog.OperationDetail)
	if err != nil {
		blog.Errorf("marshal audit log detail failed, err: %v, auditLog: %v, rid: %s", err, *auditLog, rid)
		return nil, false
	}

	var instanceID string
	if intID, err := util.GetInt64ByInterface(auditLog.ResourceID); err == nil {
		instanceID = strconv.FormatInt(intID, 10)
	} else {
		instanceID = util.GetStrByInterface(auditLog.ResourceID)
	}

	logAttributes := []*commonv1.KeyValue{
		a.convertToStrKeyValue("event_id", strconv.FormatInt(auditLog.ID, 10)),
		a.convertToStrKeyValue("event_content", string(detail)),
		a.convertToStrKeyValue("request_id", auditLog.RequestID),
		a.convertToIntKeyValue("start_time", auditLog.OperationTime.UnixMilli()),
		a.convertToIntKeyValue("end_time", auditLog.OperationTime.UnixMilli()),
		a.convertToStrKeyValue("bk_app_code", a.conf.AppCode),
		a.convertToStrKeyValue("action_id", string(auditLog.Action)),
		a.convertToStrKeyValue("resource_type_id", string(auditLog.ResourceType)),
		a.convertToStrKeyValue("instance_id", instanceID),
		a.convertToStrKeyValue("instance_name", auditLog.ResourceName),
		a.convertToIntKeyValue("instance_sensitivity", 0),
		a.convertToStrKeyValue("instance_data", "{}"),
		a.convertToStrKeyValue("instance_origin_data", "{}"),
		a.convertToStrKeyValue("extend_data", "{}"),
		a.convertToIntKeyValue("result_code", 0),
		a.convertToStrKeyValue("result_content", ""),
		a.convertToStrKeyValue("bk_log_scope", "bk_audit_event"),
	}

	logAttributes = append(logAttributes, a.parseUserInfo(auditLog)...)
	logAttributes = append(logAttributes, a.parseScopeInfo(auditLog)...)

	return &logsv1.ResourceLogs{
		Resource: &resourcev1.Resource{
			Attributes: []*commonv1.KeyValue{
				a.convertToStrKeyValue(string(semconv.ServiceNameKey), a.conf.AppCode),
				a.convertToStrKeyValue("bk.data.token", a.conf.Token),
			},
		},
		ScopeLogs: []*logsv1.ScopeLogs{{
			Scope: &commonv1.InstrumentationScope{
				Name: a.conf.AppCode,
			},
			LogRecords: []*logsv1.LogRecord{{
				TimeUnixNano:         uint64(auditLog.OperationTime.UnixNano()),
				ObservedTimeUnixNano: uint64(time.Now().UnixNano()),
				SeverityNumber:       logsv1.SeverityNumber_SEVERITY_NUMBER_INFO,
				SeverityText:         "INFO",
				Attributes:           logAttributes,
			}},
		}},
		SchemaUrl: semconv.SchemaURL,
	}, true
}

func (a *Audit) parseUserInfo(auditLog *metadata.AuditLog) []*commonv1.KeyValue {
	var accessType, userIdentifyType int64
	if auditLog.OperateFrom == metadata.FromUser {
		accessType = 0
		userIdentifyType = 0
	} else {
		accessType = -1
		userIdentifyType = 1
	}

	var userIdentifySrc, userIdentifySrcUsername string
	if auditLog.AppCode != "" {
		userIdentifySrc = "bk_app_code"
		userIdentifySrcUsername = auditLog.AppCode
		accessType = 1
	}

	return []*commonv1.KeyValue{
		a.convertToStrKeyValue("username", auditLog.User),
		a.convertToStrKeyValue("user_identify_tenant_id", auditLog.SupplierAccount),
		a.convertToIntKeyValue("user_identify_type", userIdentifyType),
		a.convertToStrKeyValue("user_identify_src", userIdentifySrc),
		a.convertToStrKeyValue("user_identify_src_username", userIdentifySrcUsername),
		a.convertToIntKeyValue("access_type", accessType),
		a.convertToStrKeyValue("access_source_ip", ""),
		a.convertToStrKeyValue("access_user_agent", ""),
	}
}

func (a *Audit) parseScopeInfo(auditLog *metadata.AuditLog) []*commonv1.KeyValue {
	var scopeType, scopeID string
	if auditLog.BusinessID != 0 {
		scopeType = "biz"
		scopeID = strconv.FormatInt(auditLog.BusinessID, 10)
	}

	return []*commonv1.KeyValue{
		a.convertToStrKeyValue("scope_type", scopeType),
		a.convertToStrKeyValue("scope_id", scopeID),
	}
}

func (a *Audit) convertToStrKeyValue(key, value string) *commonv1.KeyValue {
	return &commonv1.KeyValue{
		Key:   key,
		Value: &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: value}},
	}
}

func (a *Audit) convertToIntKeyValue(key string, value int64) *commonv1.KeyValue {
	return &commonv1.KeyValue{
		Key:   key,
		Value: &commonv1.AnyValue{Value: &commonv1.AnyValue_IntValue{IntValue: value}},
	}
}
