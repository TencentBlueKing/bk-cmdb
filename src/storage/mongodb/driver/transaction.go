/*
* Tencent is pleased to support the open source community by making 蓝鲸 available.
* Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
* Licensed under the MIT License (the ",License"); you may not use this file except
* in compliance with the License. You may obtain a copy of the License at
* http://opensource.org/licenses/MIT
* Unless required by applicable law or agreed to in writing, software distributed under
* the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
* either express or implied. See the License for the specific language governing permissions and
* limitations under the License.
 */

package driver

import (
	"context"

	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var _ mongodb.Transaction = (*transaction)(nil)

type transaction struct {
	mongocli     *client
	innerSession mongo.Session
}

func newSessionTransaction(mongocli *client, clientSession mongo.Session) *transaction {
	return &transaction{
		mongocli:     mongocli,
		innerSession: clientSession,
	}
}

func (t *transaction) StartTransaction() error {
	return t.innerSession.StartTransaction()
}
func (t *transaction) AbortTransaction() error {
	return t.innerSession.AbortTransaction(context.TODO())
}
func (t *transaction) CommitTransaction() error {
	return t.innerSession.CommitTransaction(context.TODO())
}

func (t *transaction) Collection(collName string) mongodb.CollectionInterface {
	target := newCollectionWithSession(t.mongocli.innerDB.innerDatabase, t.innerSession, collName)
	return target
}

func (t *transaction) Close() error {
	t.innerSession.EndSession(context.TODO())
	return nil
}
