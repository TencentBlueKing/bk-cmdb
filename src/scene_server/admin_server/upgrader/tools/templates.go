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

package tools

import (
	"fmt"
	"strings"

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

func InsertTemplateData(kit *rest.Kit, db dal.RDB, data []mapstr.MapStr,
	insertOps *InsertOptions, idOption *IDOptions, dataType tenanttmp.TenantTemplateType) error {

	tmpData := make([]tenanttmp.TenantTmpData[mapstr.MapStr], 0)
	for _, item := range data {
		for _, idField := range idOption.RemoveKeys {
			delete(item, idField)
		}
		tmpData = append(tmpData, tenanttmp.TenantTmpData[mapstr.MapStr]{
			Type:  dataType,
			IsPre: true,
			Data:  item,
		})
	}

	result := make([]tenanttmp.TenantTmpData[mapstr.MapStr], 0)
	cond := mapstr.MapStr{
		"type": string(dataType),
	}
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(cond).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exist data failed, table: %s, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	insertData, err := cmpTenantTmp(result, tmpData, insertOps.UniqueFields, insertOps.IgnoreKeys)
	if err != nil {
		blog.Errorf("compare data failed, err: %v", err)
		return err
	}

	// get audit resName
	resNames := make([]string, 0)
	for _, item := range insertData {
		resNames = append(resNames, util.GetStrByInterface(item.Data[idOption.ResNameField]))
	}

	if err = insertTmpData[mapstr.MapStr](kit, db, common.BKTableNameTenantTemplate, insertData,
		resNames); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	return nil
}

func InsertSvrCategoryTmp(kit *rest.Kit, db dal.RDB, data []tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp]) error {

	result := make([]tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp], 0)
	cond := mapstr.MapStr{
		"type": tenanttmp.TemplateTypeServiceCategory,
	}
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(cond).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exist data failed, table: %s, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	insertData := make([]tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp], 0)
	existUniqueMap := make(map[string]interface{}, 0)
	for _, item := range result {
		existUniqueMap[item.Data.ParentName+"*"+item.Data.Name] = struct{}{}
	}

	resNames := make([]string, 0)
	for _, item := range data {
		uniqueStr := item.Data.ParentName + "*" + item.Data.Name
		if _, ok := existUniqueMap[uniqueStr]; ok {
			continue
		}
		insertData = append(insertData, item)
		resNames = append(resNames, item.Data.Name)
	}

	if err = insertTmpData[tenanttmp.SvrCategoryTmp](kit, db, common.BKTableNameTenantTemplate,
		insertData, resNames); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	return nil
}

func InsertUniqueKeyTmp(kit *rest.Kit, db dal.RDB, data []tenanttmp.TenantTmpData[tenanttmp.UniqueKeyTmp]) error {

	result := make([]tenanttmp.TenantTmpData[tenanttmp.UniqueKeyTmp], 0)
	cond := mapstr.MapStr{
		"type": tenanttmp.TemplateTypeUniqueKeys,
	}
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(cond).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exist data failed, table: %s, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	insertData := make([]tenanttmp.TenantTmpData[tenanttmp.UniqueKeyTmp], 0)
	existUniqueMap := make(map[string]interface{}, 0)
	for _, item := range result {
		existUniqueMap[strings.Join(item.Data.Keys, "*")] = struct{}{}
	}

	resNames := make([]string, 0)
	for _, item := range data {
		uniqueStr := strings.Join(item.Data.Keys, "*")
		if _, ok := existUniqueMap[uniqueStr]; ok {
			continue
		}
		insertData = append(insertData, item)
		resNames = append(resNames, util.GetStrByInterface(item.Data.ObjectID))
	}

	if err = insertTmpData[tenanttmp.UniqueKeyTmp](kit, db, common.BKTableNameTenantTemplate, insertData,
		resNames); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameTenantTemplate, err)
		return err
	}
	return nil
}

func insertTmpData[T tenanttmp.UniqueKeyTmp | tenanttmp.SvrCategoryTmp | mapstr.MapStr](kit *rest.Kit, db dal.RDB,
	table string, insertData []tenanttmp.TenantTmpData[T], resNames []string) error {

	if len(insertData) == 0 {
		blog.Infof("no data to insert, table: %s", table)
		return nil
	}

	nextIDs, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next %d data IDs failed, table: %s, err: %v", len(insertData), table, err)
		return fmt.Errorf("get next %d data IDs failed, table: %s, err: %v", len(insertData), table, err)
	}
	for index := range insertData {
		insertData[index].ID = int64(nextIDs[index])
	}

	if err := mongodb.Dal().Shard(kit.SysShardOpts()).Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("add data for table %s failed, data: %+v, err: %v", insertData, err)
		return err
	}

	// add audit log
	auditField := &AuditStruct{
		AuditDataField: &AuditDataField{
			ResIDField: "id",
		},
		AuditTypeData: &AuditResType{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.TenantTemplateRes,
		},
		ResNames: resNames,
	}
	auditDataMap := make([]map[string]interface{}, 0)
	for _, item := range insertData {
		dataMap, err := ConvStructToMap(item)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return err
		}
		auditDataMap = append(auditDataMap, dataMap)
	}

	if err := AddTmpAuditLog(kit, db, auditDataMap, auditField); err != nil {
		blog.Errorf("add audit log failed, err: %v", err)
		return err
	}
	return nil
}

func cmpTenantTmp(existData, data []tenanttmp.TenantTmpData[mapstr.MapStr],
	uniqueFields, ignoreFields []string) ([]tenanttmp.TenantTmpData[mapstr.MapStr], error) {

	existMap := make(map[string]tenanttmp.TenantTmpData[mapstr.MapStr])
	for _, item := range existData {
		valueStr := getUniqueStr(item.Data, uniqueFields)
		if valueStr == "" {
			continue
		}
		existMap[valueStr] = item
	}

	insertData := make([]tenanttmp.TenantTmpData[mapstr.MapStr], 0)
	for _, item := range data {
		valueStr := getUniqueStr(item.Data, uniqueFields)
		if valueStr == "" {
			continue
		}
		if _, exist := existMap[valueStr]; !exist {
			insertData = append(insertData, item)
			continue
		}
		if err := cmpData(item.Data, existMap[valueStr].Data, ignoreFields); err != nil {
			return nil, err
		}
	}

	return insertData, nil
}
