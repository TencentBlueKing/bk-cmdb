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

package dal

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/aggregateopt"
	"github.com/mongodb/mongo-go-driver/mongo/changestreamopt"
	"github.com/mongodb/mongo-go-driver/mongo/countopt"
	"github.com/mongodb/mongo-go-driver/mongo/deleteopt"
	"github.com/mongodb/mongo-go-driver/mongo/distinctopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
	"github.com/mongodb/mongo-go-driver/mongo/replaceopt"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type MongoCollectionClient interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...aggregateopt.Aggregate) (mongo.Cursor, error)
	Count(ctx context.Context, filter interface{}, opts ...countopt.Count) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...distinctopt.Distinct) ([]interface{}, error)
	Drop(ctx context.Context) error
	Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) *mongo.DocumentResult
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...findopt.DeleteOne) *mongo.DocumentResult
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...findopt.ReplaceOne) *mongo.DocumentResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...findopt.UpdateOne) *mongo.DocumentResult
	Indexes() mongo.IndexView
	InsertMany(ctx context.Context, documents []interface{}, opts ...insertopt.Many) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error)
	Name() string
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...replaceopt.Replace) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...updateopt.Update) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, options ...updateopt.Update) (*mongo.UpdateResult, error)
	Watch(ctx context.Context, pipeline interface{}, opts ...changestreamopt.ChangeStream) (mongo.Cursor, error)
}
