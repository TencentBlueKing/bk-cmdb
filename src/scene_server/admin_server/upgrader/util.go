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

	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
)

// Upsert inset row but updata it whitout ignores key if exists same value with keys
func Upsert(ctx context.Context, db dal.RDB, tablename string, row interface{}, idfieldname string, keys []string, ignores []string) (instID uint64, preData map[string]interface{}, err error) {
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
	for _, key := range keys {
		condition[key] = data[key]
	}

	existOne := map[string]interface{}{}
	err = db.Table(tablename).Find(condition).One(ctx, &existOne)

	if db.IsNotFoundError(err) {
		if "" != idfieldname {
			instID, err = db.NextSequence(ctx, tablename)
			if err != nil {
				return 0, nil, err
			}
			data[idfieldname] = instID
		}
		// blog.Infof("%s insert %v", tablename, data)
		err = db.Table(tablename).Insert(ctx, data)
		if err != nil {
			return instID, nil, fmt.Errorf("insert error %v", err)
		}
		return instID, nil, nil
	}
	if nil != err {
		return 0, nil, fmt.Errorf("Find error %v", err)
	}

	ignoreset := map[string]bool{idfieldname: true}
	if "" != idfieldname {
		switch id := existOne[idfieldname].(type) {
		case nil:
			return 0, nil, errors.New("there is no " + idfieldname + " field in table " + tablename)
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
			instID, err = db.NextSequence(ctx, tablename)
			if err != nil {
				return 0, nil, fmt.Errorf("get NextSequence error %v", err)
			}
			data[idfieldname] = instID
			delete(ignoreset, idfieldname)
			blog.Infof("reset %s %s to %d", tablename, idfieldname, instID)
		}
	}

	for _, key := range ignores {
		ignoreset[key] = true
	}
	newData := map[string]interface{}{}
	for key, value := range data {
		if ignoreset[key] == true {
			continue
		}
		newData[key] = value
	}
	// blog.Infof("%s update %v", tablename, newData)
	err = db.Table(tablename).Update(ctx, condition, newData)
	if err != nil {
		return instID, existOne, fmt.Errorf("update error %v", err)
	}
	return instID, existOne, nil
}
