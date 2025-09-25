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

package metadata

// SysUserConfigItem 用户在cc_System表中的用户自定义配置
type SysUserConfigItem struct {
	Flag     bool  `json:"flag" bson:"flag"`
	ExpireAt int64 `json:"expire_at" bson:"expire_at"`
}

// ResponseSysUserConfigData response data for sys user config
type ResponseSysUserConfigData struct {
	RowType        string            `json:"type"`
	BluekingModify SysUserConfigItem `json:"blueking_modify"`
}

// ReponseSysUserConfig response for sys user config
type ReponseSysUserConfig struct {
	BaseResp `json:",inline"`
	Data     ResponseSysUserConfigData `json:"data"`
}

// CCSystemUserConfigSwitch TODO
const CCSystemUserConfigSwitch = "user_config_switch"
