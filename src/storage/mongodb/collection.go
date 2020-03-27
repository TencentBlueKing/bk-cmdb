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

package mongodb

import (
	"context"

	"configcenter/src/storage/mongodb/options/aggregateopt"
	"configcenter/src/storage/mongodb/options/deleteopt"
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/mongodb/options/insertopt"
	"configcenter/src/storage/mongodb/options/replaceopt"
	"configcenter/src/storage/mongodb/options/updateopt"
)

// CollectionInterface collection operation methods
type CollectionInterface interface {
	Name() string
	Drop(ctx context.Context) error
	CreateIndex(index Index) error
	DropIndex(indexName string) error
	GetIndexes() (*QueryIndexResult, error)
	Count(ctx context.Context, filter interface{}) (uint64, error)

	DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*DeleteResult, error)

	Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error
	FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error

	Aggregate(ctx context.Context, pipeline interface{}, opts *aggregateopt.One, output interface{}) error

	InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error
	InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error

	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*UpdateResult, error)

	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*ReplaceOneResult, error)

	Update(ctx context.Context, filter interface{}, update map[string]interface{}, opts *updateopt.Many) (*UpdateResult, error)
}
