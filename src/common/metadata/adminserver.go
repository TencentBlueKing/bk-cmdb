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

// SysUserConfigItem 用户在cc_System表中的用户自定义配置
type SysUserConfigItem struct {
	Flag     bool  `json:"flag" bson:"flag"`
	ExpireAt int64 `json:"expire_at" bson:"expire_at"`
}

type ResponseSysUserConfigData struct {
	RowType        string            `json:"type"`
	BluekingModify SysUserConfigItem `json:"blueking_modify"`
}

type ReponseSysUserConfig struct {
	BaseResp `json:",inline"`
	Data     ResponseSysUserConfigData `json:"data"`
}

const CCSystemUserConfigSwitch = "user_config_switch"
