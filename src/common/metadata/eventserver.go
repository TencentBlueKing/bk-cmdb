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

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
)

// SyncIdentifierResult sync host identifier result
type SyncIdentifierResult struct {
	TaskID    string          `json:"task_id"`
	HostInfos []HostBriefInfo `json:"host_infos"`
}

// HostBriefInfo get sync host identifier task result option
type HostBriefInfo struct {
	HostID         int64  `json:"bk_host_id"`
	Identification string `json:"identification"`
}

// GetTaskResultOption get sync host identifier task result option
type GetTaskResultOption struct {
	TaskID string `json:"task_id"`
}

// Validate validate GetTaskResultOption
func (h *GetTaskResultOption) Validate() (rawError errors.RawErrorInfo) {
	if h.TaskID == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"task_id"},
		}
	}

	return errors.RawErrorInfo{}
}

// HostIdentifierTaskResult sync host identifier task result
type HostIdentifierTaskResult struct {
	SuccessList []int64 `json:"success_list"`
	FailedList  []int64 `json:"failed_list"`
	PendingList []int64 `json:"pending_list"`
}
