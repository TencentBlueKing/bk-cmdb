/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package index TODO
package index

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/index/collections"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InstanceIndexes TODO
func InstanceIndexes() []types.Index {
	return instanceDefaultIndexes
}

// InstanceAssociationIndexes TODO
func InstanceAssociationIndexes() []types.Index {
	return associationDefaultIndexes
}

// CCFieldTypeToDBType TODO
func CCFieldTypeToDBType(typ string) string {
	switch typ {
	case common.FieldTypeSingleChar, common.FieldTypeEnum, common.FieldTypeDate, common.FieldTypeList:
		return "string"
	case common.FieldTypeInt, common.FieldTypeFloat:
		return "number"
	}

	// other type not support
	return ""
}

// GetUniqueIndexNameByID TODO
func GetUniqueIndexNameByID(id uint64) string {
	return fmt.Sprintf("%s%d", common.CCLogicUniqueIdxNamePrefix, id)
}

// DeprecatedIndexName 获取表中没有规范化前，所有索引的名字, map[collection name][]string{"索引名字"}
func DeprecatedIndexName() map[string][]string {

	return collections.DeprecatedIndexName()
}

// TableIndexes TODO
func TableIndexes() map[string][]types.Index {

	return collections.Indexes()
}

// ToDBUniqueIndex TODO
func ToDBUniqueIndex(objID string, id uint64, keys []metadata.UniqueKey,
	properties []metadata.Attribute) (types.Index, errors.CCErrorCoder) {
	dbIndex := types.Index{
		Background:              true,
		Unique:                  true,
		Name:                    GetUniqueIndexNameByID(id),
		Keys:                    make(bson.D, 0),
		PartialFilterExpression: make(map[string]interface{}),
	}
	propertiesIDMap := make(map[int64]metadata.Attribute, len(properties))
	for _, property := range properties {
		propertiesIDMap[property.ID] = property
	}

	keyLen := len(keys)
	for _, key := range keys {
		attr := propertiesIDMap[int64(key.ID)]
		if objID == common.BKInnerObjIDHost && attr.PropertyID == common.BKCloudIDField {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未foreignkey 特殊的类型
			attr.PropertyType = common.FieldTypeInt
		}
		if objID == common.BKInnerObjIDHost &&
			(attr.PropertyID == common.BKHostInnerIPField || attr.PropertyID == common.BKHostOuterIPField ||
				attr.PropertyID == common.BKOperatorField || attr.PropertyID == common.BKBakOperatorField ||
				attr.PropertyID == common.BKHostInnerIPv6Field || attr.PropertyID == common.BKHostOuterIPv6Field) {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未innerIP,OuterIP 特殊的类型
			attr.PropertyType = common.FieldTypeList
		}

		if !ValidateCCFieldType(attr.PropertyType, keyLen) {
			return dbIndex, errors.GetGlobalCCError().CreateDefaultCCErrorIf(string(common.English)).
				CCErrorf(common.CCErrCoreServiceUniqueIndexPropertyType, attr.PropertyID)
		}

		dbType := CCFieldTypeToDBType(attr.PropertyType)
		if dbType == "" {
			return dbIndex, errors.GetGlobalCCError().CreateDefaultCCErrorIf(string(common.English)).
				CCErrorf(common.CCErrCoreServiceUniqueIndexPropertyType, attr.PropertyID)
		}

		dbIndex.Keys = append(dbIndex.Keys, primitive.E{
			Key:   attr.PropertyID,
			Value: 1,
		})
		dbIndex.PartialFilterExpression[attr.PropertyID] = map[string]interface{}{common.BKDBType: dbType}
	}

	return dbIndex, nil
}

// ValidateCCFieldType returns if cc unique field type is valid, differs for union and separate unique. issue #5240
func ValidateCCFieldType(propertyType string, keyLen int) bool {
	if keyLen == 1 {
		switch propertyType {
		case common.FieldTypeSingleChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeList:
			return true
		default:
			return false
		}
	}

	switch propertyType {
	case common.FieldTypeSingleChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum,
		common.FieldTypeDate, common.FieldTypeList:
		return true
	default:
		return false
	}
}
