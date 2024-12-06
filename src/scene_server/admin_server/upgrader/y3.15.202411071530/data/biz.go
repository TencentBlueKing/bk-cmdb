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

package data

import (
	"fmt"
	"sort"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/scene_server/admin_server/upgrader/types"
	"configcenter/src/storage/dal"
)

func addBizData(kit *rest.Kit, db dal.Dal) error {

	if kit.TenantID == types.GetBlueKing() {
		if err := addBizModule(kit, db, []interface{}{blueKingBizData, resBizData},
			[]tools.AuditField{blueKingBizAudit, resBizAudit}); err != nil {
			blog.Errorf("add biz module or set data failed, %v", err)
			return err
		}
		return nil
	}
	if err := addBizModule(kit, db, []interface{}{resBizData}, []tools.AuditField{resBizAudit}); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}
	return nil

}

func addBizModule(kit *rest.Kit, db dal.Dal, data []interface{}, auditField []tools.AuditField) error {

	cmpField := &tools.CmpFiled{
		Unique:     []string{common.BKAppNameField},
		IgnoreKeys: []string{common.BKAppIDField},
		ID:         common.BKAppIDField,
	}
	ids, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameBaseApp, data, cmpField, auditField)
	if err != nil {
		blog.Errorf("insert biz data for table %s failed, err: %v, data: %+v", common.BKTableNameBaseApp, err, data)
		return err
	}
	if len(ids) == 0 {
		blog.Errorf("failed to get biz id")
		return fmt.Errorf("failed to get biz id")
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	bkBizID := ids[0]
	resBizID := ids[len(ids)-1]
	resModuleName := []string{common.DefaultResModuleName}
	if err := addBizAsstData(kit, db, int64(resBizID), resModuleName); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}

	if kit.TenantID != types.GetBlueKing() {
		return nil
	}

	bkModuleName := []string{common.DefaultResModuleName, common.DefaultFaultModuleName,
		common.DefaultRecycleModuleName}
	if err := addBizAsstData(kit, db, int64(bkBizID), bkModuleName); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}
	return nil
}

func addBizAsstData(kit *rest.Kit, db dal.Dal, bizID int64, moduleNames []string) error {
	// add resource business cluster
	ids, err := addSetBaseData(kit, db, bizID)
	if err != nil {
		blog.Errorf("add set bizSetData failed, %v", err)
		return err
	}
	if len(ids) == 0 {
		blog.Errorf("failed to get set id")
		return fmt.Errorf("failed to get set id")
	}
	// add resource business module
	if err := addModuleData(kit, db, bizID, moduleNames, ids[0]); err != nil {
		blog.Errorf("add module bizSetData failed, %v", err)
		return err
	}
	return nil
}

var (
	blueKingBizData = bizData{
		BizName:       common.BKAppName,
		BizMaintainer: "admin",
		TimeZone:      "Asia/Shanghai",
		Language:      "1",
		LifeCycle:     common.DefaultAppLifeCycleNormal,
		Default:       common.DefaultFlagDefaultValue,
	}
	resBizData = bizData{
		BizName:       common.DefaultAppName,
		BizProductor:  "admin",
		BizMaintainer: "admin",
		TimeZone:      "Asia/Shanghai",
		Language:      "1",
		LifeCycle:     common.DefaultAppLifeCycleNormal,
		Default:       common.DefaultAppFlag,
		BizID:         1,
	}
	blueKingBizAudit = tools.AuditField{
		AuditType:    metadata.BusinessType,
		ResourceType: metadata.BusinessRes,
		ResourceName: common.BKAppName,
	}
	resBizAudit = tools.AuditField{
		AuditType:    metadata.BusinessType,
		ResourceType: metadata.BusinessRes,
		ResourceName: common.DefaultAppName,
	}
)

type bizData struct {
	ID            int64       `bson:"id"`
	BizMaintainer string      `bson:"bk_biz_maintainer"`
	LifeCycle     string      `bson:"life_cycle"`
	Time          *tools.Time `bson:",inline"`
	Default       int         `bson:"default"`
	BizTester     string      `bson:"bk_biz_tester"`
	Operator      string      `bson:"operator"`
	BizProductor  string      `bson:"bk_biz_productor"`
	TimeZone      string      `bson:"time_zone"`
	Language      string      `bson:"language"`
	BizID         int64       `bson:"bk_biz_id"`
	BizName       string      `bson:"bk_biz_name"`
	BizDeveloper  string      `bson:"bk_biz_developer"`
}
