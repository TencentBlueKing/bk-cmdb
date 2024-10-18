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

// Package types defines multiple cmdb synchronize types
package types

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
)

// SyncCmdbDataOption defines sync cmdb data option
type SyncCmdbDataOption struct {
	ResType ResType          `json:"resource_type"`
	SubRes  string           `json:"sub_resource"`
	IsAll   bool             `json:"is_all"`
	Start   map[string]int64 `json:"start"`
	End     map[string]int64 `json:"end"`
}

// Validate sync cmdb data option
func (o *SyncCmdbDataOption) Validate() errors.RawErrorInfo {
	if rawErr := o.ResType.Validate(o.SubRes); rawErr.ErrCode != 0 {
		return rawErr
	}

	if o.IsAll {
		if len(o.Start) != 0 || len(o.End) != 0 {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{"is_all", "start", "end"},
			}
		}
		return errors.RawErrorInfo{}
	}

	if len(o.Start) == 0 && len(o.End) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"start", "end"},
		}
	}

	return errors.RawErrorInfo{}
}

// InfiniteEndID represent infinity for end id of id rule info
const InfiniteEndID int64 = -1
