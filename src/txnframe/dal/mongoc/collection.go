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

// #cgo CFLAGS: -I/usr/local/include/libbson-1.0
// #cgo CFLAGS: -I/usr/local/include/libmongoc-1.0
// #cgo LDFLAGS: -lmongoc-1.0
// #cgo LDFLAGS: -lbson-1.0
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
	Count() (int64, error)
	DeleteMany() (*DeleteResult, error)
	DeleteOne() (*DeleteResult, error)
	Drop() error
	Find(ctx context.Context, filter interface{}, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, output interface{}) *DocumentResult
	InsertMany(ctx context.Context, document interface{}) (*InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}) (*InsertOneResult, error)
	UpdateMany() (*UpdateResult, error)
	UpdateOne() (*UpdateResult, error)
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
func (c *collection) Count() (int64, error) {
	return 0, nil
}
func (c *collection) DeleteMany() (*DeleteResult, error) {
	return nil, nil
}

func (c *collection) DeleteOne() (*DeleteResult, error) {
	return nil, nil
}
func (c *collection) Drop() error {
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

func (c *collection) InsertMany(ctx context.Context, document interface{}) (*InsertManyResult, error) {
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
func (c *collection) UpdateMany() (*UpdateResult, error) {
	return nil, nil
}
func (c *collection) UpdateOne() (*UpdateResult, error) {
	return nil, nil
}
