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

type DeleteHostBatchOpt struct {
	HostID string `json:"bk_host_id"`
}

type FavouriteParms struct {
	ID          string `json:"id"`
	Info        string `json:"info"`
	QueryParams string `json:"query_params"`
	Name        string `json:"name"`
	IsDefault   int    `json:"is_default"`
	Count       int    `json:"count"`
}

type HostInstanceProperties struct {
	PropertyID    string      `json:"bk_property_id"`
	PropertyName  string      `json:"bk_property_name"`
	PropertyValue interface{} `json:"bk_property_value"`
}

type HostInstancePropertiesResult struct {
	BaseResp `json:",inline"`
	Data     []HostInstanceProperties `json:"data"`
}

type HostSnapResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type HostInputType string

const (
	ExecelType HostInputType = "excel"
)

type HostList struct {
	ApplicationID int64                            `json:"bk_biz_id"`
	HostInfo      map[int64]map[string]interface{} `json:"host_info"`
	SupplierID    int64                            `json:"bk_supplier_id"`
	InputType     HostInputType                    `json:"input_type"`
}

type HostsModuleRelation struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostID        []int64 `json:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id"`
	IsIncrement   bool    `json:"is_increment"`
}

type HostToAppModule struct {
	Ips         []string `json:"ips"`
	HostName    []string `json:"bk_host_name"`
	ModuleName  string   `json:"bk_module_name"`
	SetName     string   `json:"bk_set_name"`
	AppName     string   `json:"bk_biz_name"`
	OsType      string   `json:"bk_os_type"`
	OwnerID     string   `json:"bk_supplier_account"`
	PlatID      int64    `json:"bk_cloud_id"`
	IsIncrement bool     `json:"is_increment"`
}

type HostCommonSearch struct {
	AppID     int64             `json:"bk_biz_id,omitempty"`
	Ip        IPInfo            `json:"ip"`
	Condition []SearchCondition `json:"condition"`
	Page      BasePage          `json:"page"`
	Pattern   string            `json:"pattern,omitempty"`
}

//ip search info
type IPInfo struct {
	Data  []string `json:"data"`
	Exact int64    `json:"exact"`
	Flag  string   `json:"flag"`
}

//search condition
type SearchCondition struct {
	Fields    []string      `json:"fields"`
	Condition []interface{} `json:"condition"`
	ObjectID  string        `json:"bk_obj_id"`
}
