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

package auditoplog

// AuditOpType  log type
type AuditOpType int

const (
	// AuditOpTypeAdd  object add log
	AuditOpTypeAdd AuditOpType = iota + 1
	// AuditOpTypeModify object add log
	AuditOpTypeModify AuditOpType = 2
	// AuditOpTypeDel object del log
	AuditOpTypeDel AuditOpType = 3
	// AuditOpTypeHostModule host  change module
	AuditOpTypeHostModule AuditOpType = 100
)
