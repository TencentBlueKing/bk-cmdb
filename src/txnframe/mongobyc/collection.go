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
	"configcenter/src/txnframe/mongobyc/findopt"
)

// CollectionInterface collection operation methods
type CollectionInterface interface {
	Name() string
	Drop(ctx context.Context) error
	CreateIndex(index Index) error
	Count(ctx context.Context, filter interface{}) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}) (*DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts *findopt.Opts, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, opts *findopt.Opts, output interface{}) error
	FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.Opts, output interface{}) error
	InsertMany(ctx context.Context, document []interface{}) error
	InsertOne(ctx context.Context, document interface{}) error
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*UpdateResult, error)
}

func newCollectionWithoutSession(innerClient *client, collectionName string) CollectionInterface {
	coll := &collection{
		name:      collectionName,
		mongocCli: innerClient,
	}

	innerDBName := C.CString(innerClient.dbName)
	innerCollectionName := C.CString(collectionName)
	defer C.free(unsafe.Pointer(innerDBName))
	defer C.free(unsafe.Pointer(innerCollectionName))

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.db, innerCollectionName, nil, &bsonErr)
	if nil == coll.innerCollection {
		if 48 != bsonErr.code { /* code 48 is NamespaceExists, see error_codes.err in mongodb source */
			coll.err = TransformError(bsonErr)
			return coll
		}
		coll.innerCollection = C.mongoc_client_get_collection(innerClient.innerClient, innerDBName, innerCollectionName)
	}

	return coll
}

func newCollectionWithSession(innerClient *client, collectionName string, clientSession *C.mongoc_client_session_t) CollectionInterface {

	coll := &collection{
		name:          collectionName,
		mongocCli:     innerClient,
		clientSession: clientSession,
	}

	innerDBName := C.CString(innerClient.dbName)
	innerCollectionName := C.CString(collectionName)

	defer C.free(unsafe.Pointer(innerDBName))
	defer C.free(unsafe.Pointer(innerCollectionName))

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.db, innerCollectionName, nil, &bsonErr)
	if nil == coll.innerCollection {
		if 48 != bsonErr.code { /* code 48 is NamespaceExists, see error_codes.err in mongodb source */
			coll.err = TransformError(bsonErr)
			return coll
		}
		coll.innerCollection = C.mongoc_client_get_collection(innerClient.innerClient, innerDBName, innerCollectionName)
	}

	return coll
}

type collection struct {
	name            string
	mongocCli       *client
	innerCollection *C.mongoc_collection_t
	clientSession   *C.mongoc_client_session_t
	err             error
}

func (c *collection) Name() string {
	return c.name
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

func (c *collection) CreateIndex(index Index) error {

	data, err := TransformDocument(index)
	if nil != err {
		return err
	}
	defer C.bson_destroy(data)

	indexName := C.CString(index.Name)
	collName := C.CString(c.name)
	defer C.free(unsafe.Pointer(indexName))
	defer C.free(unsafe.Pointer(collName))
	var reply C.bson_t
	var bsonErr C.bson_error_t
	if !C.create_collection_index(c.mongocCli.db, collName, data, &reply, &bsonErr) {
		return TransformError(bsonErr)
	}
	C.bson_destroy(&reply)
	return nil
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

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Opts, output interface{}) error {

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

	if nil != opts {

		opts = TransformFindOpts(opts)

		if nil == operationOpts {
			operationOpts = C.bson_new()
		}

		limit := C.CString("limit")
		sort := C.CString("sort")
		skip := C.CString("skip")
		fields := C.CString("projection")

		C.bson_append_int64(operationOpts, limit, -1, C.longlong(opts.Limit))
		C.bson_append_int64(operationOpts, skip, -1, C.longlong(opts.Skip))
		var bsonSort, bsonFields *C.bson_t
		bsonFields, err := TransformDocument(opts.Fields)
		if nil != err {
			return err
		}

		if opts.Descending {
			bsonSort, err = TransformDocument(fmt.Sprintf(`{"%s":%d}`, opts.Sort, -1))
		} else {
			bsonSort, err = TransformDocument(fmt.Sprintf(`{"%s":%d}`, opts.Sort, 1))
		}
		if nil != err {
			return err
		}

		C.bson_append_document(operationOpts, sort, -1, bsonSort)
		C.bson_append_document(operationOpts, fields, -1, bsonFields)

		C.free(unsafe.Pointer(limit))
		C.free(unsafe.Pointer(sort))
		C.free(unsafe.Pointer(skip))
		C.free(unsafe.Pointer(fields))
		C.bson_destroy(bsonSort)

	}

	//fmt.Println("opts:", TransformBsonIntoGoString(operationOpts))

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

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.Opts, output interface{}) error {

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

	if nil != opts {

		opts = TransformFindOpts(opts)

		if nil == operationOpts {
			operationOpts = C.bson_new()
		}

		limit := C.CString("limit")
		sort := C.CString("sort")
		skip := C.CString("skip")
		fields := C.CString("projection")

		C.bson_append_int64(operationOpts, limit, -1, C.longlong(opts.Limit))
		C.bson_append_int64(operationOpts, skip, -1, C.longlong(opts.Skip))
		var bsonSort, bsonFields *C.bson_t
		bsonFields, err := TransformDocument(opts.Fields)
		if nil != err {
			return err
		}

		if opts.Descending {
			bsonSort, err = TransformDocument(fmt.Sprintf(`{"%s":%d}`, opts.Sort, -1))
		} else {
			bsonSort, err = TransformDocument(fmt.Sprintf(`{"%s":%d}`, opts.Sort, 1))
		}
		if nil != err {
			return err
		}

		C.bson_append_document(operationOpts, sort, -1, bsonSort)
		C.bson_append_document(operationOpts, fields, -1, bsonFields)

		C.free(unsafe.Pointer(limit))
		C.free(unsafe.Pointer(sort))
		C.free(unsafe.Pointer(skip))
		C.free(unsafe.Pointer(fields))
		C.bson_destroy(bsonSort)

	}

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
func (c *collection) FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.Opts, output interface{}) error {
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
	ok := C.mongoc_collection_insert_many(c.innerCollection, (**C.bson_t)(unsafe.Pointer(&bsonDatas[0])), C.size_t(len(document)), operationOpts, &reply, &bsonErr)
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
