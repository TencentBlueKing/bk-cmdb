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

func createAssociationKind(t *testing.T, client *httpclient.HttpClient, asstKindID string) {
	data := `{
		"data":{
				"bk_asst_id": "` + asstKindID + `",
				"bk_asst_name": "` + "name_" + asstKindID + `",
				"src_des": "属于",
				"dest_des": "被属于",
				"direction": "none",
		}
		
	}`

	client.POST("http://127.0.0.1:3308/api/v3/create/associationkind", defaultHeader, []byte(data))
}

func createModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, asstModelID, modelAssociationID string) {

	asstKindID := xid.New().String()
	createAssociationKind(t, client, asstKindID)
	modelAssociationItems := metadata.CreateModelAssociation{
		Spec: metadata.Association{
			ObjectID:             modelID,
			AssociationName:      modelAssociationID,
			AssociationAliasName: "name_" + modelAssociationID,
			AsstObjID:            asstModelID,
			AsstKindID:           asstKindID,
		},
	}

	inputParams, err := json.Marshal(modelAssociationItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create model association:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/create/modelassociation", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.CreatedOneOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.True(t, modelResult.Result, modelResult.ErrMsg)
	require.NotEqual(t, uint64(0), modelResult.Data.Created.ID)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("create model association result:%s", resultStr)

}

func setModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, asstModelID, modelAssociationID string) {

	modelAssociationItems := metadata.SetModelAssociation{
		Spec: metadata.Association{
			ObjectID:             modelID,
			AssociationName:      modelAssociationID,
			AssociationAliasName: "name_" + modelAssociationID,
		},
	}

	inputParams, err := json.Marshal(modelAssociationItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set model association:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/set/modelassociation", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.SetDataResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)

	for _, item := range modelResult.Created {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range modelResult.Updated {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range modelResult.Exceptions {
		require.NotEqual(t, int64(0), item.Code)
		require.NotNil(t, item.Data)
		require.NotEmpty(t, item.Data)
		require.NotEmpty(t, item.Message)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("set model association result:%s", resultStr)
}

func queryModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, modelAssociationID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AssociationFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AssociationFieldAsstID).Eq(modelAssociationID))

	modelItems := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("read model association:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/read/modelassociation", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.ReadModelAssociationResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.Equal(t, modelResult.Data.Count, uint64(len(modelResult.Data.Info)))
	for _, item := range modelResult.Data.Info {
		require.NotEmpty(t, item.OwnerID)
		require.NotEmpty(t, item.ObjectID)
		require.NotEmpty(t, item.AssociationName)
		require.NotEqual(t, int64(0), item.ID)
	}

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("read  model association result:%s", resultStr)
}

func updateModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, modelAssociationID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AssociationFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AssociationFieldAsstID).Eq(modelAssociationID))
	modelItems := metadata.UpdateOption{
		Data: mapstr.MapStr{
			"bk_obj_asst_name": "update_" + modelID,
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("update model association:%s", inputParams)

	dataResult, err := client.PUT("http://127.0.0.1:3308/api/v3/update/modelassociation", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.UpdatedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("update model association result:%s", resultStr)
}

func deleteModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, modelAssociationID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.AssociationFieldObjectID).Eq(modelID))
	cond.Element(mongo.Field(metadata.AssociationFieldAsstID).Eq(modelAssociationID))
	modelItems := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(modelItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("delete model association:%s", inputParams)

	dataResult, err := client.DELETE("http://127.0.0.1:3308/api/v3/delete/modelassociation", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	modelResult := &metadata.DeletedOptionResult{}
	err = json.Unmarshal(dataResult, modelResult)
	require.NoError(t, err)
	require.NotNil(t, modelResult)
	require.NotEqual(t, uint64(0), modelResult.Data.Count)

	resultStr, err := json.Marshal(modelResult)
	require.NoError(t, err)

	t.Logf("delete model association result:%s", resultStr)
}

func deleteCascadeModelAssociation(t *testing.T, client *httpclient.HttpClient, modelID, modelAssociationID string) {
	// TODO
}

func TestModelAssociationCRUD(t *testing.T) {

	// base
	startCoreService(t, "127.0.0.1", 3308)
	client := httpclient.NewHttpClient()

	classID := xid.New().String()
	createOneClassification(t, client, classID)

	modelID := xid.New().String()
	createModel(t, client, modelID, classID)
	asstModelID := xid.New().String()
	createModel(t, client, asstModelID, classID)

	modelAssociationID := xid.New().String()
	createModelAssociation(t, client, modelID, asstModelID, modelAssociationID)
	queryModelAssociation(t, client, modelID, modelAssociationID)
	updateModelAssociation(t, client, modelID, modelAssociationID)
	queryModelAssociation(t, client, modelID, modelAssociationID)
	deleteModelAssociation(t, client, modelID, modelAssociationID)
	queryModelAssociation(t, client, modelID, modelAssociationID)

	deleteModel(t, client, modelID, classID)
	deleteModel(t, client, asstModelID, classID)
	deleteClassification(t, client, classID)

}
