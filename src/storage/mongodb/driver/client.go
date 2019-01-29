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

package driver

import (
	"context"

	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/options/aggregateopt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
)

var _ mongodb.CommonClient = (*client)(nil)

type collectionName string

type client struct {
	createdByPool  bool
	uri            string
	collectionMaps map[collectionName]mongodb.CollectionInterface
	innerDB        *database
	innerClient    *mongo.Client
}

// NewClient create a mongoc client instance
func NewClient(uri string) mongodb.CommonClient {
	return &client{
		uri:            uri,
		collectionMaps: map[collectionName]mongodb.CollectionInterface{},
	}
}

func (c *client) Open() error {

	cnnstr, err := connstring.Parse(c.uri)
	if nil != err {
		return err
	}

	c.innerClient, err = mongo.NewClient(c.uri)
	if nil != err {
		return err
	}

	if err := c.innerClient.Connect(context.TODO()); nil != err {
		return err
	}

	c.innerDB = newDatabase(c.innerClient.Database(cnnstr.Database))

	return nil
}

func (c *client) Close() error {

	if nil != c.innerClient {
		return nil
	}

	return c.innerClient.Disconnect(context.TODO())
}

func (c *client) Ping() error {

	return c.innerClient.Ping(context.TODO(), nil)
}

func (c *client) Database() mongodb.Database {
	return c.innerDB
}

func (c *client) Collection(collName string) mongodb.CollectionInterface {
	target, ok := c.collectionMaps[collectionName(collName)]
	if !ok {
		target = newCollection(c.innerDB.innerDatabase, collName)
		c.collectionMaps[collectionName(collName)] = target
	}
	return target
}

func (c *client) AggregateOne(ctx context.Context, pipeline interface{}, opts *aggregateopt.One, output interface{}) error {
	aggregateOptions := &options.AggregateOptions{}
	if nil != opts {
		aggregateOptions = opts.ConvertToMongoOptions()
	}

	c.innerClient.Database("").Collection("asdf").Aggregate
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

func (c *client) Session() mongodb.SessionOperation {
	return newSessionOperation(c)
}
