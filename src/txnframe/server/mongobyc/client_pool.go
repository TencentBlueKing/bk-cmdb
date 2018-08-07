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

// ClientPool the mongo client pool
type ClientPool interface {
	OpenCloser
	Pop() Client
	Push(targetClient Client)
}

// NewClientPool create a mongoc client pool instance
func NewClientPool(uri, database string) ClientPool {
	pool := new(clientPool)
	pool.client.uri = uri
	pool.client.dbName = database
	return pool
}

type clientPool struct {
	client
	pool *C.mongoc_client_pool_t
}

func (c *clientPool) Open() error {

	// create client
	var err C.bson_error_t
	uri := C.mongoc_uri_new_with_error(C.CString(c.client.uri), &err)
	if nil == uri {
		return TransformError(err)
	}

	c.pool = C.mongoc_client_pool_new(uri)
	if nil == c.pool {
		return fmt.Errorf("can not create a client pool instance")
	}

	C.mongoc_uri_destroy(uri)

	return nil
}

func (c *clientPool) Close() error {
	if nil != c.pool {
		C.mongoc_client_pool_destroy(c.pool)
	}
	return nil
}

func (c *clientPool) Pop() Client {

	mongocClient := C.mongoc_client_pool_pop(c.pool)

	innerClient := new(client)
	innerClient.innerClient = mongocClient
	innerClient.dbName = c.client.dbName
	innerClient.uri = c.client.uri

	cName := C.CString(c.client.dbName)
	innerClient.db = C.mongoc_client_get_database(mongocClient, cName)
	C.free(unsafe.Pointer(cName))

	return innerClient
}

func (c *clientPool) Push(targetClient Client) {

	switch tmp := targetClient.(type) {
	case *client:
		if nil != tmp.db {
			C.mongoc_database_destroy(tmp.db)
		}
		C.mongoc_client_pool_push(c.pool, tmp.innerClient)
	}

}
