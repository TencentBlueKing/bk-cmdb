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

package driver

import (
	"context"

	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
)

var _ mongodb.CommonClient = (*client)(nil)

type client struct {
	createdByPool bool
	uri           string
	innerDB       *database
	innerClient   *mongo.Client
}

// NewClient create a mongoc client instance
func NewClient(uri string) mongodb.CommonClient {
	return &client{
		uri: uri,
	}
}

func (c *client) Open() error {

	cnnstr, err := connstring.Parse(c.uri)
	if nil != err {
		return err
	}

	c.innerClient, err = mongo.NewClient(c.uri)
	if nil != err {
		return err
	}

	if err := c.innerClient.Connect(context.TODO()); nil != err {
		return err
	}

	c.innerDB = newDatabase(c.innerClient.Database(cnnstr.Database))

	return nil
}

func (c *client) Close() error {

	if nil != c.innerClient {
		return nil
	}

	return c.innerClient.Disconnect(context.TODO())
}

func (c *client) Ping() error {

	return c.innerClient.Ping(context.TODO(), nil)
}

func (c *client) Database() mongodb.Database {
	return c.innerDB
}

func (c *client) Collection(collName string) mongodb.CollectionInterface {
	return newCollection(c.innerDB.innerDatabase, collName)
}

func (c *client) Session() mongodb.SessionOperation {
	return newSessionOperation(c)
}
