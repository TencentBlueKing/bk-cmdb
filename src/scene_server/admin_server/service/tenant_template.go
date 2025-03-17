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

package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

var (
	typeHandlerMap = map[tenanttmp.TenantTemplateType]func(kit *rest.Kit, db local.DB) error{
		tenanttmp.TemplateTypeAssociation:       insertAsstData,
		tenanttmp.TemplateTypeObject:            insertObjData,
		tenanttmp.TemplateTypeObjAttribute:      insertObjAttrData,
		tenanttmp.TemplateTypeObjAssociation:    insertObjAssociationData,
		tenanttmp.TemplateTypeObjClassification: insertObjClassification,
		tenanttmp.TemplateTypePlat:              insertPlatData,
		tenanttmp.TemplateTypePropertyGroup:     insertPropertyGrp,
		tenanttmp.TemplateTypeBizSet:            insertBizSetData,
		tenanttmp.TemplateTypeServiceCategory:   insertSvrCategoryData,
		tenanttmp.TemplateTypeUniqueKeys:        insertUniqueKeyData,
	}
)

func insertAsstData(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameAsstDes
	data, err := getTemplateData[metadata.AssociationKind](kit, tenanttmp.TemplateTypeAssociation)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.AssociationKind, 0)
	if err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.AssociationKindIDField).All(kit.Ctx,
		&result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.AssociationKindID] = struct{}{}
	}

	insertData := make([]metadata.AssociationKind, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.AssociationKindID]; ok {
			continue
		}
		insertData = append(insertData, item.Data)

	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].AssociationKindName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeAssociation,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	// generate audit log.
	if err := addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func insertBizSetData(kit *rest.Kit, db local.DB) error {
	table := common.BKTableNameBaseBizSet
	data, err := getTemplateData[metadata.BizSetInst](kit, tenanttmp.TemplateTypeBizSet)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.BizSetInst, 0)
	if err := db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKBizSetNameField).All(kit.Ctx,
		&result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.BizSetName] = struct{}{}
	}
	insertData := make([]metadata.BizSetInst, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.BizSetName]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].BizSetID = int64(ids[index])
		insertData[index].CreateTime = metadata.Time{Time: time.Now()}
		insertData[index].LastTime = metadata.Time{Time: time.Now()}
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].BizSetName,
			ResourceID:   insertData[index].BizSetID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeBizSet,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjAssociationData(kit *rest.Kit, db local.DB) error {
	table := common.BKTableNameObjAsst
	data, err := getTemplateData[metadata.Association](kit, tenanttmp.TemplateTypeObjAssociation)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.Association, 0)
	if err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.AssociationObjAsstIDField).All(kit.Ctx,
		&result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.AsstObjID] = struct{}{}
	}
	insertData := make([]metadata.Association, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.AsstObjID]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].AssociationName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeObjAssociation,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjAttrData(kit *rest.Kit, db local.DB) error {
	table := common.BKTableNameObjAttDes
	data, err := getTemplateData[metadata.Attribute](kit, tenanttmp.TemplateTypeObjAttribute)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}

	result := make([]metadata.Attribute, 0)
	err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKObjIDField, common.BKPropertyIDField).All(kit.Ctx,
		&result)
	if err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.ObjectID+"*"+item.PropertyID] = struct{}{}
	}
	insertData := make([]metadata.Attribute, 0)
	for _, item := range data {
		value := item.Data.ObjectID + "*" + item.Data.PropertyID
		if _, ok := existData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		insertData[index].CreateTime = &metadata.Time{Time: time.Now()}
		insertData[index].LastTime = &metadata.Time{Time: time.Now()}
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].PropertyName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeObjAttribute,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func insertObjClassification(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameObjClassification
	data, err := getTemplateData[metadata.Classification](kit, tenanttmp.TemplateTypeObjClassification)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.Classification, 0)
	err = db.Table(table).Find(mapstr.MapStr{}).Fields(metadata.ClassFieldClassificationName).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.ClassificationName] = struct{}{}
	}
	insertData := make([]metadata.Classification, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.ClassificationName]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].ClassificationName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeObjClassification,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err,
			kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertPlatData(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameBasePlat
	data, err := getTemplateData[metadata.CloudArea](kit, tenanttmp.TemplateTypePlat)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.CloudArea, 0)
	if err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKCloudNameField).All(kit.Ctx,
		&result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.CloudName] = struct{}{}
	}
	insertData := make([]metadata.CloudArea, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.CloudName]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].CloudID = int64(ids[index])
		insertData[index].CreateTime = time.Now()
		insertData[index].LastTime = time.Now()
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].CloudName,
			ResourceID:   insertData[index].CloudID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypePlat,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertPropertyGrp(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNamePropertyGroup
	data, err := getTemplateData[metadata.Group](kit, tenanttmp.TemplateTypePropertyGroup)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.Group, 0)
	err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKObjIDField, common.BKPropertyGroupIDField).All(kit.Ctx,
		&result)
	if err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.ObjectID+"*"+util.GetStrByInterface(item.GroupID)] = struct{}{}
	}
	insertData := make([]metadata.Group, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.ObjectID+"*"+item.Data.GroupID]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].GroupName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypePropertyGroup,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjData(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameObjDes
	data, err := getTemplateData[metadata.Object](kit, tenanttmp.TemplateTypeObject)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}

	result := make([]metadata.Object, 0)
	if err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKObjIDField).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for _, item := range result {
		existData[item.ObjectID] = struct{}{}
	}
	insertData := make([]metadata.Object, 0)
	for _, item := range data {
		if _, ok := existData[item.Data.ObjectID]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = int64(ids[index])
		insertData[index].CreateTime = &metadata.Time{Time: time.Now()}
		insertData[index].LastTime = &metadata.Time{Time: time.Now()}
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertData[index].ObjectName,
			ResourceID:   insertData[index].ID,
			Data:         insertData[index],
			Type:         tenanttmp.TemplateTypeObject,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func getUniqueKeysStr(keys []metadata.UniqueKey, objID string) string {
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].ID < keys[j].ID
	})
	uniqueArr := []string{objID}
	for _, key := range keys {
		uniqueArr = append(uniqueArr, fmt.Sprintf("%d", key.ID))
	}

	return strings.Join(uniqueArr, "*")
}

func insertUniqueKeyData(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameObjUnique
	uniqueData, err := getTemplateData[tenanttmp.UniqueKeyTmp](kit, tenanttmp.TemplateTypeUniqueKeys)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.ObjectUnique, 0)
	err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKObjIDField, common.BKObjectUniqueKeys).All(kit.Ctx,
		&result)
	if err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]interface{}, 0)
	for index := range result {
		existData[getUniqueKeysStr(result[index].Keys, result[index].ObjID)] = struct{}{}
	}
	// get attribute data
	attrArr := make([]metadata.Attribute, 0)
	err = db.Table(table).Find(nil).All(kit.Ctx, &attrArr)
	if err != nil {
		blog.Errorf("get host unique fields failed, err: %v", err)
		return err
	}
	attrIDMap := make(map[string]uint64)
	for _, attr := range attrArr {
		attrIDMap[generateUniqueKey(attr.ObjectID, attr.PropertyID)] = uint64(attr.ID)
	}

	insertData := make([]metadata.ObjectUnique, 0)
	for index := range uniqueData {
		keys := make([]metadata.UniqueKey, 0)
		for _, field := range uniqueData[index].Data.Keys {
			keys = append(keys, metadata.UniqueKey{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   attrIDMap[generateUniqueKey(uniqueData[index].Data.ObjectID, field)],
			})
		}

		uniqueStr := getUniqueKeysStr(keys, uniqueData[index].Data.ObjectID)
		if _, ok := existData[uniqueStr]; ok {
			continue
		}

		insertData = append(insertData, metadata.ObjectUnique{
			Keys:     keys,
			ObjID:    uniqueData[index].Data.ObjectID,
			Ispre:    true,
			LastTime: metadata.Time{Time: time.Now()},
		})
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for index := range insertData {
		insertData[index].ID = ids[index]
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			Data:         insertData[index],
			ResourceID:   int64(insertData[index].ID),
			ResourceName: "",
			Type:         tenanttmp.TemplateTypeUniqueKeys,
		})
	}
	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func getSvrCategoryUniqueStr(name string, isSubCategory bool) string {
	if isSubCategory {
		return fmt.Sprintf("%s*1", name)
	}
	return fmt.Sprintf("%s*0", name)
}

func insertSvrCategoryData(kit *rest.Kit, db local.DB) error {

	table := common.BKTableNameServiceCategory
	svrCategoryTmp, err := getTemplateData[tenanttmp.SvrCategoryTmp](kit, tenanttmp.TemplateTypeServiceCategory)
	if err != nil {
		blog.Errorf("get template data failed, err: %v", err)
		return err
	}
	result := make([]metadata.ServiceCategory, 0)
	err = db.Table(table).Find(mapstr.MapStr{}).Fields(common.BKFieldName, common.BKParentIDField).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get data from table %s failed, err: %v", table, err)
		return err
	}

	existData := make(map[string]int64, 0)
	for _, item := range result {
		existData[getSvrCategoryUniqueStr(item.Name, item.ParentID != 0)] = item.ID
	}

	insertData := make([]*metadata.ServiceCategory, 0)
	existParent := make(map[string]int64, 0)
	insertParent := make(map[string]*metadata.ServiceCategory, 0)
	insertSubCategory := make(map[string][]*metadata.ServiceCategory, 0)
	// get insert parent category
	insertCount := 0
	for _, item := range svrCategoryTmp {
		uniqueStr := getSvrCategoryUniqueStr(item.Data.Name, item.Data.ParentName != "")
		if id, ok := existData[uniqueStr]; ok {
			if item.Data.ParentName != "" {
				existParent[item.Data.ParentName] = id
			}
			continue
		}
		if item.Data.ParentName == "" {
			insertParent[item.Data.Name] = &metadata.ServiceCategory{
				Name:      item.Data.Name,
				IsBuiltIn: true,
			}
			insertCount++
			continue
		}
		insertSubCategory[item.Data.ParentName] = append(insertSubCategory[item.Data.ParentName],
			&metadata.ServiceCategory{
				Name:      item.Data.Name,
				IsBuiltIn: true,
			})
		insertCount++
	}

	if insertCount == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, insertCount)
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	idxCount := 0
	auditLogs := make([]*auditlog.TenantTmpAuditOpts, 0)
	for key := range insertParent {
		insertParent[key].ID = int64(ids[idxCount])
		insertParent[key].RootID = int64(ids[idxCount])
		existParent[key] = int64(ids[idxCount])
		insertData = append(insertData, insertParent[key])
		idxCount++
		auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
			ResourceName: insertParent[key].Name,
			ResourceID:   insertParent[key].ID,
			Data:         insertData,
			Type:         tenanttmp.TemplateTypeServiceCategory,
		})
	}

	for parentName, subValues := range insertSubCategory {
		for index := range subValues {
			subValues[index].ID = int64(ids[idxCount])
			subValues[index].ParentID = existParent[parentName]
			subValues[index].RootID = existParent[parentName]
			insertData = append(insertData, subValues[index])
			idxCount++
			auditLogs = append(auditLogs, &auditlog.TenantTmpAuditOpts{
				ResourceName: subValues[index].Name,
				ResourceID:   subValues[index].ID,
				Data:         insertData,
				Type:         tenanttmp.TemplateTypeServiceCategory,
			})
		}
	}

	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err = addAuditLog(kit, db, auditLogs); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func generateUniqueKey(objID, propertyID string) string {
	return objID + ":" + propertyID
}

func addAuditLog(kit *rest.Kit, db local.DB, insertData []*auditlog.TenantTmpAuditOpts) error {

	audit := auditlog.NewTenantTemplateAuditLog()
	auditLog := audit.GenerateAuditLog(insertData)

	// save audit log.
	err := audit.SaveAuditLog(kit, db, auditLog...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func getTemplateData[T any](kit *rest.Kit, ty tenanttmp.TenantTemplateType) ([]tenanttmp.TenantTmpData[T], error) {

	tmpData := make([]tenanttmp.TenantTmpData[T], 0)
	lastId := 0
	for {
		filter := mapstr.MapStr{
			"type": ty,
			"id":   map[string]interface{}{common.BKDBGT: lastId},
		}
		result := make([]tenanttmp.TenantTmpData[T], 0)
		err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(filter).
			Sort("id").Limit(uint64(common.BKMaxInstanceLimit)).All(kit.Ctx, &result)
		if err != nil {
			blog.Errorf("get template data for type %s failed, err: %v, rid: %s", ty, err, kit.Rid)
			return nil, err
		}

		if len(result) > 0 {
			tmpData = append(tmpData, result...)
			lastId = int(result[len(result)-1].ID)
		}
		if len(result) < common.BKMaxInstanceLimit {
			break
		}
	}
	return tmpData, nil
}
