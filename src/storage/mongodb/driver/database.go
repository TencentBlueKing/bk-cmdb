/*
* Tencent is pleased to support the open source community by making 蓝鲸 available.
* Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
* Licensed under the MIT License (the ",License"); you may not use this file except
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

	"configcenter/src/common/mapstr"
	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

var _ mongodb.Database = (*database)(nil)

type database struct {
	innerDatabase *mongo.Database
}

func newDatabase(db *mongo.Database) *database {
	return &database{innerDatabase: db}
}

func (d *database) Drop() error {
	return d.innerDatabase.Drop(context.TODO())
}

func (d *database) Name() string {
	return d.innerDatabase.Name()
}

func (d *database) HasCollection(collName string) (bool, error) {

	cursor, err := d.innerDatabase.ListCollections(context.TODO(),
		bsonx.Doc{
			bsonx.Elem{Key: "type", Value: bsonx.String("collection")},
			bsonx.Elem{Key: "name", Value: bsonx.String(collName)},
		},
	)

	if nil != err {
		return false, err
	}

	for cursor.Next(context.TODO()) {
		return true, nil
	}

	return false, nil
}

func (d *database) DropCollection(collName string) error {
	return d.innerDatabase.Collection(collName).Drop(context.TODO())
}

func (d *database) CreateEmptyCollection(collName string) error {
	ret := d.innerDatabase.RunCommand(context.TODO(), map[string]interface{}{"create": collName})
	return ret.Err()
}

func (d *database) GetCollectionNames() ([]string, error) {

	cursor, err := d.innerDatabase.ListCollections(context.TODO(),
		bsonx.Doc{
			bsonx.Elem{Key: "type", Value: bsonx.String("collection")},
		},
	)

	if nil != err {
		return nil, err
	}

	collNames := []string{}
	for cursor.Next(context.TODO()) {

		result := mapstr.New()
		if err := cursor.Decode(&result); nil != err {
			return nil, err
		}

		name, err := result.String("name")
		if nil != err {
			return nil, err
		}

		collNames = append(collNames, name)
	}

	return collNames, nil
}
