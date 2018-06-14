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

import "time"

type ID struct {
	ID string `json:"id"`
}

type IDResult struct {
	BaseResp `json:",inline"`
	Data     ID `json:"data"`
}

type HostInstanceResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type FavoriteResult struct {
	Count int           `json:"count"`
	Info  []interface{} `json:"info"`
}

type GetHostFavoriteResult struct {
	BaseResp `json:",inline"`
	Data     FavoriteResult `json:"data"`
}

type GetHostFavoriteWithIDResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type HistoryContent struct {
	Content string `json:"content"`
}

type AddHistoryResult struct {
	BaseResp `json:",inline"`
	Data     ID `json:"data"`
}

type HistoryResult struct {
	Count int           `json:"count"`
	Info  []interface{} `json:"info"`
}

type GetHistoryResult struct {
	BaseResp `json:",inline"`
	Data     HistoryResult `json:"data"`
}

type HostInfo struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

type GetHostsResult struct {
	BaseResp `json:",inline"`
	Data     HostInfo `json:"data"`
}

type HostSnap struct {
	Data string `json:"data"`
}

type GetHostSnapResult struct {
	BaseResp `json:",inline"`
	Data     HostSnap `json:"data"`
}

type GetHostModuleIDsResult struct {
	BaseResp `json:",inline"`
	Data     []int `json:"data"`
}

type ParamData struct {
	ApplicationID       int   `json:"bk_biz_id"`
	HostID              []int `json:"bk_host_id"`
	OwnerModuleID       int   `json:"bk_owner_module_id"`
	OwnerAppplicationID int   `json:"bk_owner_biz_id"`
}

type AssignHostToAppParams struct {
	ApplicationID      int   `json:"bk_biz_id"`
	HostID             []int `json:"bk_host_id"`
	ModuleID           int   `json:"bk_module_id"`
	OwnerApplicationID int   `json:"bk_owner_biz_id"`
	OwnerModuleID      int   `json:"bk_owner_module_id"`
}

type HostConfig struct {
	BaseResp `json:",inline"`
	Data     []interface{} `json:"data"`
}

type ModuleHostConfigParams struct {
	ApplicationID int64 `json:"bk_biz_id"`
	HostID        int64 `json:"bk_host_id"`
	ModuleID      []int `json:"bk_module_id"`
}

type UserConfig struct {
	Info       string    `json:"info"`
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"last_time"`
	AppID      int64     `json:"bk_biz_id"`
	CreateUser string    `json:"create_user"`
	ModifyUser string    `json:"modify_user"`
}

type UserConfigResult struct {
	Count int           `json:"count"`
	Info  []interface{} `json:"info"`
}

type GetUserConfigResult struct {
	BaseResp `json:",inline"`
	Data     UserConfigResult `json:"data"`
}

type GetUserConfigDetailResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type GetUserCustomResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}
