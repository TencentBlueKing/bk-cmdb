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
	"configcenter/src/storage/mongodb/options/aggregateopt"
	"configcenter/src/storage/mongodb/options/deleteopt"
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/mongodb/options/insertopt"
	"configcenter/src/storage/mongodb/options/replaceopt"
	"configcenter/src/storage/mongodb/options/updateopt"

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

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			return c.innerCollection.Drop(mctx)
		})
	}

	// no session
	return c.innerCollection.Drop(ctx)
}
func (c *collection) CreateIndex(index mongodb.Index) error {

	indexView := c.innerCollection.Indexes()

	keys := bsonx.Doc{}
	for key, val := range index.Keys {
		keys = keys.Append(key, bsonx.Int32(val))
	}

	indexOpts := &options.IndexOptions{
		Name:       &index.Name,
		Background: &index.Background,
		Unique:     &index.Unique,
	}

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(context.TODO(), c.innerSession, func(mctx mongo.SessionContext) error {
			_, err := indexView.CreateOne(
				mctx,
				mongo.IndexModel{
					Keys:    keys,
					Options: indexOpts,
				},
			)
			return err
		})
	}

	// no session
	_, err := indexView.CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    keys,
			Options: indexOpts,
		},
	)

	return err

}
func (c *collection) DropIndex(indexName string) error {

	indexView := c.innerCollection.Indexes()

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(context.TODO(), c.innerSession, func(mctx mongo.SessionContext) error {
			_, err := indexView.DropOne(mctx, indexName)
			return err
		})
	}

	// no session
	_, err := indexView.DropOne(context.TODO(), indexName)
	return err
}

func (c *collection) GetIndexes() (*mongodb.QueryIndexResult, error) {

	indexView := c.innerCollection.Indexes()

	// in a session
	if nil != c.innerSession {
		returnIndexResult := &mongodb.QueryIndexResult{}
		err := mongo.WithSession(context.TODO(), c.innerSession, func(mctx mongo.SessionContext) error {

			cursor, err := indexView.List(mctx)
			if nil != err {
				return err
			}

			defer cursor.Close(mctx)
			returnIndexResult, err = c.parseIndexResult(mctx, cursor)
			return err
		})

		if nil != err {
			return returnIndexResult, err
		}

		return returnIndexResult, nil
	}

	// no session
	cursor, err := indexView.List(context.TODO())
	if nil != err {
		return &mongodb.QueryIndexResult{}, err
	}

	defer cursor.Close(context.TODO())
	return c.parseIndexResult(context.TODO(), cursor)
}

func (c *collection) Count(ctx context.Context, filter interface{}) (uint64, error) {

	// in a session
	if nil != c.innerSession {
		var innerCnt uint64
		err := mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			cnt, err := c.innerCollection.Count(ctx, filter)
			innerCnt = uint64(cnt)
			return err
		})
		return innerCnt, err
	}

	// no session
	cnt, err := c.innerCollection.Count(ctx, filter)
	return uint64(cnt), err
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*mongodb.DeleteResult, error) {

	deleteOption := &options.DeleteOptions{}
	if nil != opts {
		deleteOption = opts.ConvertToMongoOptions()
	}

	returnResult := &mongodb.DeleteResult{}

	// in a session
	if nil != c.innerSession {
		err := mongo.WithSession(ctx, c.innerSession, func(mtcx mongo.SessionContext) error {
			delResult, err := c.innerCollection.DeleteOne(mtcx, filter, deleteOption)
			if nil != err {
				return err
			}
			returnResult.DeletedCount = uint64(delResult.DeletedCount)
			return nil
		})
		return returnResult, err
	}

	// no session
	delResult, err := c.innerCollection.DeleteOne(ctx, filter, deleteOption)
	if nil != err {
		return returnResult, err
	}

	returnResult.DeletedCount = uint64(delResult.DeletedCount)
	return returnResult, nil
}

func (c *collection) DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*mongodb.DeleteResult, error) {

	deleteOption := &options.DeleteOptions{}
	if nil != opts {
		deleteOption = opts.ConvertToMongoOptions()
	}

	returnResult := &mongodb.DeleteResult{}

	// in a session
	if nil != c.innerSession {
		err := mongo.WithSession(ctx, c.innerSession, func(mtcx mongo.SessionContext) error {
			delResult, err := c.innerCollection.DeleteMany(ctx, filter, deleteOption)
			if nil != err {
				return err
			}
			returnResult.DeletedCount = uint64(delResult.DeletedCount)
			return nil
		})
		return returnResult, err
	}

	// no session
	delResult, err := c.innerCollection.DeleteMany(ctx, filter, deleteOption)
	if nil != err {
		return &mongodb.DeleteResult{DeletedCount: 0}, err
	}
	returnResult.DeletedCount = uint64(delResult.DeletedCount)
	return returnResult, nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error {

	findOptions := &options.FindOptions{}
	if nil != opts {
		findOptions = opts.ConvertToMongoOptions()
	}

	datas := []mapstr.MapStr{}

	// in a session
	if nil != c.innerSession {

		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {

			cursor, err := c.innerCollection.Find(mctx, filter, findOptions)

			if nil != err {
				return err
			}

			defer cursor.Close(mctx)
			datas := []mapstr.MapStr{}
			for cursor.Next(mctx) {
				result := mapstr.New()
				if err := cursor.Decode(&result); nil != err {
					return err
				}
				datas = append(datas, result)
			}

			if err := cursor.Err(); err != nil {
				return err
			}

			return mapstr.ConvertArrayMapStrInto(datas, output)
		})
	}

	// no session
	cursor, err := c.innerCollection.Find(ctx, filter, findOptions)
	if nil != err {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		result := mapstr.New()
		if err := cursor.Decode(&result); nil != err {
			return err
		}
		datas = append(datas, result)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	// package result
	return mapstr.ConvertArrayMapStrInto(datas, output)
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error {

	findOptions := &options.FindOneOptions{}
	if nil != opts {
		findOptions = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			return c.innerCollection.FindOne(mctx, filter, findOptions).Decode(output)
		})
	}

	// no session
	return c.innerCollection.FindOne(ctx, filter, findOptions).Decode(output)
}

func (c *collection) FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error {

	findOneAndModify := &options.FindOneAndUpdateOptions{}
	if nil != opts {
		findOneAndModify = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			return c.innerCollection.FindOneAndUpdate(mctx, filter, update, findOneAndModify).Decode(output)
		})
	}

	// no session
	return c.innerCollection.FindOneAndUpdate(ctx, filter, update, findOneAndModify).Decode(output)
}

func (c *collection) AggregateOne(ctx context.Context, pipeline interface{}, opts *aggregateopt.One, output interface{}) error {
	aggregateOptions := &options.AggregateOptions{}
	if nil != opts {
		aggregateOptions = opts.ConvertToMongoOptions()
	}

	mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
		cursor, err := c.innerCollection.Aggregate(mctx, pipeline, aggregateOptions)
		if nil != err {
			return err
		}

		defer cursor.Close(mctx)
		for cursor.Next(ctx) {
			return cursor.Decode(output)
		}
		return cursor.Err()
	})

	// no session
	cursor, err := c.innerCollection.Aggregate(ctx, pipeline, aggregateOptions)
	if nil != err {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		return cursor.Decode(output)
	}
	return cursor.Err()
}

func (c *collection) InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error {

	insertOption := &options.InsertOneOptions{}
	if nil != opts {
		insertOption = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			_, err := c.innerCollection.InsertOne(mctx, document, insertOption)
			return err
		})
	}

	// no session
	_, err := c.innerCollection.InsertOne(ctx, document, insertOption)
	return err
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error {

	insertOption := &options.InsertManyOptions{}
	if nil != opts {
		insertOption = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {
		return mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			_, err := c.innerCollection.InsertMany(mctx, document, insertOption)
			return err
		})
	}

	// no session
	_, err := c.innerCollection.InsertMany(ctx, document, insertOption)
	return err
}

func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*mongodb.UpdateResult, error) {

	updateOption := &options.UpdateOptions{}
	if nil != opts {
		updateOption = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {
		returnResult := &mongodb.UpdateResult{}
		err := mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			updateResult, err := c.innerCollection.UpdateMany(mctx, filter, bson.M{"$set": update}, updateOption)
			if nil != err {
				return err
			}
			returnResult = &mongodb.UpdateResult{
				MatchedCount:  uint64(updateResult.MatchedCount),
				ModifiedCount: uint64(updateResult.ModifiedCount),
			}

			return nil
		})

		return returnResult, err

	}

	// no session
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

	// in a session
	if nil != c.innerSession {

		returnResult := &mongodb.UpdateResult{}
		err := mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {
			updateResult, err := c.innerCollection.UpdateOne(mctx, filter, bson.M{"$set": update}, updateOption)
			if nil != err {
				return err
			}

			returnResult = &mongodb.UpdateResult{
				MatchedCount:  uint64(updateResult.MatchedCount),
				ModifiedCount: uint64(updateResult.ModifiedCount),
			}

			return nil
		})

		return returnResult, err
	}

	// no session
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

	replaceOption := &options.ReplaceOptions{}
	if nil != opts {
		replaceOption = opts.ConvertToMongoOptions()
	}

	// in a session
	if nil != c.innerSession {

		returnResult := &mongodb.ReplaceOneResult{}
		err := mongo.WithSession(ctx, c.innerSession, func(mctx mongo.SessionContext) error {

			replaceResult, err := c.innerCollection.ReplaceOne(mctx, filter, replacement, replaceOption)
			if nil != err {
				return err
			}

			returnResult = &mongodb.ReplaceOneResult{
				UpdateResult: mongodb.UpdateResult{
					MatchedCount:  uint64(replaceResult.MatchedCount),
					ModifiedCount: uint64(replaceResult.ModifiedCount),
				},
			}
			return nil
		})

		return returnResult, err
	}

	// no session
	replaceResult, err := c.innerCollection.ReplaceOne(ctx, filter, replacement, replaceOption)
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
