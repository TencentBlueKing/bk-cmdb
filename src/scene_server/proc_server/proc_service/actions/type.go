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
 
package actions

type ProcessCResult struct {
	Result  bool           `json:"result"`
	Code    int            `json:"code"`
	Message interface{}    `json:"message"`
	Data    map[string]int `json:"data"`
}

type ProcessResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type ModuleSResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    ModuleData  `json:"data"`
}

type ModuleData struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

type ProcModuleConfig struct {
	ApplicationID int    `json:"bk_biz_id"`
	ModuleName    string `json:"bk_module_name"`
	processID     int    `json:"bk_process_id"`
}

type ProcModuleResult struct {
	Result  bool               `json:"result"`
	Code    int                `json:"code"`
	Message interface{}        `json:"message"`
	Data    []ProcModuleConfig `json:"data"`
}
