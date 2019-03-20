/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authcenter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/auth/meta"
)

func (ac *AuthCenter) Init(ctx context.Context, configs meta.InitConfig) error {
	header := http.Header{}
	if err := ac.authClient.RegistSystem(ctx, header, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	if err := ac.authClient.UpdateSystem(ctx, header, System{SystemID: expectSystem.SystemID, SystemName: expectSystem.SystemName}); err != nil {
		return err
	}

	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
		return err
	}

	// init model classifaction
	clsName2ID := map[string]int64{}
	for _, cls := range configs.Classifications {
		bizID, _ := cls.Metadata.Label.GetBusinessID()
		clsName2ID[fmt.Sprintf("%d:%s", bizID, cls.ClassificationID)] = cls.ID

		scopeType := ScopeTypeIDSystem
		ScopeID := SystemIDCMDB

		if bizID > 0 {
			scopeType = ScopeTypeIDBiz
			ScopeID = strconv.FormatInt(bizID, 10)
		}

		info := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysModelGroup,
					ResourceID: []ResourceID{
						{ResourceType: SysModelGroup, ResourceID: strconv.FormatInt(cls.ID, 10)},
					},
					ResourceName: cls.ClassificationName,
					ScopeInfo: ScopeInfo{
						ScopeType: scopeType,
						ScopeID:   ScopeID,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &info); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// init model description
	for _, model := range configs.Models {
		bizID, _ := model.Metadata.Label.GetBusinessID()
		clsID := clsName2ID[fmt.Sprintf("%d:%s", bizID, model.ObjCls)]
		if clsID <= 0 {
			return fmt.Errorf("classification id not found")
		}

		scopeType := ScopeTypeIDSystem
		ScopeID := SystemIDCMDB

		if bizID > 0 {
			scopeType = ScopeTypeIDBiz
			ScopeID = strconv.FormatInt(bizID, 10)
		}

		info := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysModelGroup,
					ResourceID: []ResourceID{
						{ResourceType: SysModelGroup, ResourceID: strconv.FormatInt(clsID, 10)},
						{ResourceType: SysModel, ResourceID: strconv.FormatInt(model.ID, 10)},
					},
					ResourceName: model.ObjectName,
					ScopeInfo: ScopeInfo{
						ScopeType: scopeType,
						ScopeID:   ScopeID,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &info); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// init business inst
	for _, biz := range configs.Bizs {
		bkbiz := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysBusinessInstance,
					ResourceID: []ResourceID{
						{ResourceType: SysBusinessInstance, ResourceID: strconv.FormatInt(biz.BizID, 10)},
					},
					ResourceName: biz.BizName,
					ScopeInfo: ScopeInfo{
						ScopeType: "system",
						ScopeID:   SystemIDCMDB,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &bkbiz); err != nil && err != ErrDuplicated {
			return err
		}
	}

	return nil
}
