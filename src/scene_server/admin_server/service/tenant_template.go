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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
	"github.com/mohae/deepcopy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	tableFieldsMap = map[string]*tools.InsertOptions{
		common.BKTableNameAsstDes: {
			UniqueFields: []string{common.AssociationKindIDField},
			IgnoreKeys:   []string{"id"},
			IDField:      []string{metadata.AttributeFieldID},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelAssociationRes,
			},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   "id",
				ResNameField: "bk_asst_name",
			},
		},
		common.BKTableNameObjDes: {
			UniqueFields: []string{"bk_obj_id"},
			IgnoreKeys:   []string{"id", "obj_sort_number"},
			IDField:      []string{common.BKFieldID},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModuleRes,
			},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   "id",
				ResNameField: "bk_obj_name",
			},
		},
		common.BKTableNameObjAsst: {
			UniqueFields: []string{"bk_obj_asst_id"},
			IgnoreKeys:   []string{"id"},
			IDField:      []string{common.BKFieldID},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   "id",
				ResNameField: "bk_obj_asst_id",
			},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.AssociationKindType,
				ResourceType: metadata.MainlineInstanceRes,
			},
		},
		common.BKTableNameObjAttDes: {
			UniqueFields: []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField},
			IgnoreKeys:   []string{"id", "bk_property_index"},
			IDField:      []string{common.BKFieldID},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelAttributeRes,
			},
			AuditDataField: &tools.AuditDataField{
				BizIDField:   "bk_biz_id",
				ResIDField:   common.BKFieldID,
				ResNameField: "bk_property_name",
			},
		},
		common.BKTableNameObjClassification: {
			UniqueFields: []string{"bk_classification_name"},
			IgnoreKeys:   []string{"id"},
			IDField:      []string{metadata.ClassificationFieldID},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   "id",
				ResNameField: "bk_classification_name",
			},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelClassificationRes,
			},
		},
		common.BKTableNameBasePlat: {
			UniqueFields: []string{common.BKCloudNameField},
			IgnoreKeys:   []string{common.BKCloudIDField},
			IDField:      []string{common.BKCloudIDField},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModuleRes,
			},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   common.BKCloudIDField,
				ResNameField: "bk_cloud_name",
			},
		},
		common.BKTableNamePropertyGroup: {
			UniqueFields: []string{common.BKObjIDField, common.BKAppIDField, common.BKPropertyGroupIndexField},
			IgnoreKeys:   []string{common.BKFieldID, common.BKPropertyGroupIndexField},
			IDField:      []string{common.BKFieldID},
			AuditDataField: &tools.AuditDataField{
				BizIDField:   "bk_biz_id",
				ResIDField:   common.BKFieldID,
				ResNameField: "bk_group_name",
			},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelGroupRes,
			},
		},
		common.BKTableNameServiceCategory: {
			UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
			IgnoreKeys:   []string{common.BKFieldID},
			IDField:      []string{common.BKFieldID},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.PlatformSetting,
				ResourceType: metadata.ServiceCategoryRes,
			},
			AuditDataField: &tools.AuditDataField{
				BizIDField:   "bk_biz_id",
				ResIDField:   "id",
				ResNameField: "name",
			},
		},
		common.BKTableNameObjUnique: {
			UniqueFields: []string{"keys"},
			IgnoreKeys:   []string{common.BKFieldID},
			IDField:      []string{common.BKFieldID},
			AuditTypeField: &tools.AuditResType{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelUniqueRes,
			},
			AuditDataField: &tools.AuditDataField{
				ResIDField:   "id",
				ResNameField: "bk_obj_id",
			},
		},
	}
	typeTableMap = map[string]string{
		"association":        common.BKTableNameAsstDes,
		"object":             common.BKTableNameObjDes,
		"obj_attribute":      common.BKTableNameObjAttDes,
		"obj_association":    common.BKTableNameObjAsst,
		"obj_classification": common.BKTableNameObjClassification,
		"plat":               common.BKTableNameBasePlat,
		"property_group":     common.BKTableNamePropertyGroup,
		"service_category":   common.BKTableNameServiceCategory,
		"unique_keys":        common.BKTableNameObjUnique,
	}
	dataInitiator = map[string]func(kit *rest.Kit, db dal.RDB, table string, data []mapstr.MapStr,
		insertOps *tools.InsertOptions) error{
		common.BKTableNameObjAttDes:         initTemplateData,
		common.BKTableNameServiceCategory:   addServiceCategoryData,
		common.BKTableNameObjDes:            initTemplateData,
		common.BKTableNameAsstDes:           initTemplateData,
		common.BKTableNameObjAsst:           initTemplateData,
		common.BKTableNameObjClassification: initTemplateData,
		common.BKTableNameBasePlat:          initTemplateData,
		common.BKTableNamePropertyGroup:     initTemplateData,
		common.BKTableNameObjUnique:         addObjectUniqueData,
	}
	tableSeq = []string{
		common.BKTableNameObjAttDes,
		common.BKTableNameServiceCategory,
		common.BKTableNameObjDes,
		common.BKTableNameAsstDes,
		common.BKTableNameObjAsst,
		common.BKTableNameObjClassification,
		common.BKTableNameBasePlat,
		common.BKTableNamePropertyGroup,
		common.BKTableNameObjUnique,
	}
)

func initTemplateData(kit *rest.Kit, db dal.RDB, table string, data []mapstr.MapStr,
	insertOps *tools.InsertOptions) error {

	for _, item := range data {
		if _, exists := item[common.CreateTimeField]; exists {
			item[common.CreateTimeField] = time.Now()
		}
		if _, exists := item[common.LastTimeField]; exists {
			item[common.LastTimeField] = time.Now()
		}
	}

	dataInterface := make([]interface{}, len(data))
	for i, item := range data {
		dataInterface[i] = item
	}
	_, err := tools.InsertData(kit, db, table, dataInterface, insertOps)
	if err != nil {
		blog.Errorf("add init data for tenant %s failed, table: %s, err: %v", kit.TenantID, table, err)
		return err
	}
	return nil
}

func addServiceCategoryData(kit *rest.Kit, db dal.RDB, table string, data []mapstr.MapStr,
	insertOps *tools.InsertOptions) error {

	subCategory := make([]mapstr.MapStr, 0)
	parentCategory := make([]interface{}, 0)
	// sub id: parent id
	subIDParentNameMap := make(map[int]int)
	// parent id: parent name
	parentIDNameMap := make(map[int]string)
	ids := make([]int, 0)

	for _, item := range data {
		id, err := util.GetIntByInterface(item[common.BKFieldID])
		if err != nil {
			blog.Errorf("get id int from interface failed, err: %v", err)
			return err
		}
		parentID, err := util.GetIntByInterface(item[common.BKParentIDField])
		if err != nil {
			blog.Errorf("get parent id int from interface failed, err: %v", err)
			return err
		}
		if parentID == 0 {
			parentCategory = append(parentCategory, item)
			parentIDNameMap[id] = util.GetStrByInterface(item[common.BKFieldName])
			continue
		}
		subCategory = append(subCategory, item)
		subIDParentNameMap[id] = parentID
		ids = append(ids, id)

	}

	copiedInterface := deepcopy.Copy(*insertOps)
	copiedMap, ok := copiedInterface.(tools.InsertOptions)
	if !ok {
		blog.Errorf("failed to convert interface to tools.InsertOptions, rid: %s", kit.Rid)
		return fmt.Errorf("failed to convert interface to tools.InsertOptions")
	}

	copiedMap.IgnoreKeys = append(copiedMap.IgnoreKeys, common.BKRootIDField)
	copiedMap.IDField = append(copiedMap.IDField, common.BKRootIDField)

	parentIDs, err := tools.InsertData(kit, db, common.BKTableNameServiceCategory, parentCategory, &copiedMap)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}
	for key, value := range parentIDs {
		name := strings.Split(key, "*")[0]
		parentIDs[name] = value
	}

	// add sub category data
	subCategoryData := make([]interface{}, 0)
	for index, item := range subCategory {
		parentID, err := util.GetInt64ByInterface(parentIDs[parentIDNameMap[subIDParentNameMap[ids[index]]]])
		if err != nil {
			blog.Errorf("get parent id int64 failed, err: %v", err)
			return err
		}
		item[common.BKParentIDField] = parentID
		item[common.BKRootIDField] = parentID
		subCategoryData = append(subCategoryData, item)
	}

	_, err = tools.InsertData(kit, db, common.BKTableNameServiceCategory, subCategoryData, insertOps)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}
	return nil
}

func generateUniqueKey(objID, propertyID string) string {
	return objID + ":" + propertyID
}

func getUniqueKeys(kit *rest.Kit, db dal.RDB, objUniqueKeys map[string][][]string) ([]objectUnique, error) {
	attrArr := make([]metadata.Attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(nil).All(kit.Ctx, &attrArr)
	if err != nil {
		blog.Errorf("get host unique fields failed, err: %v", err)
		return nil, err
	}

	if len(attrArr) == 0 {
		blog.Errorf("get object attribute failed")
		return nil, fmt.Errorf("get object attribute failed")
	}

	attrIDMap := make(map[string]uint64)
	for _, attr := range attrArr {
		attrIDMap[generateUniqueKey(attr.ObjectID, attr.PropertyID)] = uint64(attr.ID)
	}
	uniqueKeys := make([]objectUnique, 0)
	for objID, value := range objUniqueKeys {
		for _, property := range value {
			keys := make([]uniqueKey, 0)
			for _, field := range property {
				keys = append(keys, uniqueKey{
					Kind: "property",
					ID:   attrIDMap[generateUniqueKey(objID, field)],
				})
			}
			uniqueKeys = append(uniqueKeys, objectUnique{
				Keys:     keys,
				ObjID:    objID,
				IsPre:    true,
				LastTime: time.Now(),
			})
		}
	}

	return uniqueKeys, nil
}

func addObjectUniqueData(kit *rest.Kit, db dal.RDB, table string, data []mapstr.MapStr,
	insertOps *tools.InsertOptions) error {

	if len(data) == 0 {
		blog.Errorf("get data failed, err: %v")
		return fmt.Errorf("get data failed")
	}
	uniqueKeys := data[0]
	objUniqueKeys := make(map[string][][]string, 0)
	for key, item := range uniqueKeys {
		switch v := item.(type) {
		case primitive.A:
			subArr := make([][]string, 0)
			for _, subItem := range v {
				innerArr, err := convertToArryString(subItem)
				if err != nil {
					blog.Errorf("convert to array string failed, err: %v", err)
					return err
				}
				subArr = append(subArr, innerArr)
			}
			objUniqueKeys[key] = subArr
		default:
			blog.Errorf("invalid type, %v", item)
			return fmt.Errorf("invalid type")
		}
	}

	uniqueKeysArr, err := getUniqueKeys(kit, db, objUniqueKeys)
	if err != nil {
		blog.Errorf("get unique keys failed, err: %v", err)
		return err
	}

	objUniqueData := make([]interface{}, 0)
	for _, key := range uniqueKeysArr {
		objUniqueData = append(objUniqueData, key)
	}

	_, err = tools.InsertData(kit, db, table, objUniqueData, insertOps)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", table, err)
		return err
	}

	return nil
}

func convertToArryString(data interface{}) ([]string, error) {
	switch t := data.(type) {
	case primitive.A:
		innerArr := make([]string, 0)
		for _, d := range t {
			switch ty := d.(type) {
			case string:
				innerArr = append(innerArr, ty)
			default:
				blog.Errorf("invalid type, %T", ty)
				return nil, fmt.Errorf("invalid type, %T", ty)
			}
		}
		return innerArr, nil
	default:
		blog.Errorf("invalid type, %T", t)
		return nil, fmt.Errorf("invalid type, %T", t)
	}
}

type objectUnique struct {
	ID       uint64      `bson:"id"`
	ObjID    string      `bson:"bk_obj_id"`
	Keys     []uniqueKey `bson:"keys"`
	IsPre    bool        `bson:"ispre"`
	LastTime time.Time   `bson:"last_time"`
}

type uniqueKey struct {
	Kind string `bson:"key_kind"`
	ID   uint64 `bson:"key_id"`
}
