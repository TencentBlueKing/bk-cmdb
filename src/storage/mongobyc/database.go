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

// #include "mongo.h"
import "C"

import "unsafe"

// Database methods
type Database interface {
	Drop() error
	Name() string
	HasCollection(collName string) (bool, error)
	DropCollection(collName string) error
	CreateEmptyCollection(collName string) error
	GetCollectionNames() ([]string, error)
}

type database struct {
	dbName string
	db     *C.mongoc_database_t
}

func (d *database) Drop() error {
	return nil
}

func (d *database) Name() string {
	return d.dbName
}

func (d *database) HasCollection(collName string) (bool, error) {
	var err C.bson_error_t
	innerCollName := C.CString(collName)
	defer C.free(unsafe.Pointer(innerCollName))
	if C.mongoc_database_has_collection(d.db, innerCollName, &err) {
		return true, nil
	}
	if 0 != err.code {
		return false, TransformError(err)
	}
	return false, nil
}

func (d *database) DropCollection(collName string) error {
	return nil
}
func (d *database) CreateEmptyCollection(collName string) error {

	innerCollName := C.CString(collName)
	defer C.free(unsafe.Pointer(innerCollName))
	var bsonErr C.bson_error_t
	innerCollection := C.mongoc_database_create_collection(d.db, innerCollName, nil, &bsonErr)
	if nil == innerCollection {
		return TransformError(bsonErr)
	}
	C.mongoc_collection_destroy(innerCollection)
	return nil
}

func (d *database) GetCollectionNames() ([]string, error) {

	//readConcern := C.mongoc_read_concern_new()
	//majority := C.CString(C.MONGOC_READ_CONCERN_LEVEL_MAJORITY)
	//C.mongoc_read_concern_set_level(readConcern, majority)
	opts := C.bson_new()
	//C.mongoc_read_concern_append(readConcern, opts)
	var bsonErr C.bson_error_t
	names := C.mongoc_database_get_collection_names_with_opts(d.db, opts, &bsonErr)
	if nil == names {
		return nil, TransformError(bsonErr)
	}

	results := []string{}
	idx := 0
	pointerSize := unsafe.Sizeof((*C.char)(unsafe.Pointer(names)))
	for {

		name := (**C.char)(unsafe.Pointer((uintptr(unsafe.Pointer(names)) + pointerSize*uintptr(idx))))
		if nil != *name {
			results = append(results, C.GoString(*name))
			idx++
			continue
		}

		break
	}

	//	C.mongoc_read_concern_destroy(readConcern)
	//C.free(unsafe.Pointer(majority))
	C.bson_destroy(opts)
	C.bson_strfreev(names)

	return results, nil
}
func (d *database) Close() error {
	if nil != d.db {
		C.mongoc_database_destroy(d.db)
		d.db = nil
	}

	return nil
}
