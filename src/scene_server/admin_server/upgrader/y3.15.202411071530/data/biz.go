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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	svrutils "configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

func addBizData(kit *rest.Kit, db local.DB) error {
	if err := addBizModule(kit, db, []interface{}{blueKingBizData, resBizData}, bizAuditType); err != nil {
		blog.Errorf("add biz module or set data failed, %v", err)
		return err
	}
	return nil
}

func addBizModule(kit *rest.Kit, db local.DB, data []interface{}, auditField *svrutils.AuditResType) error {
	needField := &svrutils.InsertOptions{
		UniqueFields:   []string{common.BKAppNameField},
		IgnoreKeys:     []string{common.BKAppIDField},
		IDField:        []string{common.BKAppIDField},
		AuditTypeField: auditField,
		AuditDataField: &svrutils.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   common.BKAppIDField,
			ResNameField: "bk_biz_name",
		},
	}

	var dataMap []mapstr.MapStr
	for _, item := range data {
		itemMap, err := util.ConvStructToMap(item)
		if err != nil {
			blog.Errorf("failed to convert struct to map, err: %v", err)
			return err
		}
		dataMap = append(dataMap, itemMap)
	}

	ids, err := svrutils.InsertData(kit, db, common.BKTableNameBaseApp, dataMap, needField)
	if err != nil {
		blog.Errorf("insert biz data for table %s failed, err: %v, data: %+v", common.BKTableNameBaseApp, err, data)
		return err
	}
	if len(ids) == 0 {
		blog.Errorf("failed to get biz id, data: %+v, needField: %+v", data, *needField)
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

func addBizAsstData(kit *rest.Kit, db local.DB, bizID int64, moduleNames []string) error {
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
	// add resource business module common.BKAppIDField, common.BKSetNameField, common.BKInstParentStr
	uniqueStr := strings.Join([]string{strconv.FormatInt(bizID, 16), common.DefaultResSetName,
		strconv.FormatInt(bizID, 16)}, "*")
	value := ids[uniqueStr]
	id, err := util.GetInt64ByInterface(value)
	if err != nil {
		blog.Errorf("get set id int64 failed, %v", err)
		return err
	}
	if err := addModuleData(kit, db, bizID, moduleNames, id); err != nil {
		blog.Errorf("add module data failed, %v", err)
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
	bizAuditType = &svrutils.AuditResType{
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
