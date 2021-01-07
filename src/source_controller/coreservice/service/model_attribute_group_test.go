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

func createModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {

	modelAttributeGroupItems := metadata.CreateModelAttributeGroup{
		Data: metadata.Group{
			ObjectID:  modelID,
			GroupID:   modelAttributeGroupID,
			GroupName: "name_" + modelAttributeGroupID,
		},
	}

	inputParams, err := json.Marshal(modelAttributeGroupItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create model attribute group:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/create/model/"+modelID+"/group", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.CreatedOneOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("create one model attribute group result:%s", resultStr)

}

func setModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {

	modelAttributeGroupItems := metadata.SetModelAttributeGroup{
		Data: metadata.Group{
			ObjectID:  modelID,
			GroupID:   modelAttributeGroupID,
			GroupName: "name_" + modelAttributeGroupID,
		},
	}

	inputParams, err := json.Marshal(modelAttributeGroupItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set one model attribute group:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/set/model/"+modelID+"/group", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.SetDataResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("set one model attribute group result:%s", resultStr)
}

func queryModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.GroupFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.GroupFieldGroupID).Eq(modelAttributeGroupID))
	modelItems := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("read some model attribute group:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/read/model/"+modelID+"/group", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.ReadModelAttributeGroupResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.Equal(t, modelResult.Data.Count, int64(len(modelResult.Data.Info)))
	for _, item := range modelResult.Data.Info {
		require.NotEmpty(t, item.OwnerID)
		require.NotEmpty(t, item.ObjectID)
		require.NotEmpty(t, item.GroupID)
		require.NotEqual(t, int64(0), item.ID)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("read some model attribute group result:%s", resultStr)
}

func updateModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.GroupFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.GroupFieldGroupID).Eq(modelAttributeGroupID))
	modelItems := metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.GroupFieldGroupName: "update_" + modelID,
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("update some model attribute group:%s", inputParams)

	dataResult, err := client.PUT("http://127.0.0.1:3308/api/v3/update/model/"+modelID+"/group", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.UpdatedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("update some model attribute group result:%s", resultStr)
}

func deleteModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.GroupFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.GroupFieldGroupID).Eq(modelAttributeGroupID))
	modelItems := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("delete some model attribute group:%s", inputParams)

	dataResult, err := client.DELETE("http://127.0.0.1:3308/api/v3/delete/model/"+modelID+"/group", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.DeletedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("delete some model attribute group result:%s", resultStr)
}

func deleteCascadeModelAttributeGroup(t *testing.T, client *httpclient.HttpClient, modelID, modelAttributeGroupID string) {
	// TODO
}

func TestModelAttributeGroupCRUD(t *testing.T) {

	// base
	startCoreService(t, "127.0.0.1", 3308)
	client := httpclient.NewHttpClient()

	classID := xid.New().String()
	createOneClassification(t, client, classID)

	modelID := xid.New().String()
	createModel(t, client, modelID, classID)

	// create many
	modelAttributeGroupID := xid.New().String()
	createModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	updateModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	deleteModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)

	// create one
	modelAttributeGroupID = xid.New().String()
	createModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	updateModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	deleteModelAttributeGroup(t, client, modelID, modelAttributeGroupID)
	queryModelAttributeGroup(t, client, modelID, modelAttributeGroupID)

	deleteModel(t, client, modelID, classID)
	deleteClassification(t, client, classID)

}
