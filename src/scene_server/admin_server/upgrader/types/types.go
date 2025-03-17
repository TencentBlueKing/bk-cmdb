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
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// GetBlueKingKit 后续改为从tenant包获取
func GetBlueKingKit() *rest.Kit {
	var err error
	blueKingKit := rest.NewKit()
	blueKingKit.TenantID, err = GetBlueKing()
	if err != nil {
		blog.Errorf("get tenant id failed, err: %v", err)
		return nil
	}
	return blueKingKit
}

// GetBlueKing 获取蓝鲸租户ID,后续修改为默认租户
func GetBlueKing() (string, error) {
	tenantModeEnable, err := cc.Bool("tenant.enableMultiTenantMode")
	if err != nil {
		blog.Errorf("get enable multi tenant mode failed, err: %v", err)
		return "", err
	}
	if tenantModeEnable {
		return common.BKDefaultTenantID, nil
	}
	return common.BKUnconfiguredTenantID, nil
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
