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

package y3_11_202308101200

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
)

// deleteModuleHostApplyRule 根据 bk_module_id 删除 cc_HostApplyRule 表中所有脏数据的
func deleteModuleHostApplyRule(ctx context.Context, db dal.RDB) error {
	moduleIds, err := listDirtyDataModuleIds(ctx, db)
	if err != nil {
		blog.Errorf("ListDirtyDataModuleIds failed, err: %v", err)
		return err
	}

	filter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIds,
		},
	}

	err = db.Table(common.BKTableNameHostApplyRule).Delete(ctx, filter)
	if err != nil {
		blog.Errorf("delete HostApplyRule failed, err: %v", err)
		return err
	}

	return nil
}

// listDirtyDataModuleIds 获取 cc_HostApplyRule 表中所有脏数据的 bk_module_id
func listDirtyDataModuleIds(ctx context.Context, db dal.RDB) ([]int64, error) {
	var result []int64
	ruleIds := make([]map[string]int64, 0)

	filter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBGT: 0,
		},
	}
	err := db.Table(common.BKTableNameHostApplyRule).Find(filter).Fields(common.BKModuleIDField).All(ctx, &ruleIds)
	if err != nil {
		blog.Errorf("ListModuleIds failed, db find failed, err: %s", err.Error())
		return result, err
	}

	for _, ruleId := range ruleIds {
		filter := map[string]interface{}{
			common.BKModuleIDField: map[string]int64{
				common.BKDBEQ: ruleId[common.BKModuleIDField],
			},
		}
		count, err := db.Table(common.BKTableNameBaseModule).Find(filter).Count(ctx)
		if err != nil {
			blog.Errorf("db count failed, err: %s", err.Error())
			return result, err
		}
		if count == 0 {
			result = append(result, ruleId[common.BKModuleIDField])
		}
	}

	return result, nil
}

// deleteServiceTemplateHostApplyRule 根据 service_template_id 删除 cc_HostApplyRule 表中所有脏数据的
func deleteServiceTemplateHostApplyRule(ctx context.Context, db dal.RDB) error {

	moduleIds, err := listDirtyDataServiceTemplateIds(ctx, db)
	if err != nil {
		blog.Errorf("ListDirtyDataModuleIds failed, err: %v", err)
		return err
	}

	filter := map[string]interface{}{
		common.BKServiceTemplateIDField: map[string]interface{}{
			common.BKDBIN: moduleIds,
		},
	}

	err = db.Table(common.BKTableNameHostApplyRule).Delete(ctx, filter)
	if err != nil {
		blog.Errorf("delete HostApplyRule failed, err: %v", err)
		return err
	}

	return nil
}

// listDirtyDataServiceTemplateIds 获取 cc_HostApplyRule 表中所有脏数据的 service_template_id
func listDirtyDataServiceTemplateIds(ctx context.Context, db dal.RDB) ([]int64, error) {
	var result []int64
	ruleIds := make([]map[string]int64, 0)

	filter := map[string]interface{}{
		common.BKServiceTemplateIDField: map[string]interface{}{
			common.BKDBGT: 0,
		},
	}
	err := db.Table(common.BKTableNameHostApplyRule).Find(filter).
		Fields(common.BKServiceTemplateIDField).All(ctx, &ruleIds)
	if err != nil {
		blog.Errorf("ListModuleIds failed, db find failed, err: %s", err.Error())
		return result, err
	}

	for _, ruleId := range ruleIds {
		filter := map[string]interface{}{
			common.BKFieldID: map[string]int64{
				common.BKDBEQ: ruleId[common.BKServiceTemplateIDField],
			},
		}
		count, err := db.Table(common.BKTableNameServiceTemplate).Find(filter).Count(ctx)
		if err != nil {
			blog.Errorf("db count failed, err: %s", err.Error())
			return result, err
		}
		if count == 0 {
			result = append(result, ruleId[common.BKServiceTemplateIDField])
		}
	}

	return result, nil
}
