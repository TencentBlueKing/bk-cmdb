/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package driver

import (
	"context"

	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func (c *collection) parseIndexResult(ctx context.Context, cursor *mongo.Cursor) (*mongodb.QueryIndexResult, error) {

	// this struct from mongodb go driver about index
	type index struct {
		Key  map[string]int
		NS   string
		Name string
	}

	returnIndexResult := &mongodb.QueryIndexResult{}

	for cursor.Next(ctx) {
		elem := index{}
		if err := cursor.Decode(&elem); err != nil {
			return returnIndexResult, err
		}
		idxResult := mongodb.IndexResult{Name: elem.Name, Namespace: elem.NS}
		for name := range elem.Key {
			idxResult.Key = append(idxResult.Key, name)
		}
		returnIndexResult.Indexes = append(returnIndexResult.Indexes, idxResult)
	}

	if err := cursor.Err(); err != nil {
		return returnIndexResult, err
	}

	return returnIndexResult, nil
}
