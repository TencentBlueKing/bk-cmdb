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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/scene_server/admin_server/upgrader/types"
	"configcenter/src/storage/dal"
)

func addBizData(kit *rest.Kit, db dal.Dal) error {

	if kit.TenantID == types.GetBlueKing() {
		if err := addBizModule(kit, db, []interface{}{blueKingBizData, resBizData},
			[]tools.AuditType{blueKingBizAudit, resBizAudit}); err != nil {
			blog.Errorf("add biz module or set data failed, %v", err)
			return err
		}
		return nil
	}
	if err := addBizModule(kit, db, []interface{}{resBizData}, []tools.AuditType{resBizAudit}); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}
	return nil

}

func addBizModule(kit *rest.Kit, db dal.Dal, data []interface{}, auditField []tools.AuditType) error {

	cmpField := &tools.CmpFiled{
		UniqueFields: []string{common.BKAppNameField},
		IgnoreKeys:   []string{common.BKAppIDField},
		IDField:      common.BKAppIDField,
	}
	auditDataField := &tools.AuditDataField{
		BusinessID:   "bk_biz_id",
		ResourceID:   common.BKAppIDField,
		ResourceName: "bk_biz_name",
	}
	ids, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameBaseApp, data, cmpField, auditField,
		auditDataField)
	if err != nil {
		blog.Errorf("insert biz data for table %s failed, err: %v, data: %+v", common.BKTableNameBaseApp, err, data)
		return err
	}
	if len(ids) == 0 {
		blog.Errorf("failed to get biz id, data: %+v, cmpField: %+v", data, *cmpField)
		return fmt.Errorf("failed to get biz id")
	}

	resBizID, err := util.GetInt64ByInterface(ids[common.DefaultAppName])
	if err != nil {
		blog.Errorf("get biz id failed, %v", err)
		return err
	}

	resModuleName := []string{common.DefaultResModuleName}
	if err := addBizAsstData(kit, db, resBizID, resModuleName); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}

	if kit.TenantID != types.GetBlueKing() {
		return nil
	}

	bkBizID, err := util.GetInt64ByInterface(ids[common.BKAppName])
	if err != nil {
		blog.Errorf("get biz id failed, %v", err)
		return err
	}
	bkModuleName := []string{common.DefaultResModuleName, common.DefaultFaultModuleName,
		common.DefaultRecycleModuleName}
	if err := addBizAsstData(kit, db, bkBizID, bkModuleName); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}
	return nil
}

func addBizAsstData(kit *rest.Kit, db dal.Dal, bizID int64, moduleNames []string) error {
	// add resource business cluster
	ids, err := addSetBaseData(kit, db, bizID)
	if err != nil {
		blog.Errorf("add set data failed, %v", err)
		return err
	}
	if len(ids) == 0 {
		blog.Errorf("failed to get set id")
		return fmt.Errorf("failed to get set id")
	}
	// add resource business module
	for _, value := range ids {
		if err := addModuleData(kit, db, bizID, moduleNames, value.(uint64)); err != nil {
			blog.Errorf("add module data failed, %v", err)
			return err
		}
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
	blueKingBizAudit = tools.AuditType{
		AuditType:    metadata.BusinessType,
		ResourceType: metadata.BusinessRes,
	}
	resBizAudit = tools.AuditType{
		AuditType:    metadata.BusinessType,
		ResourceType: metadata.BusinessRes,
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
