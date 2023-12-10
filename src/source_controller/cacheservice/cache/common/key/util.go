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

package key

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseMapStrData(data interface{}) (mapstr.MapStr, error) {
	switch t := data.(type) {
	case *mapstr.MapStr:
		if t == nil {
			return nil, errors.New("data is nil")
		}
		return *t, nil
	case *map[string]interface{}:
		if t == nil {
			return nil, errors.New("data is nil")
		}
		return *t, nil
	case mapstr.MapStr:
		return t, nil
	case map[string]interface{}:
		return t, nil
	default:
		return nil, errors.New("data is not of mapstr.MapStr type")
	}
}

func commonIDGenerator(data interface{}, field string) (string, float64, error) {
	commonData, err := parseMapStrData(data)
	if err != nil {
		return "", 0, err
	}

	id, err := util.GetInt64ByInterface(commonData[field])
	if err != nil {
		return "", 0, err
	}

	return strconv.FormatInt(id, 10), float64(id), nil
}

func commonOidGenerator(data interface{}) (string, float64, error) {
	commonData, err := parseMapStrData(data)
	if err != nil {
		return "", 0, err
	}

	_, exists := commonData[common.MongoMetaID]
	if !exists {
		return "", 0, errors.New("data oid field is not exists")
	}

	switch t := commonData[common.MongoMetaID].(type) {
	case primitive.ObjectID:
		return t.Hex(), 0, nil
	case string:
		return t, 0, nil
	default:
		return "", 0, errors.New("data oid field is invalid")
	}
}

func commonKeyGenerator(data interface{}, fields ...string) (string, error) {
	if dataStr, ok := data.(string); ok {
		return dataStr, nil
	}

	commonData, err := parseMapStrData(data)
	if err != nil {
		return "", err
	}

	var keys []string
	for _, field := range fields {
		keys = append(keys, util.GetStrByInterface(commonData[field]))
	}

	key := strings.Join(keys, "||")
	if len(key) == 0 {
		return "", errors.New("key is empty")
	}

	return key, nil
}

func commonIDDataGetter(db dal.DB, table, idField string, keys ...string) ([]interface{}, error) {
	ids := make([]int64, len(keys))
	var err error

	for i, key := range keys {
		ids[i], err = strconv.ParseInt(key, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	cond := mapstr.MapStr{
		idField: mapstr.MapStr{common.BKDBIN: ids},
	}

	return getCommonDBDataByCond(db, table, cond)
}

func commonKeyDataGetter(db dal.DB, table, field string, keys ...string) ([]interface{}, error) {
	cond := mapstr.MapStr{
		field: mapstr.MapStr{common.BKDBIN: keys},
	}

	return getCommonDBDataByCond(db, table, cond)
}

func commonOidDataGetter(db dal.DB, table string, keys ...string) ([]interface{}, error) {
	oids := make([]primitive.ObjectID, len(keys))
	for i, key := range keys {
		oid, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			return nil, err
		}
		oids[i] = oid
	}

	cond := mapstr.MapStr{
		common.MongoMetaID: mapstr.MapStr{common.BKDBIN: oids},
	}

	opts := types.NewFindOpts().SetWithObjectID(true)
	return getCommonDBDataByCond(db, table, cond, opts)
}

func commonKeyWithOidGetter(db dal.DB, table, field string, keys ...string) ([]interface{}, error) {
	cond := mapstr.MapStr{
		field: mapstr.MapStr{common.BKDBIN: keys},
	}

	opts := types.NewFindOpts().SetWithObjectID(true)
	return getCommonDBDataByCond(db, table, cond, opts)
}

func getCommonDBDataByCond(db dal.DB, table string, cond mapstr.MapStr, opts ...*types.FindOpts) ([]interface{},
	error) {

	dataArr := make([]mapstr.MapStr, 0)
	err := db.Table(table).Find(cond, opts...).All(context.Background(), &dataArr)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(dataArr))
	for i, data := range dataArr {
		result[i] = data
	}

	return result, nil
}
