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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

var (
	typeHandlerMap = map[string]func(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error{
		metadata.TemplateTypeAssociation:       insertAsstData,
		metadata.TemplateTypeObject:            insertObjData,
		metadata.TemplateTypeObjAttribute:      insertObjAttrData,
		metadata.TemplateTypeObjAssociation:    insertObjAssociationData,
		metadata.TempalteTypeObjClassification: insertObjClassification,
		metadata.TempalteTypePlat:              insertPlatData,
		metadata.TempalteTypePropertyGroup:     insertPropertyGrp,
		metadata.TemplateTypeBizSet:            insertBizSetData,
	}
)

func insertAsstData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.AssociationKindIDField},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{metadata.AttributeFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAssociationRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_asst_name",
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameAsstDes, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameAsstDes, err)
	}
	return nil
}

func insertBizSetData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	for _, item := range data {
		item["create_time"] = time.Now()
		item["last_time"] = time.Now()
	}
	needField := &utils.InsertOptions{
		UniqueFields: []string{common.BKBizSetNameField},
		IgnoreKeys:   []string{common.BKBizSetIDField},
		IDField:      []string{common.BKBizSetIDField},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.BizSetType,
			ResourceType: metadata.BizSetRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "bk_biz_set_id",
			ResNameField: "bk_biz_set_name",
		},
	}

	_, err := utils.InsertData(kit, db, common.BKTableNameBaseBizSet, data, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameBaseBizSet, err)
	}
	return nil
}

func insertObjAssociationData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	insertOps := &utils.InsertOptions{
		UniqueFields: []string{"bk_obj_asst_id"},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{common.BKFieldID},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_asst_id",
		},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.AssociationKindType,
			ResourceType: metadata.MainlineInstanceRes,
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameObjAsst, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameAsstDes, err)
	}
	return nil
}

func insertObjAttrData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	for _, item := range data {
		item["create_time"] = time.Now()
		item["last_time"] = time.Now()
	}
	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField},
		IgnoreKeys:   []string{"id", "bk_property_index"},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAttributeRes,
		},
		AuditDataField: &utils.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   common.BKFieldID,
			ResNameField: "bk_property_name",
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameObjAttDes, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameAsstDes, err)
	}
	return nil
}

func insertObjClassification(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	insertOps := &utils.InsertOptions{
		UniqueFields: []string{"bk_classification_name"},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{metadata.ClassificationFieldID},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_classification_name",
		},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelClassificationRes,
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameObjClassification, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjClassification, err)
	}
	return nil
}

func insertPlatData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	for _, item := range data {
		item["create_time"] = time.Now()
		item["last_time"] = time.Now()
	}
	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.BKCloudNameField},
		IgnoreKeys:   []string{common.BKCloudIDField},
		IDField:      []string{common.BKCloudIDField},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   common.BKCloudIDField,
			ResNameField: "bk_cloud_name",
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameBasePlat, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameBasePlat, err)
	}
	return nil
}

func insertPropertyGrp(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.BKObjIDField, common.BKAppIDField, common.BKPropertyGroupIndexField},
		IgnoreKeys:   []string{common.BKFieldID, common.BKPropertyGroupIndexField},
		IDField:      []string{common.BKFieldID},
		AuditDataField: &utils.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   common.BKFieldID,
			ResNameField: "bk_group_name",
		},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelGroupRes,
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNamePropertyGroup, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNamePropertyGroup, err)
	}
	return nil
}

func insertObjData(kit *rest.Kit, db local.DB, data []mapstr.MapStr) error {

	for _, item := range data {
		item["create_time"] = time.Now()
		item["last_time"] = time.Now()
	}
	insertOps := &utils.InsertOptions{
		UniqueFields: []string{"bk_obj_id"},
		IgnoreKeys:   []string{"id", "obj_sort_number"},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_name",
		},
	}
	_, err := utils.InsertData(kit, db, common.BKTableNameObjDes, data, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjDes, err)
	}
	return nil
}

func insertUniqueKey(kit *rest.Kit, db local.DB) error {

	filter := mapstr.MapStr{
		"type": "unique_keys",
	}
	data := make([]metadata.UniqueKeyTmp, 0)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(filter).All(kit.Ctx, &data)
	if err != nil {
		blog.Errorf("get template data for types unique keys failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	attrArr := make([]metadata.Attribute, 0)
	err = db.Table(common.BKTableNameObjAttDes).Find(nil).All(kit.Ctx, &attrArr)
	if err != nil {
		blog.Errorf("get host unique fields failed, err: %v", err)
		return err
	}

	attrIDMap := make(map[string]uint64)
	for _, attr := range attrArr {
		attrIDMap[generateUniqueKey(attr.ObjectID, attr.PropertyID)] = uint64(attr.ID)
	}

	var insertData []mapstr.MapStr
	for _, item := range data {
		keys := make([]uniqueKey, 0)
		for _, property := range item.Attributes {
			keys = append(keys, uniqueKey{
				Kind: "property",
				ID:   attrIDMap[generateUniqueKey(item.ObjectID, property)],
			})
		}
		item.Data[common.BKObjectUniqueKeys] = keys
		insertData = append(insertData, item.Data)
	}
	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.BKObjectUniqueKeys},
		IgnoreKeys:   []string{common.BKFieldID},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelUniqueRes,
		},
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_id",
		},
	}

	_, err = utils.InsertData(kit, db, common.BKTableNameObjUnique, insertData, insertOps)
	if err != nil {
		blog.Errorf("insert unique keys data for table %s failed, err: %v, rid: %s", common.BKTableNameObjUnique, err,
			kit.Rid)
		return err
	}
	return nil
}

func generateUniqueKey(objID, propertyID string) string {
	return objID + ":" + propertyID
}

type uniqueKey struct {
	Kind string `bson:"key_kind"`
	ID   uint64 `bson:"key_id"`
}
