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

package util

import (
	"context"
	"net/http"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
)

// Kit is cmdb data syncer kit
type Kit struct {
	Rid    string
	Header http.Header
	Ctx    context.Context
}

// NewKit creates a new kit for cmdb data syncer
func NewKit() *Kit {
	rid := util.GenerateRID()

	header := make(http.Header)

	httpheader.SetRid(header, rid)
	httpheader.SetUser(header, common.CCSystemOperatorUserName)
	httpheader.SetTenantID(header, common.BKDefaultTenantID)
	header.Add("Content-Type", "application/json")

	ctx := util.NewContextFromHTTPHeader(header)
	ctx, header = util.SetReadPreference(context.Background(), header, common.SecondaryPreferredMode)

	return &Kit{
		Rid:    rid,
		Header: header,
		Ctx:    ctx,
	}
}

// ConvertKit converts rest.Kit to Kit
func ConvertKit(kit *rest.Kit) *Kit {
	ctx, header := util.SetReadPreference(kit.Ctx, kit.Header, common.SecondaryPreferredMode)
	return &Kit{
		Rid:    kit.Rid,
		Header: header,
		Ctx:    ctx,
	}
}
