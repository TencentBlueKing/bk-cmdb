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

import (
	"time"

	"configcenter/src/common/mapstr"
)

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
	Count uint64                   `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

type GetHostFavoriteResult struct {
	BaseResp `json:",inline"`
	Data     FavoriteResult `json:"data"`
}

type GetHostFavoriteWithIDResult struct {
	BaseResp `json:",inline"`
	Data     FavouriteMeta `json:"data"`
}

type HistoryContent struct {
	Content string `json:"content"`
}

type AddHistoryResult struct {
	BaseResp `json:",inline"`
	Data     ID `json:"data"`
}

type HistoryMeta struct {
	ID         string    `json:"id,omitempty" bson:"id,omitempty" `
	User       string    `json:"user,omitempty" bson:"user,omitempty"`
	Content    string    `json:"content,omitempty" bson:"content,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty" bson:"create_time,omitempty"`
	OwnerID    string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type HistoryResult struct {
	Count uint64        `json:"count"`
	Info  []HistoryMeta `json:"info"`
}

type GetHistoryResult struct {
	BaseResp `json:",inline"`
	Data     HistoryResult `json:"data"`
}

type HostInfo struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
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
	Data     []int64 `json:"data"`
}

type ParamData struct {
	ApplicationID       int64   `json:"bk_biz_id"`
	HostID              []int64 `json:"bk_host_id"`
	OwnerModuleID       int64   `json:"bk_owner_module_id"`
	OwnerAppplicationID int64   `json:"bk_owner_biz_id"`
}

type AssignHostToAppParams struct {
	ApplicationID      int64   `json:"bk_biz_id"`
	HostID             []int64 `json:"bk_host_id"`
	ModuleID           int64   `json:"bk_module_id"`
	OwnerApplicationID int64   `json:"bk_owner_biz_id"`
	OwnerModuleID      int64   `json:"bk_owner_module_id"`
}

type ModuleHost struct {
	AppID    int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	HostID   int64  `json:"bk_host_id" bson:"bk_host_id"`
	ModuleID int64  `json:"bk_module_id" bson:"bk_module_id"`
	SetID    int64  `json:"bk_set_id" bson:"bk_set_id"`
	OwnerID  string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type HostConfig struct {
	BaseResp `json:",inline"`
	Data     HostConfigData `json:"data"`
}

type HostConfigData struct {
	Count int64        `json:"count"`
	Info  []ModuleHost `json:"data"`
	Page  BasePage     `json:"page"`
}

type ModuleHostConfigParams struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostID        int64   `json:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id"`
	OwnerID       string  `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type UserConfig struct {
	Info       string    `json:"info" bson:"info"`
	Name       string    `json:"name" bson:"name"`
	ID         string    `json:"id" bson:"id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	UpdateTime time.Time `json:"last_time" bson:"last_time"`
	AppID      int64     `json:"bk_biz_id" bson:"bk_biz_id"`
	CreateUser string    `json:"create_user" bson:"create_user"`
	ModifyUser string    `json:"modify_user" bson:"modify_user"`
}

type UserConfigResult struct {
	Count uint64        `json:"count"`
	Info  []interface{} `json:"info"`
}

type GetUserConfigResult struct {
	BaseResp `json:",inline"`
	Data     UserConfigResult `json:"data"`
}

type GetUserCustomResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type FavouriteParms struct {
	ID          string `json:"id,omitempty"`
	Info        string `json:"info,omitempty"`
	QueryParams string `json:"query_params,omitempty"`
	Name        string `json:"name,omitempty"`
	IsDefault   int    `json:"is_default,omitempty"`
	Count       int    `json:"count,omitempty"`
	BizID       int64  `json:"bk_biz_id"`
}

type FavouriteMeta struct {
	BizID       int64     `json:"bk_biz_id" bson:"bk_biz_id"`
	ID          string    `json:"id,omitempty" bson:"id,omitempty"`
	Info        string    `json:"info,omitempty" bson:"info,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
	Count       int       `json:"count,omitempty" bson:"count,omitempty"`
	User        string    `json:"user,omitempty" bson:"user,omitempty"`
	OwnerID     string    `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account,omitempty"`
	QueryParams string    `json:"query_params,omitempty" bson:"query_params,omitempty"`
	CreateTime  time.Time `json:"create_time,omitempty" bson:"create_time,omitempty"`
	UpdateTime  time.Time `json:"last_time,omitempty" bson:"last_time,omitempty"`
}

type GetUserConfigDetailResult struct {
	BaseResp `json:",inline"`
	Data     UserConfigMeta `json:"data"`
}

type UserConfigMeta struct {
	AppID      int64     `json:"bk_biz_id,omitempty" bson:"bk_biz_id,omitempty"`
	Info       string    `json:"info,omitempty" bson:"info,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	ID         string    `json:"id,omitempty" bson:"id,omitempty"`
	CreateTime time.Time `json:"create_time" bson:"create_time,omitempty"`
	CreateUser string    `json:"create_user" bson:"create_user,omitempty"`
	ModifyUser string    `json:"modify_user" bson:"modify_user,omitempty"`
	UpdateTime time.Time `json:"last_time" bson:"last_time,omitempty"`
	OwnerID    string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type AddConfigQuery struct {
	AppID      int64  `json:"bk_biz_id,omitempty"`
	Info       string `json:"info,omitempty"`
	Name       string `json:"name,omitempty"`
	CreateUser string `json:"create_user,omitempty"`
}

type CloudTaskSearch struct {
	Count uint64          `json:"count"`
	Info  []CloudTaskInfo `json:"info"`
}

type CloudTaskInfo struct {
	User            string `json:"bk_user" bson:"bk_user"`
	TaskName        string `json:"bk_task_name" bson:"bk_task_name"`
	TaskID          int64  `json:"bk_task_id" bson:"bk_task_id"`
	AccountType     string `json:"bk_account_type" bson:"bk_account_type"`
	AccountAdmin    string `json:"bk_account_admin" bson:"bk_account_admin"`
	PeriodType      string `json:"bk_period_type" bson:"bk_period_type"`
	Period          string `json:"bk_period" bson:"bk_period"`
	LastSyncTime    string `json:"bk_last_sync_time" bson:"bk_last_sync_time"`
	ObjID           string `json:"bk_obj_id" bson:"bk_obj_id"`
	Status          bool   `json:"bk_status" bson:"bk_status"`
	ResourceConfirm bool   `json:"bk_confirm" bson:"bk_confirm"`
	AttrConfirm     bool   `json:"bk_attr_confirm" bson:"bk_attr_confirm"`
	SecretID        string `json:"bk_secret_id" bson:"bk_secret_id"`
	SecretKey       string `json:"bk_secret_key" bson:"bk_secret_key"`
	SyncStatus      string `json:"bk_sync_status" bson:"bk_sync_status"`
	NewAdd          int64  `json:"new_add" bson:"new_add"`
	AttrChanged     int64  `json:"attr_changed" bson:"attr_changed"`
	OwnerID         string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// TransferHostToInnerModule transfer host to inner module eg:idle module ,fault module
type TransferHostToInnerModule struct {
	ApplicationID int64   `json:"bk_biz_id"`
	ModuleID      int64   `json:"bk_module_id"`
	HostID        []int64 `json:"bk_host_id"`
}
