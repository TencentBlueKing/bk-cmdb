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
	"configcenter/src/storage/mongobyc/deleteopt"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/mongobyc/insertopt"
	"configcenter/src/storage/mongobyc/replaceopt"
	"configcenter/src/storage/mongobyc/updateopt"
)

// CollectionInterface collection operation methods
type CollectionInterface interface {
	Name() string
	Drop(ctx context.Context) error
	CreateIndex(index Index) error
	DropIndex(indexName string) error
	GetIndexes() (*GetIndexResult, error)
	Count(ctx context.Context, filter interface{}) (uint64, error)

	DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*DeleteResult, error)

	Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error
	FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error

	InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error
	InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error

	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*UpdateResult, error)

	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*ReplaceOneResult, error)
}

func newCollectionWithoutSession(innerClient *client, collectionName string) CollectionInterface {
	coll := &collection{
		name:      collectionName,
		mongocCli: innerClient,
	}

	innerDBName := C.CString(innerClient.innerDB.dbName)
	innerCollectionName := C.CString(collectionName)
	defer C.free(unsafe.Pointer(innerDBName))
	defer C.free(unsafe.Pointer(innerCollectionName))

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.innerDB.db, innerCollectionName, nil, &bsonErr)
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

	innerDBName := C.CString(innerClient.innerDB.dbName)
	innerCollectionName := C.CString(collectionName)

	defer C.free(unsafe.Pointer(innerDBName))
	defer C.free(unsafe.Pointer(innerCollectionName))

	var bsonErr C.bson_error_t
	coll.innerCollection = C.mongoc_database_create_collection(innerClient.innerDB.db, innerCollectionName, nil, &bsonErr)
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
	defer C.bson_destroy(&reply)
	if !C.create_collection_index(c.mongocCli.innerDB.db, collName, data, &reply, &bsonErr) {
		return TransformError(bsonErr)
	}

	//fmt.Println("reply:", TransformBsonIntoGoString(&reply))

	return nil
}

func (c *collection) DropIndex(indexName string) error {
	cstrIndexName := C.CString(indexName)
	defer C.free(unsafe.Pointer(cstrIndexName))
	var bsonErr C.bson_error_t
	if !C.mongoc_collection_drop_index(c.innerCollection, cstrIndexName, &bsonErr) {
		return TransformError(bsonErr)
	}
	return nil
}

func (c *collection) GetIndexes() (*GetIndexResult, error) {

	indexRst := &GetIndexResult{}
	cursor := C.mongoc_collection_find_indexes_with_opts(c.innerCollection, nil)
	for {
		var doc *C.bson_t
		if !C.mongoc_cursor_next(cursor, &doc) {
			break
		}

		docResult, err := mapstr.NewFromInterface(TransformBsonIntoGoString(doc))
		if nil != err {
			C.bson_destroy(doc)
			return nil, err
		}
		C.bson_destroy(doc)
		rst := IndexResult{}
		if err := docResult.MarshalJSONInto(&rst); nil != err {
			return nil, err
		}

		indexRst.Indexes = append(indexRst.Indexes, rst)
	}

	return indexRst, nil

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

func (c *collection) Count(ctx context.Context, filter interface{}) (uint64, error) {

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

	return uint64(count), nil
}
func (c *collection) DeleteMany(ctx context.Context, filter interface{}, opts *deleteopt.Many) (*DeleteResult, error) {

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

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	deleteRst := &DeleteResult{}
	err = result.MarshalJSONInto(deleteRst)
	return deleteRst, err
}

func (c *collection) DeleteOne(ctx context.Context, filter interface{}, opts *deleteopt.One) (*DeleteResult, error) {

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

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	deleteRst := &DeleteResult{}
	err = result.MarshalJSONInto(deleteRst)
	return deleteRst, err
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts *findopt.Many, output interface{}) error {

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

		opts.Opts = *TransformFindOpts(&opts.Opts)

		if nil == operationOpts {
			operationOpts = C.bson_new()
		}

		limit := C.CString("limit")
		sort := C.CString("sort")
		skip := C.CString("skip")

		C.bson_append_int64(operationOpts, limit, -1, C.int64_t(opts.Limit))
		C.bson_append_int64(operationOpts, skip, -1, C.int64_t(opts.Skip))
		var bsonSort *C.bson_t

		if 0 != len(opts.Fields) {
			fields := C.CString("projection")
			defer C.free(unsafe.Pointer(fields))
			var bsonFields *C.bson_t
			bsonFields, err := TransformDocument(opts.Fields)
			if nil != err {
				return err
			}

			C.bson_append_document(operationOpts, fields, -1, bsonFields)
			C.bson_destroy(bsonFields)
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
		C.free(unsafe.Pointer(limit))
		C.free(unsafe.Pointer(sort))
		C.free(unsafe.Pointer(skip))
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
	//fmt.Println("datas:", datas)
	TransformMapStrIntoResult(datas, output)

	return nil
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts *findopt.One, output interface{}) error {

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

		opts.Opts = *TransformFindOpts(&opts.Opts)

		if nil == operationOpts {
			operationOpts = C.bson_new()
		}

		limit := C.CString("limit")
		sort := C.CString("sort")
		skip := C.CString("skip")

		C.bson_append_int64(operationOpts, limit, -1, C.int64_t(opts.Limit))
		C.bson_append_int64(operationOpts, skip, -1, C.int64_t(opts.Skip))
		var bsonSort *C.bson_t

		if 0 != len(opts.Fields) {
			fields := C.CString("projection")
			defer C.free(unsafe.Pointer(fields))
			var bsonFields *C.bson_t
			bsonFields, err := TransformDocument(opts.Fields)
			if nil != err {
				return err
			}
			C.bson_append_document(operationOpts, fields, -1, bsonFields)
			C.bson_destroy(bsonFields)
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

		C.free(unsafe.Pointer(limit))
		C.free(unsafe.Pointer(sort))
		C.free(unsafe.Pointer(skip))

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
func (c *collection) FindAndModify(ctx context.Context, filter interface{}, update interface{}, opts *findopt.FindAndModify, output interface{}) error {

	if nil != c.err {
		return c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return err
	}

	defer C.bson_destroy(bsonFilter)

	bsonUpdate, err := TransformDocument(update)
	if nil != err {
		return err
	}

	defer C.bson_destroy(bsonUpdate)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return err
	}
	defer C.bson_destroy(operationOpts)

	findModifyOpts := C.mongoc_find_and_modify_opts_new()
	defer C.mongoc_find_and_modify_opts_destroy(findModifyOpts)

	if nil != opts {
		opts.Opts = *TransformFindOpts(&opts.Opts)
		if nil == operationOpts {
			operationOpts = C.bson_new()
		}

		limit := C.CString("limit")
		//sort := C.CString("sort")
		skip := C.CString("skip")
		//fields := C.CString("projection")

		C.bson_append_int64(operationOpts, limit, -1, C.int64_t(opts.Limit))
		C.bson_append_int64(operationOpts, skip, -1, C.int64_t(opts.Skip))
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
		/*
			C.bson_append_document(operationOpts, sort, -1, bsonSort)
			C.bson_append_document(operationOpts, fields, -1, bsonFields)

			C.free(unsafe.Pointer(limit))
			C.free(unsafe.Pointer(sort))
			C.free(unsafe.Pointer(skip))
			C.free(unsafe.Pointer(fields))
			C.bson_destroy(bsonSort)
		*/
		C.mongoc_find_and_modify_opts_set_sort(findModifyOpts, bsonSort)
		C.mongoc_find_and_modify_opts_set_fields(findModifyOpts, bsonFields)
	}

	C.mongoc_find_and_modify_opts_set_update(findModifyOpts, bsonUpdate)
	flag := C.mongoc_find_and_modify_opts_get_flags(findModifyOpts)
	if opts.Upsert {
		flag = flag | C.MONGOC_FIND_AND_MODIFY_UPSERT
	}
	if opts.Remove {
		flag = flag | C.MONGOC_FIND_AND_MODIFY_REMOVE
	}
	if opts.New {
		flag = flag | C.MONGOC_FIND_AND_MODIFY_RETURN_NEW
	}
	C.mongoc_find_and_modify_opts_set_flags(findModifyOpts, flag)

	var reply C.bson_t
	var bsonErr C.bson_error_t
	ok := C.mongoc_collection_find_and_modify_with_opts(c.innerCollection, bsonFilter, findModifyOpts, &reply, &bsonErr)
	defer C.bson_destroy(&reply)
	if !ok {
		return TransformError(bsonErr)
	}

	docResult, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return err
	}

	valStr, err := docResult.MapStr("value")
	if nil != err {
		return err
	}

	return valStr.MarshalJSONInto(output)
}
func (c *collection) InsertMany(ctx context.Context, document []interface{}, opts *insertopt.Many) error {

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
func (c *collection) InsertOne(ctx context.Context, document interface{}, opts *insertopt.One) error {

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

func (c *collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.Many) (*UpdateResult, error) {

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

	updateRst := &UpdateResult{}
	err = result.MarshalJSONInto(updateRst)
	return updateRst, err

}
func (c *collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *updateopt.One) (*UpdateResult, error) {

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

	updateRst := &UpdateResult{}
	err = result.MarshalJSONInto(updateRst)
	return updateRst, err
}

func (c *collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts *replaceopt.One) (*ReplaceOneResult, error) {

	if nil != c.err {
		return nil, c.err
	}

	bsonFilter, err := TransformDocument(filter)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonFilter)

	bsonReplacement, err := TransformDocument(replacement)
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(bsonReplacement)

	operationOpts, err := c.getOperationOpts()
	if nil != err {
		return nil, err
	}
	defer C.bson_destroy(operationOpts)

	var bsonErr C.bson_error_t
	var reply C.bson_t
	ok := C.mongoc_collection_replace_one(c.innerCollection, bsonFilter, bsonReplacement, operationOpts, &reply, &bsonErr)
	if !ok {
		return nil, TransformError(bsonErr)
	}
	defer C.bson_destroy(&reply)

	result, err := mapstr.NewFromInterface(TransformBsonIntoGoString(&reply))
	if nil != err {
		return nil, err
	}

	replaceRst := &ReplaceOneResult{}
	err = result.MarshalJSONInto(replaceRst)
	return replaceRst, err
}

func (c *collection) Close() error {
	if nil != c.innerCollection {
		C.mongoc_collection_destroy(c.innerCollection)
		c.innerCollection = nil
	}

	return nil
}
