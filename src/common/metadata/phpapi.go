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

type UpdateHostParams struct {
	ProxyList []interface{} `json:"bk_proxy_list"`
	CloudID   int64         `json:"bk_cloud_id"`
}

type HostSearchByIPParams struct {
	IpList  []string `json:"ip_list"`
	CloudID *int64   `json:"bk_cloud_id"`
	AppID   []int64  `json:"bk_biz_id"`
}

//  HostSearchByAppIDParams host search by app
type HostSearchByAppIDParams struct {
	ApplicationID *int64 `json:"bk_biz_id"`
}

type HostSearchBySetIDParams struct {
	ApplicationID *int64  `json:"bk_biz_id"`
	SetID         []int64 `json:"bk_set_id"`
}

type HostSearchByModuleIDParams struct {
	ApplicationID *int64  `json:"bk_biz_id"`
	ModuleID      []int64 `json:"bk_module_id"`
}

// GetIPAndProxyByCompanyParams get id and proxy by company
type GetIPAndProxyByCompanyParams struct {
	Ips        []string `json:"ips"`
	AppIDStr   *string  `json:"bk_biz_id"`
	CloudIDStr *string  `json:"bk_cloud_id"`
}

type GetHostAppByCompanyIDParams struct {
	CompaynID  string `json:"bk_supplier_account"`
	IPs        string `json:"ip"`
	CloudIDStr string `json:"bk_cloud_id"`
}

var DelHostInAppParams struct {
	AppID  string `json:"appId"`
	HostID string `json:"hostId"`
}

type GitServerIpParams struct {
	AppName    string `json:"bk_biz_name"`
	SetName    string `json:"bk_set_name"`
	ModuleName string `json:"bk_module_name"`
}

type GetAgentStatusResult struct {
	AgentNorCnt    int                      `json:"agentNorCnt"`
	AgentAbnorCnt  int                      `json:"agentAbnorCnt"`
	AgentNorList   []map[string]interface{} `json:"agentNorList"`
	AgentAbnorList []map[string]interface{} `json:"agentAbnorList"`
}
