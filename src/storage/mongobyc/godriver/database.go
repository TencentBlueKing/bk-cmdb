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

package godriver

import (
	"context"

	"configcenter/src/storage/mongobyc"

	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

var _ mongobyc.Database = (*database)(nil)

type database struct {
	innerDatabase *mgo.Database
}

func newDatabase(db *mgo.Database) *database {
	return &database{innerDatabase: db}
}

func (d *database) Drop() error {
	return d.innerDatabase.Drop(context.TODO())
}

func (d *database) Name() string {
	return d.innerDatabase.Name()
}

func (d *database) HasCollection(collName string) (bool, error) {

	return false, nil
}

func (d *database) DropCollection(collName string) error {
	return nil
}

func (d *database) CreateEmptyCollection(collName string) error {
	return nil
}

func (d *database) GetCollectionNames() ([]string, error) {
	return nil, nil
}
