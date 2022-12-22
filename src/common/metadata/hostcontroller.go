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

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
)

// ID TODO
type ID struct {
	ID string `json:"id"`
}

// IDResult TODO
type IDResult struct {
	BaseResp `json:",inline"`
	Data     ID `json:"data"`
}

// HostInstanceResult TODO
type HostInstanceResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

// FavoriteResult TODO
type FavoriteResult struct {
	Count uint64                   `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

// GetHostFavoriteResult TODO
type GetHostFavoriteResult struct {
	BaseResp `json:",inline"`
	Data     FavoriteResult `json:"data"`
}

// GetHostFavoriteWithIDResult TODO
type GetHostFavoriteWithIDResult struct {
	BaseResp `json:",inline"`
	Data     FavouriteMeta `json:"data"`
}

// HistoryContent TODO
type HistoryContent struct {
	Content string `json:"content"`
}

// AddHistoryResult TODO
type AddHistoryResult struct {
	BaseResp `json:",inline"`
	Data     ID `json:"data"`
}

// HistoryMeta TODO
type HistoryMeta struct {
	ID         string    `json:"id,omitempty" bson:"id,omitempty" `
	User       string    `json:"user,omitempty" bson:"user,omitempty"`
	Content    string    `json:"content,omitempty" bson:"content,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty" bson:"create_time,omitempty"`
	OwnerID    string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// HistoryResult TODO
type HistoryResult struct {
	Count uint64        `json:"count"`
	Info  []HistoryMeta `json:"info"`
}

// GetHistoryResult TODO
type GetHistoryResult struct {
	BaseResp `json:",inline"`
	Data     HistoryResult `json:"data"`
}

// HostInfo TODO
type HostInfo struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

// ExtractHostIDs extract hostIDs
func (h HostInfo) ExtractHostIDs() ([]int64, error) {
	ids := make([]int64, len(h.Info))
	for idx, h := range h.Info {
		id, err := h.Int64(common.BKHostIDField)
		if err != nil {
			return nil, err
		}
		ids[idx] = id
	}
	return ids, nil
}

// GetHostsResult TODO
type GetHostsResult struct {
	BaseResp `json:",inline"`
	Data     HostInfo `json:"data"`
}

// GetHostModuleIDsResult TODO
type GetHostModuleIDsResult struct {
	BaseResp `json:",inline"`
	Data     []int64 `json:"data"`
}

// ParamData TODO
type ParamData struct {
	ApplicationID       int64   `json:"bk_biz_id"`
	HostID              []int64 `json:"bk_host_id"`
	OwnerModuleID       int64   `json:"bk_owner_module_id"`
	OwnerAppplicationID int64   `json:"bk_owner_biz_id"`
}

// AssignHostToAppParams TODO
type AssignHostToAppParams struct {
	ApplicationID      int64   `json:"bk_biz_id"`
	HostID             []int64 `json:"bk_host_id"`
	ModuleID           int64   `json:"bk_module_id"`
	OwnerApplicationID int64   `json:"bk_owner_biz_id"`
	OwnerModuleID      int64   `json:"bk_owner_module_id"`
}

// ModuleHost TODO
type ModuleHost struct {
	AppID    int64  `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
	HostID   int64  `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	ModuleID int64  `json:"bk_module_id,omitempty" bson:"bk_module_id"`
	SetID    int64  `json:"bk_set_id,omitempty" bson:"bk_set_id"`
	OwnerID  string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

// HostConfig TODO
type HostConfig struct {
	BaseResp `json:",inline"`
	Data     HostConfigData `json:"data"`
}

// HostConfigData TODO
type HostConfigData struct {
	Count int64        `json:"count"`
	Info  []ModuleHost `json:"data"`
	Page  BasePage     `json:"page"`
}

// HostModuleResp TODO
type HostModuleResp struct {
	BaseResp `json:",inline"`
	Data     []ModuleHost `json:"data"`
}

// ModuleHostConfigParams TODO
type ModuleHostConfigParams struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostID        int64   `json:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id"`
	OwnerID       string  `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// GetUserCustomResult TODO
type GetUserCustomResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

// FavouriteType host query favorite condition type
type FavouriteType string

const (
	// Container container topo type about host query favorite condition
	Container FavouriteType = "container"
	// Tradition tradition topo type about host query favorite condition
	Tradition FavouriteType = "tradition"
)

// FavouriteParms host query favorite condition parameter
type FavouriteParms struct {
	ID          string        `json:"id,omitempty"`
	Info        string        `json:"info,omitempty"`
	QueryParams string        `json:"query_params,omitempty"`
	Name        string        `json:"name,omitempty"`
	IsDefault   int           `json:"is_default,omitempty"`
	Count       int           `json:"count,omitempty"`
	BizID       int64         `json:"bk_biz_id"`
	Type        FavouriteType `json:"type"`
}

// FavouriteMeta host query favorite condition metadata
type FavouriteMeta struct {
	BizID       int64         `json:"bk_biz_id" bson:"bk_biz_id"`
	ID          string        `json:"id,omitempty" bson:"id,omitempty"`
	Info        string        `json:"info,omitempty" bson:"info,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Count       int           `json:"count,omitempty" bson:"count,omitempty"`
	User        string        `json:"user,omitempty" bson:"user,omitempty"`
	OwnerID     string        `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account,omitempty"`
	Type        FavouriteType `json:"type,omitempty" bson:"type,omitempty"`
	QueryParams string        `json:"query_params,omitempty" bson:"query_params,omitempty"`
	CreateTime  time.Time     `json:"create_time,omitempty" bson:"create_time,omitempty"`
	UpdateTime  time.Time     `json:"last_time,omitempty" bson:"last_time,omitempty"`
}

// TransferHostToInnerModule transfer host to inner module eg:idle module ,fault module
type TransferHostToInnerModule struct {
	ApplicationID int64   `json:"bk_biz_id"`
	ModuleID      int64   `json:"bk_module_id"`
	HostID        []int64 `json:"bk_host_id"`
}

// DistinctIDResponse TODO
type DistinctIDResponse struct {
	BaseResp `json:",inline"`
	Data     DistinctID `json:"data"`
}

// DistinctID TODO
type DistinctID struct {
	IDArr []int64 `json:"id_arr"`
}
