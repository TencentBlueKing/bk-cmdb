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
import (
	"fmt"
	"unsafe"
)

// InitMongoc init the mongc lib
func InitMongoc() {
	C.mongoc_init()
}

// CleanupMongoc cleanup the mongc lib
func CleanupMongoc() {
	C.mongoc_cleanup()
}

// Client client for mongo
type Client interface {
	Ping() error
	SessionOperation() SessionOperation
	Collection(collName string) CollectionInterface
}

// CommonClient single client instance
type CommonClient interface {
	OpenCloser
	Client
}

// NewClient create a mongoc client instance
func NewClient(uri, database string) CommonClient {
	return &client{
		uri:    uri,
		dbName: database,
	}
}

type client struct {
	uri         string
	dbName      string
	db          *C.mongoc_database_t
	innerClient *C.mongoc_client_t
}

func (c *client) SessionOperation() SessionOperation {
	return newSessionOperation(c)
}

func (c *client) Open() error {

	// create client
	var err C.bson_error_t
	uri := C.mongoc_uri_new_with_error(C.CString(c.uri), &err)
	if nil == uri {
		return TransformError(err)
	}

	c.innerClient = C.mongoc_client_new_from_uri(uri)
	if nil == c.innerClient {
		return fmt.Errorf("can not create a client instance")
	}

	C.mongoc_uri_destroy(uri)

	// set app name
	cName := C.CString(c.dbName)
	C.mongoc_client_set_appname(c.innerClient, cName)

	// get database by name
	c.db = C.mongoc_client_get_database(c.innerClient, cName)
	C.free(unsafe.Pointer(cName))

	return nil
}

func (c *client) Close() error {
	if nil != c.innerClient {
		C.mongoc_client_destroy(c.innerClient)
	}

	if nil != c.db {
		C.mongoc_database_destroy(c.db)
	}
	return nil
}

func (c *client) Collection(collName string) CollectionInterface {
	return newCollectionWithoutSession(c, collName)
}

func (c *client) Ping() error {

	pingCStr := C.CString("ping")
	pingCmd := C.create_bcon_new_int32(pingCStr, 1)

	var reply C.bson_t
	var err C.bson_error_t
	adminCStr := C.CString("admin")
	ok := C.mongoc_client_command_simple(c.innerClient, adminCStr, pingCmd, nil, &reply, &err)

	C.free(unsafe.Pointer(pingCStr))
	C.free(unsafe.Pointer(adminCStr))
	if !ok {
		return TransformError(err)
	}

	return nil
}
