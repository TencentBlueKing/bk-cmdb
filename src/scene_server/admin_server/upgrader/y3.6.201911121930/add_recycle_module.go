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

package y3_6_201911121930

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type BizSimplify struct {
	BKAppIDField   int64  `field:"bk_biz_id" bson:"bk_biz_id"`
	BKAppNameField string `field:"bk_biz_name" bson:"bk_biz_name"`
}

func AddRecycleModule(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	defaultServiceCategoryID, err := getDefaultServiceCategoryID(ctx, db, conf)
	if err != nil {
		blog.Errorf("AddRecycleModule failed, getDefaultServiceCategoryID failed, err: %s", err.Error())
		return fmt.Errorf("getDefaultServiceCategoryID failed, err: %s", err.Error())
	}

	// 资源池业务也有空闲机模块
	bizFilter := map[string]interface{}{}
	businessList := make([]BizSimplify, 0)

	start := uint64(0)
	limit := uint64(50)
	for {
		if err := db.Table(common.BKTableNameBaseApp).Find(bizFilter).Start(start).Limit(limit).All(ctx, &businessList); err != nil {
			blog.ErrorJSON("AddRecycleModule failed, find businesses failed, filter: %s, err: %s", bizFilter, err)
			return fmt.Errorf("find businesses failed, err: %s", err.Error())
		}

		if len(businessList) == 0 {
			break
		}

		for _, biz := range businessList {
			if err := ensureRecycleModuleForBiz(ctx, db, conf, biz, defaultServiceCategoryID); err != nil {
				blog.ErrorJSON("AddRecycleModule failed, ensureRecycleModuleForBiz failed, biz: %s, defaultServiceCategoryID: %s, err: %s", biz, defaultServiceCategoryID, err)
				return fmt.Errorf("ensureRecycleModuleForBiz failed, err: %s", err.Error())
			}
		}

		start += limit
	}

	return nil
}

func getDefaultServiceCategoryID(ctx context.Context, db dal.RDB, conf *upgrader.Config) (int64, error) {
	serviceCategory := metadata.ServiceCategory{}
	filter := map[string]interface{}{
		"is_built_in": true,
		"name":        "Default",
		"bk_parent_id": map[string]interface{}{
			common.BKDBGT: 0,
		},
	}
	if err := db.Table(common.BKTableNameServiceCategory).Find(filter).One(ctx, &serviceCategory); err != nil {
		blog.ErrorJSON("getDefaultServiceCategoryID failed, find service category failed, filter: %s, err: %s", filter, err.Error())
		return 0, fmt.Errorf("get default service category failed, err: %s", err)
	}
	return serviceCategory.ID, nil
}

func ensureRecycleModuleForBiz(ctx context.Context, db dal.RDB, conf *upgrader.Config, biz BizSimplify, defaultServiceCategoryID int64) error {
	moduleFilter := map[string]interface{}{
		common.BKDefaultField: common.DefaultRecycleModuleFlag,
		common.BKAppIDField:   biz.BKAppIDField,
	}
	count, err := db.Table(common.BKTableNameBaseModule).Find(moduleFilter).Count(ctx)
	if err != nil {
		blog.ErrorJSON("ensureRecycleModuleForBiz failed, moduleFilter: %s, err: %s", moduleFilter, err.Error())
		return fmt.Errorf("count default module failed, bizID: %d, err: %s", biz.BKAppIDField, err.Error())
	}
	if count > 0 {
		return nil
	}

	resourcePoolSetID, err := getResourceSetID(ctx, db, conf, biz)
	if err != nil {
		blog.Errorf("ensureRecycleModuleForBiz failed, getResourceSetID failed, bizID: %s, err: %s", biz.BKAppIDField, err.Error())
		return fmt.Errorf("get resource pool set failed, bizID: %d, err: %s", biz.BKAppIDField, err.Error())
	}

	newModuleID, err := db.NextSequence(ctx, common.BKTableNameBaseModule)
	if err != nil {
		blog.ErrorJSON("ensureRecycleModuleForBiz failed, moduleFilter: %s, err: %s", moduleFilter, err.Error())
		return fmt.Errorf("count default module failed, bizID: %d, err: %s", biz.BKAppIDField, err.Error())
	}

	// add recycle module
	now := time.Now()
	doc := map[string]interface{}{
		common.BKSetIDField:             resourcePoolSetID,
		common.BKParentIDField:          resourcePoolSetID,
		common.BkSupplierAccount:        conf.OwnerID,
		common.CreateTimeField:          now,
		common.LastTimeField:            now,
		common.BKChildStr:               nil,
		common.BKModuleTypeField:        common.DefaultModuleType,
		common.BKOperatorField:          "",
		common.BKAppIDField:             biz.BKAppIDField,
		common.BKModuleNameField:        common.DefaultRecycleModuleName,
		common.BKDefaultField:           common.DefaultRecycleModuleFlag,
		common.BKBakOperatorField:       "",
		common.BKModuleIDField:          newModuleID,
		common.BKServiceCategoryIDField: defaultServiceCategoryID,
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
		common.BKSetTemplateIDField:     common.SetTemplateIDNotSet,
	}
	if err := db.Table(common.BKTableNameBaseModule).Insert(ctx, doc); err != nil {
		blog.ErrorJSON("ensureRecycleModuleForBiz failed, doc: %s, err: %s", doc, err.Error())
		return fmt.Errorf("insert module failed, err: %s", err)
	}
	return nil
}

func getResourceSetID(ctx context.Context, db dal.RDB, conf *upgrader.Config, biz BizSimplify) (int64, error) {
	setFilter := map[string]interface{}{
		common.BKDefaultField: common.DefaultResSetFlag,
		common.BKAppIDField:   biz.BKAppIDField,
	}
	set := struct {
		BizID   int64  `bson:"bk_biz_id" json:"bk_biz_id" mapstructure:"bk_biz_id"`
		SetID   int64  `bson:"bk_set_id" json:"bk_set_id" mapstructure:"bk_set_id"`
		SetName string `bson:"bk_set_name" json:"bk_set_name" mapstructure:"bk_set_name"`
	}{}
	err := db.Table(common.BKTableNameBaseSet).Find(setFilter).One(ctx, &set)
	if err != nil {
		blog.ErrorJSON("getResourceSetID failed, setFilter: %s, err: %s", setFilter, err.Error())
		return 0, fmt.Errorf("get resource pool set failed, bizID: %d, err: %s", biz.BKAppIDField, err.Error())
	}
	return set.SetID, nil
}
