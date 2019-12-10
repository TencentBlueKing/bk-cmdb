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

type BkSystemInstallRequest struct {
	SetName    string                            `json:"bk_set_name"`
	ModuleName string                            `json:"bk_module_name"`
	InnerIP    string                            `json:"bk_host_innerip"`
	CloudID    int64                             `json:"bk_cloud_id"`
	HostInfo   map[string]interface{}            `json:"host_info"`
	ProcInfo   map[string]map[string]interface{} `json:"proc_info"`
}
