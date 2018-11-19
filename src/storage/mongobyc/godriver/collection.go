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

	"configcenter/src/common/mapstr"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/deleteopt"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/mongobyc/insertopt"
	"configcenter/src/storage/mongobyc/replaceopt"
	"configcenter/src/storage/mongobyc/updateopt"

	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

var _ mongobyc.CollectionInterface = (*collection)(nil)

type collection struct {
	innerSession     mgo.Session
	innerCollection *mgo.Collection
	err             error
}

func newCollection(db *mgo.Database, collectionName string) mongobyc.CollectionInterface {

	return &collection{
		innerCollection: db.Collection(collectionName),
	}
}

func newCollectionWithSession(db *mgo.Database, innerSession mgo.Session, collectionName string)mongobyc.CollectionInterface{
	return &collection{
		innerSession: innerSession,
		innerCollection: db.Collection(collectionName),
	}
}

func (c *collection) Name() string {
	return c.innerCollection.Name()
}
func (c *collection) Drop(ctx context.Context) error {
	return c.innerCollection.Drop(context.TODO())
}
func (c *collection) CreateIndex(index mongobyc.Index) error {

	indexView := c.innerCollection.Indexes()

	keys := bsonx.Doc{}
	for indexKey := range index.Keys {
		keys.Append(indexKey, bsonx.Int32(-1))
	}

	options := bsonx.Doc{}
	options.Append("name", bsonx.String(index.Name))
	options.Append("unique", bsonx.Boolean(index.Unique))
	options.Append("background", bsonx.Boolean(index.Backgroupd))

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
	_, err := indexView.DropOne(context.TODO(), indexName)
	return err
}

func (c *collection) GetIndexes() (*mongobyc.GetIndexResult, error) {

	indexView := c.innerCollection.Indexes()
	cursor, err := indexView.List(context.TODO())

	if nil != err {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		elem := bsonx.Doc{}
		if err := cursor.Decode(elem); err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return nil, nil
}

func (c *collection) Count(ctx context.Context, filter interface{}) (uint64, error) {
	cnt, err := c.innerCollection.Count(ctx, filter)
	return uint64(cnt), err
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*mongobyc.DeleteResult, error) {

	delResult, err := c.innerCollection.DeleteOne(ctx, filter)
	if nil != err {
		return &mongobyc.DeleteResult{}, err
	}

	return &mongobyc.DeleteResult{DeletedCount: uint64(delResult.DeletedCount)}, nil
}

func (c *collection) DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*mongobyc.DeleteResult, error) {

	delResult, err := c.innerCollection.DeleteMany(ctx, filter)
	if nil != err {
		return &mongobyc.DeleteResult{DeletedCount: 0}, err
	}

	return &mongobyc.DeleteResult{DeletedCount: uint64(delResult.DeletedCount)}, nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error {

	switch tmp := filter.(type) {
	case string, []byte:
		condMap, err := mapstr.NewFromInterface(tmp)
		if nil != err {
			return err
		}
		filter = condMap
	}

	cursor, err := c.innerCollection.Find(ctx, filter)
	if nil != err {
		return err
	}

	datas := []mapstr.MapStr{}
	for cursor.Next(ctx) {
		result := mapstr.New()
		if err := cursor.Decode(&result); nil != err {
			return err
		}
		datas = append(datas, result)
	}

	mongobyc.TransformMapStrIntoResult(datas, output)
	return nil
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error {
	switch tmp := filter.(type) {
	case string, []byte:
		condMap, err := mapstr.NewFromInterface(tmp)
		if nil != err {
			return err
		}
		filter = condMap
	}
	return c.innerCollection.FindOne(ctx, filter).Decode(output)
}

func (c *collection) FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error {
	return c.innerCollection.FindOneAndUpdate(ctx, filter, update).Decode(nil)
}

func (c *collection) InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error {

	_, err := c.innerCollection.InsertOne(ctx, document)
	return err
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error {
	_, err := c.innerCollection.InsertMany(ctx, document)
	return err
}

func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*mongobyc.UpdateResult, error) {

	updateResult, err := c.innerCollection.UpdateMany(ctx, filter, update)
	if nil != err {
		return &mongobyc.UpdateResult{}, err
	}
	return &mongobyc.UpdateResult{
		MatchedCount:  uint64(updateResult.MatchedCount),
		ModifiedCount: uint64(updateResult.ModifiedCount),
	}, nil
}

func (c *collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*mongobyc.UpdateResult, error) {

	updateResult, err := c.innerCollection.UpdateOne(ctx, filter, update)
	if nil != err {
		return &mongobyc.UpdateResult{}, err
	}

	return &mongobyc.UpdateResult{
		MatchedCount:  uint64(updateResult.MatchedCount),
		ModifiedCount: uint64(updateResult.ModifiedCount),
	}, nil
}

func (c *collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*mongobyc.ReplaceOneResult, error) {

	replaceResult, err := c.innerCollection.ReplaceOne(ctx, filter, replacement)
	if nil != err {
		return &mongobyc.ReplaceOneResult{}, err
	}

	return &mongobyc.ReplaceOneResult{
		UpdateResult: mongobyc.UpdateResult{
			MatchedCount:  uint64(replaceResult.MatchedCount),
			ModifiedCount: uint64(replaceResult.ModifiedCount),
		},
	}, nil
}

func (c *collection) Close() error {
	return nil
}
