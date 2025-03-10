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

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	typeHandlerMap = map[tenant.TenantTemplateType]func(kit *rest.Kit, db local.DB,
		data []tenant.TenantTmpData[mapstr.MapStr], coreAPI coreservice.CoreServiceClientInterface) error{
		tenant.TemplateTypeAssociation:       insertAsstData,
		tenant.TemplateTypeObject:            insertObjData,
		tenant.TemplateTypeObjAttribute:      insertObjAttrData,
		tenant.TemplateTypeObjAssociation:    insertObjAssociationData,
		tenant.TemplateTypeObjClassification: insertObjClassification,
		tenant.TemplateTypePlat:              insertPlatData,
		tenant.TemplateTypePropertyGroup:     insertPropertyGrp,
		tenant.TemplateTypeBizSet:            insertBizSetData,
		tenant.TemplateTypeServiceCategory:   insertSvrCategoryData,
		tenant.TemplateTypeUniqueKeys:        insertUniqueKeyData,
	}
)

func insertAsstData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.AssociationKind, 0)
	if err := db.Table(common.BKTableNameAsstDes).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameAsstDes, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.AssociationKindID] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.AssociationKindIDField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.AssociationKindIDField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.AssociationKindIDField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameAsstDes,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKFieldID] = ids[index]
	}
	if err = db.Table(common.BKTableNameAsstDes).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameAsstDes, err, kit.Rid)
		return err
	}

	// generate audit log.
	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.AssociationKindNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelAssociationRes}
	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func insertBizSetData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.BizSetInst, 0)
	if err := db.Table(common.BKTableNameBaseBizSet).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameBaseBizSet, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.BizSetName] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.BKBizSetNameField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKBizSetNameField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKBizSetNameField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameBaseBizSet,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKBizSetIDField] = ids[index]
		insertData[index][common.CreateTimeField] = time.Now()
		insertData[index][common.LastTimeField] = time.Now()
	}
	if err = db.Table(common.BKTableNameBaseBizSet).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameBaseBizSet, err, kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKBizSetIDField,
		ResourceName: common.BKBizSetNameField,
		AuditType:    metadata.BizSetType,
		ResourceType: metadata.BizSetRes,
	}
	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjAssociationData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.Association, 0)
	if err := db.Table(common.BKTableNameObjAsst).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameObjAsst, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.AsstObjID] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.AssociationObjAsstIDField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.AssociationObjAsstIDField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.AssociationObjAsstIDField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameObjAsst,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKFieldID] = ids[index]
	}
	if err = db.Table(common.BKTableNameObjAsst).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameObjAsst, err, kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.AssociationObjAsstIDField,
		AuditType:    metadata.AssociationKindType,
		ResourceType: metadata.MainlineInstanceRes}

	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjAttrData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.Attribute, 0)
	if err := db.Table(common.BKTableNameObjAttDes).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameObjAttDes, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.ObjectID] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.BKObjIDField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKObjIDField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKObjIDField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameObjAttDes,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKFieldID] = ids[index]
		insertData[index][common.CreateTimeField] = time.Now()
		insertData[index][common.LastTimeField] = time.Now()
	}
	if err = db.Table(common.BKTableNameObjAttDes).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameObjAttDes, err, kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.BKPropertyNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelAttributeRes}
	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func insertObjClassification(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.Classification, 0)
	if err := db.Table(common.BKTableNameObjClassification).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameObjClassification, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.ClassificationName] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[metadata.ClassFieldClassificationName]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", metadata.ClassificationFieldID, kit.Rid)
			return fmt.Errorf("not find field %s in data", metadata.ClassificationFieldID)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameObjClassification,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][metadata.ClassificationFieldID] = ids[index]
	}
	if err = db.Table(common.BKTableNameObjClassification).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameObjClassification, err,
			kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.BKClassificationNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelClassificationRes}

	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertPlatData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.CloudArea, 0)
	if err := db.Table(common.BKTableNameBasePlat).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameBasePlat, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.CloudName] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.BKCloudNameField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKCloudNameField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKCloudNameField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameBasePlat,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKCloudIDField] = ids[index]
		insertData[index][common.CreateTimeField] = time.Now()
		insertData[index][common.LastTimeField] = time.Now()

	}
	if err = db.Table(common.BKTableNameBasePlat).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameBasePlat, err, kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKCloudIDField,
		ResourceName: common.BKCloudNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModuleRes}

	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertPropertyGrp(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.Group, 0)
	if err := db.Table(common.BKTableNamePropertyGroup).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNamePropertyGroup, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.ObjectID+"*"+util.GetStrByInterface(item.GroupIndex)] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		objValue, ok := item.Data[common.BKObjIDField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKObjIDField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKObjIDField)
		}
		propertyIdxValue, ok := item.Data[common.BKPropertyGroupIndexField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKPropertyGroupIndexField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKPropertyGroupIndexField)
		}

		if _, ok := exitData[util.GetStrByInterface(objValue)+"*"+util.GetStrByInterface(propertyIdxValue)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNamePropertyGroup,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKFieldID] = ids[index]
	}
	if err = db.Table(common.BKTableNamePropertyGroup).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNamePropertyGroup, err,
			kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.BKPropertyGroupNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelGroupRes}

	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertObjData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.Object, 0)
	if err := db.Table(common.BKTableNameObjDes).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameObjDes, err)
		return err
	}

	exitData := make(map[string]interface{}, 0)
	for _, item := range result {
		exitData[item.ObjectID] = struct{}{}
	}
	insertData := make([]mapstr.MapStr, 0)
	for _, item := range data {
		value, ok := item.Data[common.BKObjIDField]
		if !ok {
			blog.Errorf("not find field %s in data, rid: %s", common.BKObjIDField, kit.Rid)
			return fmt.Errorf("not find field %s in data", common.BKObjIDField)
		}

		if _, ok := exitData[util.GetStrByInterface(value)]; ok {
			continue
		}
		insertData = append(insertData, item.Data)
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameObjDes,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index][common.BKFieldID] = ids[index]
		insertData[index][common.CreateTimeField] = time.Now()
		insertData[index][common.LastTimeField] = time.Now()
	}
	if err = db.Table(common.BKTableNameObjDes).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameObjDes, err, kit.Rid)
		return err
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.BKObjNameField,
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModuleRes}

	if err := addAuditLog(kit, db, coreAPI, insertData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func insertUniqueKeyData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.ObjectUnique, 0)
	if err := db.Table(common.BKTableNameObjUnique).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameObjUnique, err)
		return err
	}

	uniqueData := make([]tenant.UniqueKeyTmp, len(data))
	for index, item := range data {
		bsonData, err := bson.Marshal(item.Data)
		if err != nil {
			blog.Errorf("bson marshal failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		bson.Unmarshal(bsonData, &uniqueData[index])
	}

	exitData := make(map[string]interface{}, 0)
	for index := range result {
		sort.Slice(result[index].Keys, func(i, j int) bool {
			return result[index].Keys[i].Kind < result[index].Keys[j].Kind
		})
		uniqueArr := make([]string, 0)
		for _, key := range result[index].Keys {
			uniqueArr = append(uniqueArr, key.Kind)
		}
		exitData[strings.Join(uniqueArr, "*")] = struct{}{}
	}
	// get attribute data
	attrArr := make([]metadata.Attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(nil).All(kit.Ctx, &attrArr)
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
		sort.Strings(uniqueData[index].Keys)
		uniqueStr := strings.Join(uniqueData[index].Keys, "*")
		if _, ok := exitData[uniqueStr]; ok {
			continue
		}

		keys := make([]metadata.UniqueKey, 0)
		for _, field := range uniqueData[index].Keys {
			keys = append(keys, metadata.UniqueKey{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   attrIDMap[generateUniqueKey(uniqueData[index].ObjectID, field)],
			})
		}
		insertData = append(insertData, metadata.ObjectUnique{
			Keys:     keys,
			ObjID:    uniqueData[index].ObjectID,
			Ispre:    true,
			LastTime: metadata.Time{Time: time.Now()},
		})
	}

	if len(insertData) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameObjUnique,
		len(insertData))
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	for index := range insertData {
		insertData[index].ID = ids[index]
	}
	if err = db.Table(common.BKTableNameObjUnique).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameObjUnique, err, kit.Rid)
		return err
	}

	auditData := make([]mapstr.MapStr, len(insertData))
	for index, item := range insertData {
		itemMap, err := structToMap(item)
		if err != nil {
			blog.Errorf("struct to map failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditData[index] = itemMap
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   "id",
		ResourceName: "bk_obj_id",
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelUniqueRes}

	if err := addAuditLog(kit, db, coreAPI, auditData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func insertSvrCategoryData(kit *rest.Kit, db local.DB, data []tenant.TenantTmpData[mapstr.MapStr],
	coreAPI coreservice.CoreServiceClientInterface) error {

	result := make([]metadata.ServiceCategory, 0)
	if err := db.Table(common.BKTableNameServiceCategory).Find(mapstr.MapStr{}).All(kit.Ctx, &result); err != nil {
		blog.Errorf("get data from table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	svrCategoryTmp := make([]tenant.SvrCategoryTmp, len(data))
	for index, item := range data {
		bsonData, err := bson.Marshal(item.Data)
		if err != nil {
			blog.Errorf("bson marshal failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		bson.Unmarshal(bsonData, &svrCategoryTmp[index])
	}

	exitData := make(map[string]int64, 0)
	for _, item := range result {
		subFlag := "0"
		if item.ParentID != 0 {
			subFlag = "1"
		}

		exitData[item.Name+"*"+subFlag] = item.ID
	}

	insertData := make([]*metadata.ServiceCategory, 0)
	exitParent := make(map[string]int64, 0)
	insertParent := make(map[string]*metadata.ServiceCategory, 0)
	insertSubCategory := make(map[string][]*metadata.ServiceCategory, 0)
	// get insert parent category
	insertCount := 0
	for _, item := range svrCategoryTmp {
		subFlag := "0"
		if item.ParentName != "" {
			subFlag = "1"
		}
		uniqueStr := item.Name + "*" + subFlag
		if id, ok := exitData[uniqueStr]; ok {
			if subFlag == "1" {
				exitParent[item.ParentName] = id
			}
			continue
		}
		if subFlag != "1" {
			insertParent[item.Name] = &metadata.ServiceCategory{
				Name:      item.Name,
				IsBuiltIn: true,
			}
			insertCount++
			continue
		}
		insertSubCategory[item.ParentName] = append(insertSubCategory[item.ParentName], &metadata.ServiceCategory{
			Name:      item.Name,
			IsBuiltIn: true,
		})
		insertCount++
	}

	if len(insertSubCategory) == 0 {
		return nil
	}
	ids, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameServiceCategory,
		insertCount)
	if err != nil {
		blog.Errorf("get next sequence failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	idxCount := 0
	for key := range insertParent {
		insertParent[key].ID = int64(ids[idxCount])
		insertParent[key].RootID = int64(ids[idxCount])
		exitParent[key] = int64(ids[idxCount])
		insertData = append(insertData, insertParent[key])
		idxCount++
	}

	for parentName, subValues := range insertSubCategory {
		for index := range subValues {
			subValues[index].ID = int64(ids[idxCount])
			subValues[index].ParentID = exitParent[parentName]
			subValues[index].RootID = exitParent[parentName]
			insertData = append(insertData, subValues[index])
			idxCount++
		}
	}

	if err = db.Table(common.BKTableNameServiceCategory).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("insert data for table %s failed, err: %v, rid: %s", common.BKTableNameServiceCategory, err,
			kit.Rid)
		return err
	}

	auditData := make([]mapstr.MapStr, len(insertData))
	for index, item := range insertData {
		itemMap, err := structToMap(item)
		if err != nil {
			blog.Errorf("struct to map failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditData[index] = itemMap
	}

	auditLog := &auditlog.AuditOpts{
		ResourceID:   common.BKFieldID,
		ResourceName: common.BKFieldName,
		AuditType:    metadata.PlatformSetting,
		ResourceType: metadata.ServiceCategoryRes}

	if err := addAuditLog(kit, db, coreAPI, auditData, auditLog); err != nil {
		blog.Errorf("add audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func generateUniqueKey(objID, propertyID string) string {
	return objID + ":" + propertyID
}

func addAuditLog(kit *rest.Kit, db local.DB, coreAPI coreservice.CoreServiceClientInterface, insertData []mapstr.MapStr,
	auditOpt *auditlog.AuditOpts) error {

	audit := auditlog.NewTenantTemplateAudit(coreAPI)
	generateAuditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog := audit.GenerateAuditLog(generateAuditParam, insertData, auditOpt)

	// save audit log.
	err := auditlog.SaveAuditLog(kit, db, auditLog...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = bson.Unmarshal(data, &result)
	return result, err
}
