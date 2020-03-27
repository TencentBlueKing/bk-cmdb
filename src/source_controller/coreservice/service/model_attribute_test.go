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

func createModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {

	modelAttributeItems := metadata.CreateModelAttributes{
		Attributes: []metadata.Attribute{
			metadata.Attribute{
				PropertyID:   modelAttributeID,
				PropertyName: xid.New().String(),
				ObjectID:     modelID,
			},
		},
	}

	inputParams, err := json.Marshal(modelAttributeItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create one model attribute:%s", inputParams)

	dataResult, err := client.POST("http://host:port/api/v3/create/model/"+modelID+"/attributes", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.CreatedManyOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)

	for _, attrCreated := range modelResult.Data.CreateManyInfoResult.Created {
		require.NotEqual(t, uint64(0), attrCreated.ID)
	}

	for _, attrException := range modelResult.Data.Exceptions {
		require.NotEqual(t, int64(0), attrException.Code)
		require.NotEmpty(t, attrException.Message)
		require.NotNil(t, attrException.Data)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("create one model attribute result:%s", resultStr)

}

func setModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {

	modelItems := metadata.SetModelAttributes{
		Attributes: []metadata.Attribute{
			metadata.Attribute{
				PropertyID:   modelAttributeID,
				PropertyName: xid.New().String(),
				ObjectID:     modelID,
			},
		},
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set one model attributes:%s", inputParams)

	dataResult, err := client.POST("http://host:port/api/v3/set/model/"+modelID+"/attributes", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.SetDataResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	for _, item := range modelResult.Created {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range modelResult.Updated {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range modelResult.Exceptions {
		require.NotEqual(t, uint64(0), item.Code)
		require.NotNil(t, item.Data)
		require.NotEmpty(t, item.Data)
		require.NotEmpty(t, item.Message)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("set one model result:%s", resultStr)
}

func queryModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AttributeFieldPropertyID).Eq(modelAttributeID))
	modelItems := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("read some model attributes:%s", inputParams)

	dataResult, err := client.POST("http://host:port/api/v3/read/model/"+modelID+"/attributes", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.ReadModelAttrResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.Equal(t, modelResult.Data.Count, int64(len(modelResult.Data.Info)))
	require.NotNil(t, modelResult.Data.Info)

	for _, item := range modelResult.Data.Info {
		require.NotEqual(t, int64(0), item.ID)
		require.NotEmpty(t, item.OwnerID)
		require.NotEmpty(t, item.PropertyID)
		require.NotEmpty(t, item.ObjectID)
		require.NotNil(t, item.LastTime)
		require.NotNil(t, item.CreateTime)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("read some model attributes result:%s", resultStr)
}

func updateModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AttributeFieldPropertyID).Eq(modelAttributeID))
	modelItems := metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.AttributeFieldPropertyName: "update_" + modelID,
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("update some model attributes:%s", inputParams)

	dataResult, err := client.PUT("http://host:port/api/v3/update/model/"+modelID+"/attributes", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.UpdatedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)
	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("update some model attributes result:%s", resultStr)
}

func deleteModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AttributeFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AttributeFieldPropertyID).Eq(modelAttributeID))
	modelItems := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("delete some model attributes:%s", inputParams)

	dataResult, err := client.DELETE("http://host:port/api/v3/delete/model/"+modelID+"/attributes", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.DeletedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("delete some model attributes result:%s", resultStr)
}

func deleteCascadeModelAttributes(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeID string) {
	// TODO
}

func TestModelAttributeCRUD(t *testing.T) {

	// base
	startCoreService(t, "host", 3308)
	client := httpclient.NewHttpClient()

	classID := xid.New().String()
	createOneClassification(t, client, classID)

	modelID := xid.New().String()
	createModel(t, client, modelID, classID)

	// create many
	modelAttributeID := xid.New().String()
	createModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)
	updateModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)
	deleteModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)

	// create one
	modelAttributeID = xid.New().String()
	createModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)
	updateModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)
	deleteModelAttributes(t, client, modelID, modelAttributeID)
	queryModelAttributes(t, client, modelID, modelAttributeID)

	deleteModel(t, client, modelID, classID)
	deleteClassification(t, client, classID)

}
