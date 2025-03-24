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

	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
)

var ignoreKeysArr = []string{"create_time", "last_time", "_id"}

// InsertData insert data for upgrade
func InsertData(kit *rest.Kit, db local.DB, table string, data []mapstr.MapStr, insertOps *InsertOptions) (
	map[string]interface{}, error) {

	result := make([]mapstr.MapStr, 0)
	err := db.Table(table).Find(mapstr.MapStr{}).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exist data failed, table: %s, err: %v", table, err)
		return nil, err
	}

	idFields := make(map[string]interface{})
	if len(insertOps.IDField) > 0 {
		for _, item := range result {
			valueStr := getUniqueStr(item, insertOps.UniqueFields)
			if valueStr == "" {
				continue
			}
			if id, ok := item[insertOps.IDField[0]]; ok {
				idFields[valueStr] = id
			}
		}
	}

	existData := make(map[string]mapstr.MapStr)
	for _, item := range result {
		valueStr := getUniqueStr(item, insertOps.UniqueFields)
		existData[valueStr] = item
	}

	insertData, err := getInsertData(existData, data, insertOps)
	if err != nil {
		return nil, err
	}

	if len(insertData) == 0 {
		blog.Infof("no data to insert, table: %s", table)
		return idFields, nil
	}

	if len(insertOps.IDField) > 0 {
		nextIDs, err := mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
		if err != nil {
			blog.Errorf("get next %d data IDs failed, table: %s, err: %v", len(insertData), table, err)
			return nil, fmt.Errorf("get next %d data IDs failed, table: %s, err: %v", len(insertData), table, err)
		}
		for index, value := range insertData {
			for _, idField := range insertOps.IDField {
				value[idField] = nextIDs[index]
			}
			idFields[getUniqueStr(value, insertOps.UniqueFields)] = nextIDs[index]
		}
	}

	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("add data for table %s failed, data: %+v, err: %v", insertData, err)
		return nil, err
	}

	// add audit log
	if insertOps.AuditDataField == nil {
		return idFields, nil
	}
	auditField := &AuditStruct{
		AuditDataField: insertOps.AuditDataField,
		AuditTypeData:  insertOps.AuditTypeField,
	}

	if err = AddCreateAuditLog(kit, db, insertData, auditField); err != nil {
		blog.Errorf("add audit log failed, err: %v", err)
		return nil, err
	}
	return idFields, nil
}

// SplitUniqueStr split unique str by unique fields
func SplitUniqueStr(uniqueStr string, uniqueFields []string) (map[string]string, error) {
	valueMap := make(map[string]string, len(uniqueFields))

	arr := strings.Split(uniqueStr, "*")
	if len(arr) != len(uniqueFields) {
		return nil, fmt.Errorf("unique str %s is invalid", uniqueStr)
	}

	for index, uniqueField := range uniqueFields {
		valueMap[uniqueField] = arr[index]
	}

	return valueMap, nil
}

func getUniqueStr(item mapstr.MapStr, uniqueFields []string) string {
	var strArr []string
	for _, uniqueValue := range uniqueFields {
		if _, ok := item[uniqueValue]; !ok {
			continue
		}
		str := util.GetStrByInterface(item[uniqueValue])
		strArr = append(strArr, str)
	}
	return strings.Join(strArr, "*")
}

func getInsertData(existData map[string]mapstr.MapStr, data []mapstr.MapStr, compareFiled *InsertOptions) (
	[]map[string]interface{}, error) {

	insertData := make([]map[string]interface{}, 0)
	for _, item := range data {
		mapStrData, err := InterfaceToMapStr(item)
		if err != nil {
			blog.Errorf("interface to mapStr failed, err: %v, data: %+v", err, item)
			return nil, err
		}
		valueStr := getUniqueStr(mapStrData, compareFiled.UniqueFields)
		if _, exist := existData[valueStr]; !exist {
			insertData = append(insertData, mapStrData)
			continue
		}

		if err = cmpData(mapStrData, existData[valueStr], compareFiled.IgnoreKeys); err != nil {
			return nil, err
		}
	}
	return insertData, nil
}

func cmpData(data mapstr.MapStr, existData mapstr.MapStr, ignoreKeys []string) error {
	ignoreKeys = append(ignoreKeys, ignoreKeysArr...)
	for _, key := range ignoreKeys {
		delete(existData, key)
		delete(data, key)
	}

	var err error
	data, err = dataNumericConvert(data)
	if err != nil {
		blog.Errorf("data numeric convert error, data: %+v, err: %v", data, err)
		return err
	}
	existData, err = dataNumericConvert(existData)
	if err != nil {
		blog.Errorf("data numeric convert error, data: %+v, err: %v", existData, err)
		return err
	}

	if !cmp.Equal(data, existData) {
		blog.Errorf("the data in database is different from the data to be inserted, existData: %+v, insertData: %+v",
			existData, data)
		return fmt.Errorf("data in database is different from the data to be inserted")
	}
	return nil
}

func dataNumericConvert(data mapstr.MapStr) (mapstr.MapStr, error) {
	for key, value := range data {
		if !util.IsNumeric(value) {
			continue
		}
		valueInt64, err := util.GetInt64ByInterface(value)
		if err != nil {
			blog.Errorf("get value int64 error, key: %s, value: %v, type: %T, err: %v", key, value, value, err)
			return nil, err
		}
		data[key] = valueInt64
	}
	return data, nil
}

// InterfaceToMapStr interface to mapstr
func InterfaceToMapStr(data interface{}) (map[string]interface{}, error) {

	resultData := make(map[string]interface{})
	switch value := data.(type) {
	case map[string]interface{}:
		return value, nil
	default:
		out, err := bson.Marshal(data)
		if err != nil {
			blog.Errorf("marshal error %v, data: %v", err, data)
			return nil, fmt.Errorf("marshal error %v", err)
		}
		if err = bson.Unmarshal(out, &resultData); err != nil {
			blog.Errorf("marshal error %v, data: %v", err, data)
			return nil, fmt.Errorf("unmarshal error %v", err)
		}
	}

	return resultData, nil
}

// InsertOptions the options of insert field for audit and data
type InsertOptions struct {
	UniqueFields   []string
	IgnoreKeys     []string
	IDField        []string
	AuditDataField *AuditDataField `bson:",inline"`
	AuditTypeField *AuditResType   `bson:",inline"`
}

// ConvStructToMap convert struct to map
func ConvStructToMap(obj interface{}) (map[string]interface{}, error) {
	if data, isOk := obj.(map[string]interface{}); isOk {
		return data, nil
	}

	data, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = bson.Unmarshal(data, &result)
	return result, err
}

// IDOptions the options of data template id
type IDOptions struct {
	ResNameField string
	RemoveKeys   []string
}
