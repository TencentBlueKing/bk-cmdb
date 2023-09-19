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

package y3_11_202309050900

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

// deleteModuleHostApplyRule 根据 bk_module_id 删除 cc_HostApplyRule 表中所有脏数据
func deleteModuleHostApplyRule(ctx context.Context, db dal.RDB) error {
	// 已经处理过最大的id
	maxModuleID := int64(0)
	const pageSize uint64 = 100
	for {
		rules := make([]metadata.HostApplyRule, 0)
		ruleFilter := map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{common.BKDBGT: maxModuleID},
		}
		if err := db.Table(common.BKTableNameHostApplyRule).Find(ruleFilter).Fields(common.BKModuleIDField).
			Limit(pageSize).Sort(common.BKModuleIDField).All(ctx, &rules); err != nil {
			blog.Errorf("list module ids failed, db find failed, err: %v", err)
			return err
		}
		if len(rules) == 0 {
			break
		}
		maxModuleID = rules[len(rules)-1].ModuleID

		moduleIDs := make([]int64, 0)
		for _, rule := range rules {
			moduleIDs = append(moduleIDs, rule.ModuleID)
		}
		moduleIDs = util.IntArrayUnique(moduleIDs)

		moduleField := map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
		}
		moduleResult := make([]metadata.ModuleInst, 0)
		if err := db.Table(common.BKTableNameBaseModule).Find(moduleField).Fields(common.BKModuleIDField).All(ctx,
			&moduleResult); err != nil {
			blog.Errorf("list module ids failed, db find failed, err: %v", err)
			return err
		}
		dbModuleIDs := make([]int64, 0)
		for _, module := range moduleResult {
			dbModuleIDs = append(dbModuleIDs, module.ModuleID)
		}

		deleteIDs := util.IntArrDiff(moduleIDs, dbModuleIDs)
		if len(deleteIDs) == 0 {
			continue
		}
		// 执行删除
		deleteFilter := map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: deleteIDs,
			},
		}
		if err := db.Table(common.BKTableNameHostApplyRule).Delete(ctx, deleteFilter); err != nil {
			blog.Errorf("delete host apply rules failed, err: %v", err)
			return err
		}
	}
	return nil
}

// deleteServiceTemplateHostApplyRule 根据 service_template_id 删除 cc_HostApplyRule 表中所有脏数据
func deleteServiceTemplateHostApplyRule(ctx context.Context, db dal.RDB) error {
	maxSrvTmplID := int64(0)
	const pageSize uint64 = 100
	for {
		rules := make([]metadata.HostApplyRule, 0)
		ruleFilter := map[string]interface{}{
			common.BKServiceTemplateIDField: map[string]interface{}{common.BKDBGT: maxSrvTmplID},
		}
		if err := db.Table(common.BKTableNameHostApplyRule).Find(ruleFilter).Fields(common.BKServiceTemplateIDField).
			Limit(pageSize).Sort(common.BKServiceTemplateIDField).All(ctx, &rules); err != nil {
			blog.Errorf("list service template ids failed, db find failed, err: %v", err)
			return err
		}
		if len(rules) == 0 {
			break
		}
		maxSrvTmplID = rules[len(rules)-1].ServiceTemplateID

		srvTmplIDs := make([]int64, 0)
		for _, rule := range rules {
			srvTmplIDs = append(srvTmplIDs, rule.ServiceTemplateID)
		}
		srvTmplIDs = util.IntArrayUnique(srvTmplIDs)

		serviceTemplateField := map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: srvTmplIDs,
			},
		}
		serviceTemplateResult := make([]metadata.ServiceTemplate, 0)
		if err := db.Table(common.BKTableNameServiceTemplate).Find(serviceTemplateField).Fields(common.BKFieldID).All(ctx,
			&serviceTemplateResult); err != nil {
			blog.Errorf("list service template ids failed, db find failed, err: %v", err)
			return err
		}
		dbSrvTmplIDs := make([]int64, 0)
		for _, serviceTemplate := range serviceTemplateResult {
			dbSrvTmplIDs = append(dbSrvTmplIDs, serviceTemplate.ID)
		}

		deleteIDs := util.IntArrDiff(srvTmplIDs, dbSrvTmplIDs)
		if len(deleteIDs) == 0 {
			continue
		}
		// 执行删除
		deleteFilter := map[string]interface{}{
			common.BKServiceTemplateIDField: map[string]interface{}{
				common.BKDBIN: deleteIDs,
			},
		}
		if err := db.Table(common.BKTableNameHostApplyRule).Delete(ctx, deleteFilter); err != nil {
			blog.Errorf("delete host apply rules failed, err: %v", err)
			return err
		}
	}
	return nil
}
