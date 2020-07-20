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

package upgrader

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"

	"gopkg.in/mgo.v2/bson"
)

// Upsert inset row but update it without ignores key if exists same value with keys
func Upsert(ctx context.Context, db dal.RDB, tableName string, row interface{}, idField string, uniqueKeys []string, ignores []string) (instID uint64, preData map[string]interface{}, err error) {
	data := map[string]interface{}{}
	switch value := row.(type) {
	case map[string]interface{}:
		data = value
	default:
		out, err := bson.Marshal(row)
		if err != nil {
			return 0, nil, fmt.Errorf("marshal error %v", err)
		}
		if err = bson.Unmarshal(out, data); err != nil {
			return 0, nil, fmt.Errorf("unmarshal error %v", err)
		}
	}

	condition := map[string]interface{}{}
	for _, key := range uniqueKeys {
		condition[key] = data[key]
	}

	existOne := map[string]interface{}{}
	err = db.Table(tableName).Find(condition).One(ctx, &existOne)

	if db.IsNotFoundError(err) {
		if "" != idField {
			instID, err = db.NextSequence(ctx, tableName)
			if err != nil {
				return 0, nil, err
			}
			data[idField] = instID
		}
		// blog.Infof("%s insert %v", tablename, data)
		err = db.Table(tableName).Insert(ctx, data)
		if err != nil {
			return instID, nil, fmt.Errorf("insert error %v", err)
		}
		return instID, nil, nil
	}
	if nil != err {
		return 0, nil, fmt.Errorf("find error %v", err)
	}

	ignoreSet := map[string]bool{idField: true}
	if "" != idField {
		switch id := existOne[idField].(type) {
		case nil:
			return 0, nil, errors.New("there is no " + idField + " field in table " + tableName)
		case int:
			instID = uint64(id)
		case int16:
			instID = uint64(id)
		case int32:
			instID = uint64(id)
		case int64:
			instID = uint64(id)
		case float32:
			instID = uint64(id)
		case float64:
			instID = uint64(id)
		}
		if instID <= 0 {
			instID, err = db.NextSequence(ctx, tableName)
			if err != nil {
				return 0, nil, fmt.Errorf("get NextSequence error %v", err)
			}
			data[idField] = instID
			delete(ignoreSet, idField)
			blog.Infof("reset %s %s to %d", tableName, idField, instID)
		}
	}

	for _, key := range ignores {
		ignoreSet[key] = true
	}
	newData := map[string]interface{}{}
	for key, value := range data {
		if ignoreSet[key] {
			continue
		}
		newData[key] = value
	}
	// blog.Infof("%s update %v", tablename, newData)
	err = db.Table(tableName).Update(ctx, condition, newData)
	if err != nil {
		return instID, existOne, fmt.Errorf("update error %v", err)
	}
	return instID, existOne, nil
}

// Insert insert the row to db if it doesn't exist, else do nothing
// row is the data to insert
// idField used to be field name of generated instance id
// uniqueKeys used to judge whether the row exist in db
func Insert(ctx context.Context, db dal.RDB, tableName string, row interface{}, idField string, uniqueKeys []string) error {
	if idField == "" {
		blog.Errorf("idField is empty, it can't be empty")
		return errors.New("idField can't be empty")
	}

	data := map[string]interface{}{}
	switch value := row.(type) {
	case map[string]interface{}:
		data = value
	default:
		out, err := bson.Marshal(row)
		if err != nil {
			blog.Errorf("marshal error:%v", err)
			return err
		}
		if err = bson.Unmarshal(out, data); err != nil {
			blog.Errorf("unmarshal error:%v", err)
			return err
		}
	}

	condition := map[string]interface{}{}
	for _, key := range uniqueKeys {
		condition[key] = data[key]
	}

	count, err := db.Table(tableName).Find(condition).Count(ctx)
	if err != nil {
		blog.Errorf("find count error:%v", err)
		return err
	}
	// if exist, return directly
	if count > 0 {
		return nil
	}

	instID, err := db.NextSequence(ctx, tableName)
	if err != nil {
		return err
	}
	data[idField] = instID

	err = db.Table(tableName).Insert(ctx, data)
	if err != nil {
		blog.Errorf("insert error %v", err)
		return err
	}

	return nil
}
