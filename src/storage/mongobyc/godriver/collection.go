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
	"log"

	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/deleteopt"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/mongobyc/insertopt"
	"configcenter/src/storage/mongobyc/replaceopt"
	"configcenter/src/storage/mongobyc/updateopt"

	"github.com/holmeswang/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

var _ mongobyc.CollectionInterface = (*collection)(nil)

type collection struct {
	name      string
	mongocCli *client
	// innerCollection *C.mongoc_collection_t
	// clientSession   *C.mongoc_client_session_t
	innerCollection *mgo.Collection

	err error
}

func (c *collection) Name() string {
	return c.innerCollection.Name()
}
func (c *collection) Drop(ctx context.Context) error {
	return c.innerCollection.Drop(context.TODO())
}
func (c *collection) CreateIndex(index mongobyc.Index) error {

	indexView := c.innerCollection.Indexes()

	keys := bson.NewDocument()
	for indexKey, indexValue := range index.Keys {
		keys.Append(bson.EC.Interface(indexKey, indexValue))
	}

	options := bson.NewDocument()
	options.Append(bson.EC.Interface("name", index.Name))
	options.Append(bson.EC.Interface("unique", index.Unique))
	options.Append(bson.EC.Interface("background", index.Backgroupd))

	_, err := indexView.CreateOne(
		context.TODO(),
		mgo.IndexModel{
			Keys:    keys,
			Options: options,
		},
	)

	return err
}
func (c *collection) DropIndex(indexName string) error {

	indexView := c.innerCollection.Indexes()
	_, err := indexView.DropIndex(context.TODO(), indexName)
	return err
}

func (c *collection) GetIndexes() (*mongobyc.GetIndexResult, error) {

	indexView := c.innerCollection.Indexes()
	cursor, err := indexView.List(context.TODO())

	if nil != err {
		return nil, err
	}

	ctx := context.Background()
	defer cur.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		elem := bson.NewDocument()
		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

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
