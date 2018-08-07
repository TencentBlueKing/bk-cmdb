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

package mongobyc

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
	FindOne(ctx context.Context, filter interface{}, output interface{}) error
	InsertMany(ctx context.Context, document []interface{}) error
	InsertOne(ctx context.Context, document interface{}) error
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
}

func newCollectionWithoutSession(innerClient *client, collectionName string) CollectionInterface {
	coll := &collection{
		dbName:      innerClient.dbName,
		name:        collectionName,
		innerClient: innerClient.innerClient,
	}

	innerDBName := C.CString(innerClient.dbName)
	innerCollectionName := C.CString(collectionName)

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.db, innerCollectionName, nil, &bsonErr)
	if nil == coll.innerCollection {
		if 48 != bsonErr.code { /* code 48 is NamespaceExists, see error_codes.err in mongodb source */
			coll.err = TransformError(bsonErr)
			return coll
		}
		coll.innerCollection = C.mongoc_client_get_collection(innerClient.innerClient, innerDBName, innerCollectionName)
	}

	C.free(unsafe.Pointer(innerDBName))
	C.free(unsafe.Pointer(innerCollectionName))

	return coll
}

func newCollectionWithSession(innerClient *client, collectionName string, clientSession *C.mongoc_client_session_t) CollectionInterface {

	coll := &collection{
		dbName:        innerClient.dbName,
		name:          collectionName,
		innerClient:   innerClient.innerClient,
		clientSession: clientSession,
	}

	innerDBName := C.CString(innerClient.dbName)
	innerCollectionName := C.CString(collectionName)

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.db, innerCollectionName, nil, &bsonErr)
	if nil == coll.innerCollection {
		if 48 != bsonErr.code { /* code 48 is NamespaceExists, see error_codes.err in mongodb source */
			coll.err = TransformError(bsonErr)
			return coll
		}
		coll.innerCollection = C.mongoc_client_get_collection(innerClient.innerClient, innerDBName, innerCollectionName)
	}

	C.free(unsafe.Pointer(innerDBName))
	C.free(unsafe.Pointer(innerCollectionName))

	return coll
}

type collection struct {
	name            string
	dbName          string
	innerClient     *C.mongoc_client_t
	innerCollection *C.mongoc_collection_t
	clientSession   *C.mongoc_client_session_t
	err             error
}

func (c *collection) Name() string {
	return c.name
}

func (c *collection) getOperationOpts() (*C.bson_t, error) {

	if nil != c.err {
		return nil, c.err
	}

	var bsonErr C.bson_error_t
	if nil != c.clientSession {
		operationOpts := C.bson_new()
		if !C.mongoc_client_session_append(c.clientSession, operationOpts, &bsonErr) {
			return nil, TransformError(bsonErr)
		}
		return operationOpts, nil
	}
	return nil, nil
}

func (c *collection) Count(ctx context.Context, filter interface{}) (int64, error) {

	if nil != c.err {
		return 0, c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return 0, err
	}
	defer C.bson_destroy(bsonFilter)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return 0, err
	}
	defer C.bson_destroy(operationOpts)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	count := C.mongoc_collection_count_documents(c.innerCollection, bsonFilter, operationOpts, nil, &reply, &bsonErr)
	if 0 > count {
		return 0, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	return int64(count), nil
}
func (c *collection) DeleteMany(ctx context.Context, filter interface{}) (*DeleteResult, error) {

	if nil != c.err {
		return nil, c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_delete_many(c.innerCollection, bsonFilter, operationOpts, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	return nil, nil
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}) (*DeleteResult, error) {

	if nil != c.err {
		return nil, c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_delete_one(c.innerCollection, bsonFilter, operationOpts, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	return nil, nil
}
func (c *collection) Drop(ctx context.Context) error {

	if nil != c.err {
		return c.err
	}

	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_drop(c.innerCollection, &bsonErr)
	if !ok {
		return TransformError(bsonErr)
	}
	return nil
}

func (c *collection) Find(ctx context.Context, filter interface{}, output interface{}) error {

	if nil != c.err {
		return c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return err
	}

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return err
	}
	defer C.bson_destroy(operationOpts)

	defer C.bson_destroy(bsonFilter)
	datas := []mapstr.MapStr{}
	cursor := C.mongoc_collection_find_with_opts(c.innerCollection, bsonFilter, operationOpts, nil)
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

func (c *collection) FindOne(ctx context.Context, filter interface{}, output interface{}) error {

	if nil != c.err {
		return c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return err
	}

	defer C.bson_destroy(bsonFilter)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return err
	}
	defer C.bson_destroy(operationOpts)

	cursor := C.mongoc_collection_find_with_opts(c.innerCollection, bsonFilter, operationOpts, nil)
	for {
		var doc *C.bson_t
		if !C.mongoc_cursor_next(cursor, &doc) {
			break
		}
		defer C.bson_destroy(doc)
		docResult, err := mapstr.NewFromInterface(TransformBsonIntoGoString(doc))
		if nil != err {
			return err
		}

		return docResult.MarshalJSONInto(output) // only return one document
	}

	return nil
}

func (c *collection) InsertMany(ctx context.Context, document []interface{}) error {

	if nil != c.err {
		return c.err
	}

	bsonDatas := make([]*C.bson_t, len(document))
	for idx, doc := range document {
		bsonData, err := TransformDocument(doc)
		if nil != err {
			return err
		}
		defer C.bson_destroy(bsonData)
		bsonDatas[idx] = bsonData
	}

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_insert_many(c.innerCollection, (**C.bson_t)(unsafe.Pointer(&bsonDatas[0])), C.ulong(len(document)), operationOpts, &reply, &bsonErr)
	if !ok {
		return TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	return nil
}
func (c *collection) InsertOne(ctx context.Context, document interface{}) error {

	if nil != c.err {
		return c.err
	}

	bsonData, err := TransformDocument(document)
	if nil != err {
		return err
	}
	defer C.bson_destroy(bsonData)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_insert_one(c.innerCollection, bsonData, operationOpts, &reply, &bsonErr)
	if !ok {
		return TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	return nil
}
func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error) {

	if nil != c.err {
		return nil, c.err
	}

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

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_update_many(c.innerCollection, bsonFilter, bsonUpdate, operationOpts, &reply, &bsonErr)
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

	if nil != c.err {
		return nil, c.err
	}

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

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_update_one(c.innerCollection, bsonFilter, bsonUpdate, operationOpts, &reply, &bsonErr)
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
