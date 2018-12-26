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

	"configcenter/src/common/mapstr"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/deleteopt"
	"configcenter/src/storage/mongodb/findopt"
	"configcenter/src/storage/mongodb/insertopt"
	"configcenter/src/storage/mongodb/replaceopt"
	"configcenter/src/storage/mongodb/updateopt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

var _ mongodb.CollectionInterface = (*collection)(nil)

type collection struct {
	innerSession    mongo.Session
	innerCollection *mongo.Collection
	err             error
}

func newCollection(db *mongo.Database, collectionName string) mongodb.CollectionInterface {

	return &collection{
		innerCollection: db.Collection(collectionName),
	}
}

func newCollectionWithSession(db *mongo.Database, innerSession mongo.Session, collectionName string) mongodb.CollectionInterface {
	return &collection{
		innerSession:    innerSession,
		innerCollection: db.Collection(collectionName),
	}
}

func (c *collection) Name() string {
	return c.innerCollection.Name()
}
func (c *collection) Drop(ctx context.Context) error {
	return c.innerCollection.Drop(ctx)
}
func (c *collection) CreateIndex(index mongodb.Index) error {

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
		mongo.IndexModel{
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

func (c *collection) GetIndexes() (*mongodb.GetIndexResult, error) {

	indexView := c.innerCollection.Indexes()
	cursor, err := indexView.List(context.TODO())

	if nil != err {
		return &mongodb.GetIndexResult{}, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		elem := bsonx.Doc{}
		if err := cursor.Decode(elem); err != nil {
			return &mongodb.GetIndexResult{}, err
		}
	}

	if err := cursor.Err(); err != nil {
		return &mongodb.GetIndexResult{}, err
	}

	return nil, nil
}

func (c *collection) Count(ctx context.Context, filter interface{}) (uint64, error) {
	cnt, err := c.innerCollection.Count(ctx, filter)
	return uint64(cnt), err
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*mongodb.DeleteResult, error) {

	deleteOption := &options.DeleteOptions{}
	if nil != opts {
		deleteOption = opts.ConvertToMongoOptions()
	}
	delResult, err := c.innerCollection.DeleteOne(ctx, filter, deleteOption)
	if nil != err {
		return &mongodb.DeleteResult{}, err
	}

	return &mongodb.DeleteResult{DeletedCount: uint64(delResult.DeletedCount)}, nil
}

func (c *collection) DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*mongodb.DeleteResult, error) {

	deleteOption := &options.DeleteOptions{}
	if nil != opts {
		deleteOption = opts.ConvertToMongoOptions()
	}

	delResult, err := c.innerCollection.DeleteMany(ctx, filter, deleteOption)
	if nil != err {
		return &mongodb.DeleteResult{DeletedCount: 0}, err
	}

	return &mongodb.DeleteResult{DeletedCount: uint64(delResult.DeletedCount)}, nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error {

	findOptions := &options.FindOptions{}
	if nil != opts {
		findOptions = opts.ConvertToMongoOptions()
	}

	cursor, err := c.innerCollection.Find(ctx, filter, findOptions)
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

	mongodb.TransformMapStrIntoResult(datas, output)
	return nil
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error {

	findOptions := &options.FindOneOptions{}
	if nil != opts {
		findOptions = opts.ConvertToMongoOptions()
	}

	return c.innerCollection.FindOne(ctx, filter, findOptions).Decode(output)
}

func (c *collection) FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error {

	findOneAndModify := &options.FindOneAndUpdateOptions{}
	if nil != opts {
		findOneAndModify = opts.ConvertToMongoOptions()
	}
	return c.innerCollection.FindOneAndUpdate(ctx, filter, update, findOneAndModify).Decode(output)
}

func (c *collection) InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error {

	insertOption := &options.InsertOneOptions{}
	if nil != opts {
		insertOption = opts.ConvertToMongoOptions()
	}
	_, err := c.innerCollection.InsertOne(ctx, document, insertOption)
	return err
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error {

	insertOption := &options.InsertManyOptions{}
	if nil != opts {
		insertOption = opts.ConvertToMongoOptions()
	}
	_, err := c.innerCollection.InsertMany(ctx, document, insertOption)
	return err
}

func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*mongodb.UpdateResult, error) {

	updateOption := &options.UpdateOptions{}
	if nil != opts {
		updateOption = opts.ConvertToMongoOptions()
	}

	updateResult, err := c.innerCollection.UpdateMany(ctx, filter, bson.M{"$set": update}, updateOption)
	if nil != err {
		return &mongodb.UpdateResult{}, err
	}
	return &mongodb.UpdateResult{
		MatchedCount:  uint64(updateResult.MatchedCount),
		ModifiedCount: uint64(updateResult.ModifiedCount),
	}, nil
}

func (c *collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*mongodb.UpdateResult, error) {

	updateOption := &options.UpdateOptions{}
	if nil != opts {
		updateOption = opts.ConvertToMongoOptions()
	}

	updateResult, err := c.innerCollection.UpdateOne(ctx, filter, bson.M{"$set": update}, updateOption)
	if nil != err {
		return &mongodb.UpdateResult{}, err
	}

	return &mongodb.UpdateResult{
		MatchedCount:  uint64(updateResult.MatchedCount),
		ModifiedCount: uint64(updateResult.ModifiedCount),
	}, nil
}

func (c *collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*mongodb.ReplaceOneResult, error) {

	replaceResult, err := c.innerCollection.ReplaceOne(ctx, filter, replacement)
	if nil != err {
		return &mongodb.ReplaceOneResult{}, err
	}

	return &mongodb.ReplaceOneResult{
		UpdateResult: mongodb.UpdateResult{
			MatchedCount:  uint64(replaceResult.MatchedCount),
			ModifiedCount: uint64(replaceResult.ModifiedCount),
		},
	}, nil
}

func (c *collection) Close() error {
	return nil
}
