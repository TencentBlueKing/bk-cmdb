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

package service_test

import (
	"encoding/json"
	"testing"

	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func createModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {

	modelItems := metadata.CreateModel{
		Spec: metadata.Object{
			ObjectID:   modelID,
			ObjCls:     classID,
			ObjectName: "name_" + modelID,
		},
		Attributes: []metadata.Attribute{
			metadata.Attribute{
				PropertyID:   xid.New().String(),
				PropertyName: xid.New().String(),
				ObjectID:     modelID,
			},
		},
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create one model:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/create/model", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.CreatedOneOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, 0, modelResult.Data.Created.ID)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("create one model result:%s", resultStr)

}

func setModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {

	modelItems := metadata.SetModel{
		Spec: metadata.Object{
			ObjectID:   modelID,
			ObjCls:     classID,
			ObjectName: "name_" + modelID,
		},
		Attributes: []metadata.Attribute{
			metadata.Attribute{
				PropertyID:   xid.New().String(),
				PropertyName: xid.New().String(),
				ObjectID:     modelID,
			},
		},
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set one model:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/set/model", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.SetDataResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)

	for _, modelCreated := range modelResult.Created {
		require.NotEqual(t, uint64(0), modelCreated.ID)
	}
	for _, modelUpdated := range modelResult.Updated {
		require.NotEqual(t, uint64(0), modelUpdated.ID)
	}
	for _, modelException := range modelResult.Exceptions {
		require.NotEqual(t, int64(0), modelException.Code)
		require.NotNil(t, modelException.Data)
		require.NotEmpty(t, modelException.Message)
	}
	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("set one model result:%s", resultStr)
}

func queryModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ModelFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.ModelFieldObjCls).Eq(classID))
	modelItems := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("read some models:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/read/model", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.ReadModelResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.Equal(t, modelResult.Data.Count, int64(len(modelResult.Data.Info)))
	for _, item := range modelResult.Data.Info {
		require.NotEqual(t, int64(0), item.Spec.ID)
		require.NotEmpty(t, item.Spec.OwnerID)
		require.NotEmpty(t, item.Spec.ObjectID)
		require.NotEmpty(t, item.Spec.ObjCls)
		require.NotNil(t, item.Spec.CreateTime)
		require.NotNil(t, item.Spec.LastTime)
		t.Logf("create time:%s", item.Spec.CreateTime.String())

		for _, attr := range item.Attributes {
			require.NotEmpty(t, attr.OwnerID)
			require.NotEmpty(t, attr.PropertyID)
			require.NotEqual(t, int64(0), attr.ID)
			require.NotNil(t, attr.CreateTime)
			require.NotNil(t, attr.LastTime)
		}
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("read some models result:%s", resultStr)
}

func updateModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ModelFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.ModelFieldObjCls).Eq(classID))
	modelItems := metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.ModelFieldObjectName: "update_" + modelID,
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("update some models:%s", inputParams)

	dataResult, err := client.PUT("http://127.0.0.1:3308/api/v3/update/model", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.UpdatedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("update some models result:%s", resultStr)
}

func deleteModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ModelFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.ModelFieldObjCls).Eq(classID))
	modelItems := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("delete some models:%s", inputParams)

	dataResult, err := client.DELETE("http://127.0.0.1:3308/api/v3/delete/model", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.DeletedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("delete some models result:%s", resultStr)
}

func deleteCascadeModel(t *testing.T, client *httpclient.HttpClient, modelID, classID string) {
	// TODO
}

func TestModelCRUD(t *testing.T) {

	// base
	startCoreService(t, "127.0.0.1", 3308)
	client := httpclient.NewHttpClient()

	classID := xid.New().String()
	createOneClassification(t, client, classID)

	// create many
	modelID := xid.New().String()
	createModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)
	updateModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)
	deleteModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)

	// create one
	modelID = xid.New().String()
	createModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)
	updateModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)
	deleteModel(t, client, modelID, classID)
	queryModel(t, client, modelID, classID)

	deleteClassification(t, client, classID)

}
