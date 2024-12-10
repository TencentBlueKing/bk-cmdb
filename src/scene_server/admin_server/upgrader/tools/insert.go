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

	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
	"github.com/google/go-cmp/cmp"

	"go.mongodb.org/mongo-driver/bson"
)

var ignoreKeysArr = []string{"create_time", "last_time", "_id"}

// InsertData insert data for upgrade
func InsertData(kit *rest.Kit, db dal.RDB, table string, data []interface{}, compareFiled *CmpFiled,
	auditField []AuditType, audit *AuditDataField) (map[string]interface{}, error) {

	result := make([]mapstr.MapStr, 0)
	err := db.Table(table).Find(mapstr.MapStr{}).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exit data failed, table: %s, err: %v", table, err)
		return nil, err
	}

	idFields := getUniqueField(result, compareFiled)
	exitData := map[string]mapstr.MapStr{}
	for _, item := range result {
		valueStr := getUniqueStr(item, compareFiled)
		if valueStr != "" {
			exitData[valueStr[1:]] = item
		}
	}

	insertData, insertAuditData, err := getInsertData(exitData, data, auditField, compareFiled)
	if err != nil {
		return nil, err
	}

	if len(insertData) == 0 {
		blog.Infof("no data to insert, table: %s", table)
		return idFields, nil
	}

	nextIDs := make([]uint64, 0)
	if len(compareFiled.IDField) > 0 {
		nextIDs, err = mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
		if err != nil {
			return nil, fmt.Errorf("init service template, get next set ids error, err: %v", err)
		}
		for index, value := range insertData {
			value[compareFiled.IDField] = nextIDs[index]
			idFields[getUniqueStr(value, compareFiled)] = nextIDs[index]
		}
	}

	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("add data for table %s failed, data: %+v, err: %v", insertData, err)
		return nil, err
	}

	// add audit log
	if audit == nil {
		return idFields, nil
	}
	if err = AddCreateAuditLog(kit, mongodb.Dal().Shard(kit.ShardOpts()), insertData, insertAuditData,
		audit); err != nil {
		blog.Errorf("add audit log failed, err: %v", err)
		return nil, err
	}
	return idFields, nil
}

func getUniqueStr(item mapstr.MapStr, compareFiled *CmpFiled) string {
	valueStr := ""
	for _, uniqueValue := range compareFiled.UniqueFields {
		if _, ok := item[uniqueValue]; !ok {
			continue
		}
		str := util.GetStrByInterface(item[uniqueValue])
		valueStr = valueStr + "*" + str
	}
	if valueStr == "" {
		return valueStr
	}
	return valueStr[1:]
}

func getInsertData(existData map[string]mapstr.MapStr, data []interface{},
	auditField []AuditType, compareFiled *CmpFiled) ([]map[string]interface{}, []AuditType, error) {

	insertData := make([]map[string]interface{}, 0)
	insertAuditData := make([]AuditType, 0)
	uniqueIDs := map[string]interface{}{}
	for index, item := range data {
		mapStrData, err := InterfaceToMapStr(item)
		if err != nil {
			return nil, nil, err
		}
		valueStr := ""
		for _, uniqueValue := range compareFiled.UniqueFields {
			str := util.GetStrByInterface(mapStrData[uniqueValue])
			valueStr = valueStr + "*" + str
		}
		valueStr = valueStr[1:]
		if _, exit := existData[valueStr]; !exit {
			insertData = append(insertData, mapStrData)
			if len(auditField) > 0 {
				insertAuditData = append(insertAuditData, auditField[index])
			}
			if id, ok := mapStrData[compareFiled.IDField]; !ok {
				uniqueIDs[valueStr] = id
			}
			continue
		}

		if err = cmpData(mapStrData, existData[valueStr], compareFiled.IgnoreKeys); err != nil {
			return nil, nil, err
		}
	}
	return insertData, insertAuditData, nil
}

func cmpData(data mapstr.MapStr, exitData mapstr.MapStr, ignoreKeys []string) error {
	ignoreKeys = append(ignoreKeys, ignoreKeysArr...)
	for _, key := range ignoreKeys {
		delete(exitData, key)
		delete(data, key)
	}
	for key, value := range data {
		valueInt64, err := util.GetInt64ByInterface(value)
		if err != nil {
			continue
		}
		data[key] = valueInt64
	}
	for key, value := range exitData {
		valueInt64, err := util.GetInt64ByInterface(value)
		if err != nil {
			continue
		}
		exitData[key] = valueInt64
	}

	if !cmp.Equal(data, exitData) {
		blog.Errorf("the data in database is different from the data to be inserted, exitData: %+v, insertData: %+v",
			exitData, data)
		return fmt.Errorf("data in database is different from the data to be inserted")
	}
	return nil
}

func getUniqueField(result []mapstr.MapStr, compareFiled *CmpFiled) map[string]interface{} {
	idsMap := map[string]interface{}{}
	for _, item := range result {
		valueStr := getUniqueStr(item, compareFiled)
		if valueStr != "" {
			if id, ok := item[compareFiled.IDField]; !ok {
				idsMap[valueStr[1:]] = id
			}

		}
	}
	return idsMap
}

func getExitDataMap(result []mapstr.MapStr, unique []string) (map[string]mapstr.MapStr, error) {
	exitData := map[string]mapstr.MapStr{}
	for _, item := range result {
		valueStr := ""
		for _, uniqueValue := range unique {
			if _, ok := item[uniqueValue]; !ok {
				continue
			}
			str := util.GetStrByInterface(item[uniqueValue])
			valueStr = valueStr + "*" + str
		}
		valueStr = valueStr[1:]
		if valueStr != "" {
			exitData[valueStr[1:]] = item
		}
	}
	return exitData, nil
}

// InterfaceToMapStr interface to mapstr
func InterfaceToMapStr(data interface{}) (map[string]interface{}, error) {

	resultData := map[string]interface{}{}
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

// CmpFiled the compare filed
type CmpFiled struct {
	UniqueFields []string
	IgnoreKeys   []string
	IDField      string
}
