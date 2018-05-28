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

package params

type Gprivilege struct {
	ModelConfig    map[string]map[string][]string `json:"model_config"`
	SysConfig      SysConfigStruct                `json:"sys_config,omitempty"`
	IsHostCrossBiz bool                           `json:"is_host_cross_biz"`
}

type Privilege struct {
	ModelConfig map[string]map[string][]string `json:"model_config,omitempty"`
	SysConfig   *SysConfigStruct               `json:"sys_config,omitempty"`
}

type SysConfigStruct struct {
	Globalbusi []string `json:"global_busi"`
	BackConfig []string `json:"back_config"`
}

type UserPrivilege struct {
	GroupID     string                         `json:"bk_group_id"`
	ModelConfig map[string]map[string][]string `json:"model_config"`
	SysConfig   SysConfigStruct                `json:"sys_config"`
}

type UserPriviResult struct {
	Result  bool          `json:"result"`
	Code    int           `json:"code"`
	Message interface{}   `json:"message"`
	Data    UserPrivilege `json:"data"`
}

type GroupPrivilege struct {
	GroupID   string
	OwnerID   string
	Privilege Privilege `json:"privilege"`
}

type GroupPriviResult struct {
	Result  bool           `json:"result"`
	Code    int            `json:"code"`
	Message interface{}    `json:"message"`
	Data    GroupPrivilege `json:"data"`
}

type SearchGroup struct {
	Code    int         `json:"code"`
	Result  bool        `json:"result"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type SearchMainLine struct {
	Code    int                      `json:"code"`
	Result  bool                     `json:"result"`
	Message interface{}              `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}
