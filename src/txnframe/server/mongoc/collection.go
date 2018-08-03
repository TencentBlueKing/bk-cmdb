/*
 * Tencent is pleased to support the open source community by making čé˛¸ available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongoc

// #include <stdlib.h>
// #include "mongo.h"
import "C"

import (
	"context"
	"fmt"
	"unsafe"

	"configcenter/src/common/mapstr"
)

// CollectionInterface collection operation methods
type CollectionInterface interface {
	Name() string
	Count(ctx context.Context, filter interface{}) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}) (*DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*DeleteResult, error)
	Drop(ctx context.Context) error
	Find(ctx context.Context, filter interface{}, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, output interface{}) *DocumentResult
	InsertMany(ctx context.Context, document []interface{}) (*InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}) (*InsertOneResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
}

func newCollection(client *C.mongoc_client_t, dbName, collectionName string) CollectionInterface {
	coll := &collection{
		dbName:      dbName,
		name:        collectionName,
		innerClient: client,
	}

	innerDBName := C.CString(dbName)
	innerCollectionName := C.CString(collectionName)
	coll.innerCollection = C.mongoc_client_get_collection(client, innerDBName, innerCollectionName)

	C.free(unsafe.Pointer(innerDBName))
	C.free(unsafe.Pointer(innerCollectionName))

	return coll
}

type collection struct {
	name            string
	dbName          string
	innerClient     *C.mongoc_client_t
	innerCollection *C.mongoc_collection_t
}

func (c *collection) Name() string {
	return c.name
}
func (c *collection) Count(ctx context.Context, filter interface{}) (int64, error) {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return 0, err
	}
	defer C.bson_destroy(bsonFilter)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	count := C.mongoc_collection_count_documents(c.innerCollection, bsonFilter, nil, nil, &reply, &bsonErr)
	if 0 > count {
		return 0, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return 0, err
	}

	fmt.Println("delete result:", result)

	return int64(count), nil
}
func (c *collection) DeleteMany(ctx context.Context, filter interface{}) (*DeleteResult, error) {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_delete_many(c.innerCollection, bsonFilter, nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	fmt.Println("delete result:", result)

	return nil, nil
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}) (*DeleteResult, error) {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_delete_one(c.innerCollection, bsonFilter, nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	fmt.Println("delete result:", result)

	return nil, nil
}
func (c *collection) Drop(ctx context.Context) error {

	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_drop(c.innerCollection, &bsonErr)
	if !ok {
		return TransformError(bsonErr)
	}
	return nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, output interface{}) error {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return err
	}
	defer C.bson_destroy(bsonFilter)
	datas := []mapstr.MapStr{}
	cursor := C.mongoc_collection_find_with_opts(c.innerCollection, bsonFilter, nil, nil)
	for {
		var doc *C.bson_t
		if !C.mongoc_cursor_next(cursor, &doc) {
			break
		}
		docResult, err := mapstr.NewFromInterface(TransformBsonIntoGoString(doc))
		if nil != err {
			C.bson_destroy(doc)
			return err
		}

		datas = append(datas, docResult)
		C.bson_destroy(doc)
	}

	TransformMapStrIntoResult(datas, output)

	return nil
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, output interface{}) *DocumentResult {
	return nil
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}) (*InsertManyResult, error) {

	bsonDatas := make([]*C.bson_t, len(document))
	for idx, doc := range document {
		bsonData, err := TransformDocument(doc)
		if nil != err {
			return nil, err
		}
		defer C.bson_destroy(bsonData)
		bsonDatas[idx] = bsonData
	}

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_insert_many(c.innerCollection, (**C.bson_t)(unsafe.Pointer(&bsonDatas[0])), C.ulong(len(document)), nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	fmt.Println("insert manay:", result)

	return nil, nil
}
func (c *collection) InsertOne(ctx context.Context, document interface{}) (*InsertOneResult, error) {

	bsonData, err := TransformDocument(document)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonData)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_insert_one(c.innerCollection, bsonData, nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	cnt, err := result.Int64("insertedCount")
	if nil != err {
		return nil, fmt.Errorf("not found the insertedCount, %s", err.Error())
	}

	return &InsertOneResult{Count: cnt}, nil
}
func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error) {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	bsonUpdate, err := TransformDocument(update)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonUpdate)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_update_many(c.innerCollection, bsonFilter, bsonUpdate, nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	fmt.Println("update result:", result)

	return nil, nil
}
func (c *collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error) {

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	bsonUpdate, err := TransformDocument(update)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonUpdate)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_update_one(c.innerCollection, bsonFilter, bsonUpdate, nil, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	fmt.Println("update result:", result)
	return nil, nil
}
