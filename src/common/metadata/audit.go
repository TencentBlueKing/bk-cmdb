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

package metadata

import (
	"configcenter/src/common/auditoplog"
)

// AuditHostLogParams add single host log parammeter
type AuditHostLogParams struct {
	Content interface{}            `json:"content"`
	OpDesc  string                 `json:"op_desc"`
	InnerIP string                 `json:"bk_host_innerip"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	HostID  int64                  `json:"inst_id"`
}

// AuditHostsLogParams add multiple host log parameter
type AuditHostsLogParams struct {
	Content []auditoplog.AuditLogExt `json:"content"`
	OpDesc  string                   `json:"op_desc"`
	OpType  auditoplog.AuditOpType   `json:"op_type"`
}

// AuditObjParams add object single log parameter
type AuditObjParams struct {
	Content  interface{}            `json:"content"`
	OpDesc   string                 `json:"op_desc"`
	OpType   auditoplog.AuditOpType `json:"op_type"`
	OpTarget string                 `json:"op_target"`
	InstID   int64                  `json:"inst_id"`
}

// AuditObjsParams add object multiple log parameter
type AuditObjsParams struct {
	Content  []auditoplog.AuditLogContext `json:"content"`
	OpDesc   string                       `json:"op_desc"`
	OpType   auditoplog.AuditOpType       `json:"op_type"`
	OpTarget string                       `json:"op_target"`
}

// AuditProcParams add process single log parameter
type AuditProcParams struct {
	Content interface{}            `json:"content"`
	OpDesc  string                 `json:"op_desc"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	ProcID  int64                  `json:"inst_id"`
}

// AuditProcsParams add process multiple log parameter
type AuditProcsParams struct {
	Content []auditoplog.AuditLogContext `json:"bk_content"`
	OpDesc  string                       `json:"bk_op_desc"`
	OpType  auditoplog.AuditOpType       `json:"bk_op_type"`
}

// AuditModuleParams add module  single log parammete
type AuditModuleParams struct {
	Content  interface{}            `json:"content"`
	OpDesc   string                 `json:"op_desc"`
	OpType   auditoplog.AuditOpType `json:"op_type"`
	ModuleID int64                  `json:"inst_id"`
}

// AuditModuleParams add module multiple log parammete
type AuditModulesParams struct {
	Content []auditoplog.AuditLogContext `json:"content"`
	OpDesc  string                       `json:"op_desc"`
	OpType  auditoplog.AuditOpType       `json:"op_type"`
}

// AuditAppParams add application log parameter
type AuditAppParams struct {
	Content string                 `json:"content"`
	OpDesc  string                 `json:"op_desc"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	AppID   int64                  `json:"inst_id"`
}

// AuditSetParams add set single log parameter
type AuditSetParams struct {
	Content interface{}            `json:"content"`
	OpDesc  string                 `json:"op_desc"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	SetID   int64                  `json:"inst_id"`
}

// AuditSetParams add set multiple log parameter
type AuditSetsParams struct {
	Content []auditoplog.AuditLogContext `json:"content"`
	OpDesc  string                       `json:"op_desc"`
	OpType  auditoplog.AuditOpType       `json:"op_type"`
}

type AuditQueryResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int            `json:"count"`
		Info  []OperationLog `json:"info"`
	} `json:"data"`
}
