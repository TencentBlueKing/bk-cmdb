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
	"configcenter/src/storage/mongobyc/deleteopt"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/mongobyc/insertopt"
	"configcenter/src/storage/mongobyc/replaceopt"
	"configcenter/src/storage/mongobyc/updateopt"
)

var _ mongobyc.CollectionInterface = (*collection)(nil)

type collection struct {
	name      string
	mongocCli *client
	// innerCollection *C.mongoc_collection_t
	// clientSession   *C.mongoc_client_session_t
	err error
}

func (c *collection) Name() string {
	return c.name
}
func (c *collection) Drop(ctx context.Context) error {
	return nil
}
func (c *collection) CreateIndex(index mongobyc.Index) error {
	return nil
}
func (c *collection) DropIndex(indexName string) error {
	return nil
}

func (c *collection) GetIndexes() (*mongobyc.GetIndexResult, error) {
	return nil, nil
}

func (c *collection) Count(ctx context.Context, filter interface{}) (uint64, error) {
	return 0, nil
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*mongobyc.DeleteResult, error) {
	return nil, nil
}

func (c *collection) DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*mongobyc.DeleteResult, error) {
	return nil, nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error {
	return nil
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error {
	return nil
}

func (c *collection) FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error {
	return nil
}

func (c *collection) InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error {
	return nil
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error {
	return nil
}

func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*mongobyc.UpdateResult, error) {
	return nil, nil
}

func (c *collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*mongobyc.UpdateResult, error) {
	return nil, nil
}

func (c *collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*mongobyc.ReplaceOneResult, error) {
	return nil, nil
}
