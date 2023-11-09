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

package fulltextsearch

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
)

// SyncDataPageSize is the size of one sync data operation page
const SyncDataPageSize = 500

// SyncDataOption defines the sync full-text search data options
type SyncDataOption struct {
	// IsAll defines if sync all data
	IsAll bool `json:"is_all"`
	// Index defines which index's data to sync
	Index string `json:"index"`
	// Collection defines which collection's data to sync
	Collection string `json:"collection"`
	// Oids defines the specific oids of data to sync in collection, it must be set with the Collection field
	Oids []string `json:"oids"`
}

// Validate sync data options
func (o *SyncDataOption) Validate() errors.RawErrorInfo {
	if o == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"option"}}
	}

	if o.IsAll {
		if len(o.Collection) > 0 || len(o.Index) > 0 || len(o.Oids) > 0 {
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
				Args: []interface{}{"only one of the sync options can be set"}}
		}

		return errors.RawErrorInfo{}
	}

	if len(o.Index) > 0 {
		if len(o.Collection) > 0 || len(o.Oids) > 0 {
			return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
				Args: []interface{}{"only one of the sync options can be set"}}
		}
		return errors.RawErrorInfo{}
	}

	if len(o.Collection) == 0 || len(o.Oids) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid,
			Args: []interface{}{"one of the sync options must be set"}}
	}

	if len(o.Oids) > SyncDataPageSize {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit,
			Args: []interface{}{"ids length", SyncDataPageSize}}
	}

	return errors.RawErrorInfo{}
}

// MigrateResult defines the sync full-text search migrate result
type MigrateResult struct {
	PreVersion       int   `json:"pre_version,omitempty"`
	CurrentVersion   int   `json:"current_version,omitempty"`
	FinishedVersions []int `json:"finished_migrations,omitempty"`

	Message string `json:"message,omitempty"`
}
