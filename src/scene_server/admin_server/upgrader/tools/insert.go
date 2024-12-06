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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
	"github.com/google/go-cmp/cmp"

	"go.mongodb.org/mongo-driver/bson"
)

var ignoreKeysArr = []string{"ID", "create_time", "last_time", "_id", "bk_property_index"}

// InsertData insert data for upgrade
func InsertData(kit *rest.Kit, db dal.RDB, table string, data []interface{}, compareFiled *CmpFiled,
	auditField []AuditField) ([]uint64, error) {

	result := make([]mapstr.MapStr, 0)
	err := db.Table(table).Find(mapstr.MapStr{}).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("find exit data failed, table: %s, err: %v", table, err)
		return nil, err
	}

	exitData, err := getExitDataMap(result, compareFiled.Unique)
	if err != nil {
		return nil, err
	}
	insertData, insertAuditData, err := getInsertData(exitData, data, auditField, compareFiled.IgnoreKeys,
		compareFiled.Unique)
	if err != nil {
		return nil, err
	}

	if len(insertData) == 0 {
		ids := make([]uint64, 0)
		if id, ok := result[0][compareFiled.ID]; ok {
			dataID, err := util.GetIntByInterface(id)
			if err != nil {
				blog.Errorf("failed to get id, err: %v", err)
				return nil, err
			}
			ids = append(ids, uint64(dataID))
		}
		blog.Infof("no data to insert, table: %s", table)
		return ids, nil
	}

	nextSetIDs := make([]uint64, 0)
	if len(compareFiled.ID) > 0 {
		nextSetIDs, err = mongodb.Dal().Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, len(insertData))
		if err != nil {
			return nil, fmt.Errorf("init service template, get next set ids error, err: %v", err)
		}
		for index, value := range insertData {
			value[compareFiled.ID] = nextSetIDs[index]
			if len(insertAuditData) > 0 {
				insertAuditData[index].ResourceID = value[compareFiled.ID]
			}
		}
	}

	if err = db.Table(table).Insert(kit.Ctx, insertData); err != nil {
		blog.Errorf("add data for table %s failed, data: %+v, err: %v", insertData, err)
		return nil, err
	}

	// add audit log
	for index, audit := range insertAuditData {
		if dataID, ok := insertData[index][common.BKAppIDField]; ok {
			audit.BusinessID, err = util.GetInt64ByInterface(dataID)
			if err != nil {
				blog.Errorf("failed to convert data to int64, err: %v", err)
				return nil, err
			}
		}
		if err = AddCreateAuditLog(kit, mongodb.Dal().Shard(kit.ShardOpts()), insertData[index], &audit); err != nil {
			blog.Errorf("add audit log failed, err: %v", err)
			return nil, err
		}
	}
	return nextSetIDs, nil
}

func getInsertData(exitData map[string]mapstr.MapStr, data []interface{}, auditField []AuditField,
	ignoreKeys []string, unique []string) ([]mapstr.MapStr, []AuditField, error) {

	insertData := make([]mapstr.MapStr, 0)
	insertAuditData := make([]AuditField, 0)
	for index, item := range data {
		mapStrData, err := interfaceToMapStr(item)
		if err != nil {
			return nil, nil, err
		}
		valueStr := ""
		for _, uniqueValue := range unique {
			str := util.GetStrByInterface(mapStrData[uniqueValue])
			if err != nil {
				return nil, nil, err
			}
			valueStr = valueStr + "*" + str
		}
		if _, exit := exitData[valueStr]; !exit {
			insertData = append(insertData, mapStrData)
			if len(auditField) > 0 {
				insertAuditData = append(insertAuditData, auditField[index])
			}
			continue
		}

		if err = cmpData(mapStrData, exitData[valueStr], ignoreKeys); err != nil {
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

	if !cmp.Equal(data, exitData) {
		blog.Errorf("the data in database is different from the data to be inserted, exitData: %+v, insertData: %+v",
			exitData, data)
		return fmt.Errorf("data in database is different from the data to be inserted")
	}
	return nil
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
		if valueStr != "" {
			exitData[valueStr] = item
		}
	}
	return exitData, nil
}

func interfaceToMapStr(data interface{}) (map[string]interface{}, error) {

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

func toInt64(data interface{}) (int64, error) {
	switch value := data.(type) {
	case int:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int64:
		return value, nil
	case uint64:
		return int64(value), nil
	default:
		blog.Errorf("the data type is not supported, data: %v", data)
		return 0, fmt.Errorf("the data type is not supported, data: %v", data)
	}
}

// CmpFiled the compare filed
type CmpFiled struct {
	Unique     []string
	IgnoreKeys []string
	ID         string
}
