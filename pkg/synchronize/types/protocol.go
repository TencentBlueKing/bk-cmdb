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
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/errors"
)

// CreateSyncDataOption defines create sync data option
type CreateSyncDataOption struct {
	ResourceType ResType           `json:"resource_type"`
	SubResource  string            `json:"sub_resource"`
	Data         []json.RawMessage `json:"data"`
}

// Validate create sync data option
func (o *CreateSyncDataOption) Validate() errors.RawErrorInfo {
	if rawErr := o.ResourceType.Validate(o.SubResource); rawErr.ErrCode != 0 {
		return rawErr
	}

	if len(o.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(o.Data) > common.BKMaxLimitSize {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", common.BKMaxLimitSize},
		}
	}

	return errors.RawErrorInfo{}
}
