/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package types

import (
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// GetBlueKingKit 后续改为从tenant包获取
func GetBlueKingKit() *rest.Kit {
	blueKingKit := rest.NewKit()
	blueKingKit.TenantID = GetBlueKing()
	return blueKingKit
}

// GetBlueKing 获取蓝鲸租户ID,后续修改为默认租户
func GetBlueKing() string {
	return common.BKDefaultTenantID
}

// MigrateResp migrate response
type MigrateResp struct {
	metadata.BaseResp `bson:",inline"`
	Data              *MigrateInfo `json:"data"`
}

// MigrateInfo migrate info
type MigrateInfo struct {
	PreVersion       string   `json:"pre_version"`
	CurrentVersion   string   `json:"current_version"`
	FinishedVersions []string `json:"finished_migrations"`
}
