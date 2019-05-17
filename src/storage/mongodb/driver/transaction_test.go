/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package driver_test

import (
	"context"
	"testing"

	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {

	executeCommand(t, func(dbClient mongodb.CommonClient) {
		session := dbClient.Session().Create()
		session.Open()
		defer func() {
			session.Close()
		}()

		err := session.StartTransaction()
		require.NoError(t, err)
		coll := session.Collection("cc_tranTest")

		dataItems := []interface{}{}
		for idx := 10; idx > 0; idx-- {
			dataItems = append(dataItems, bson.M{"key": "value-many_" + xid.New().String()})
		}
		err = coll.InsertMany(context.TODO(), dataItems, nil)
		require.NoError(t, err)
		//session.AbortTransaction()
		session.CommitTransaction()

	})

}
