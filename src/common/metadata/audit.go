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

type SaveAuditLogParams struct {
	ID      int64                  `json:"inst_id"`
	Model   string                 `json:"op_target"`
	Content interface{}            `json:"content"`
	ExtKey  string                 `json:"ext"`
	OpDesc  string                 `json:"op_desc"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	BizID   int64                  `json:"biz_id"`
}

// AuditQueryResult add single host log paramm
type AuditQueryResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int            `json:"count"`
		Info  []OperationLog `json:"info"`
	} `json:"data"`
}
