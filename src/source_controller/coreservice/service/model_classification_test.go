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

func setManyClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {
	// set many
	classItems := metadata.SetManyModelClassification{
		Data: []metadata.Classification{
			metadata.Classification{
				Metadata: metadata.Metadata{
					Label: metadata.Label{metadata.LabelBusinessID: "test_biz"},
				},
				ClassificationID: classificationID,
			},
		},
	}

	inputParams, err := json.Marshal(classItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set many classificaiton:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/setmany/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.SetOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)

	for _, item := range clsResult.Data.Created {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range clsResult.Data.Updated {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range clsResult.Data.Exceptions {
		require.NotEmpty(t, item.Message)
		require.NotNil(t, item.Data)
		require.NotEqual(t, int64(0), item.Code)
	}

	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("set many data result:%s", resultStr)
}

func setOneClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {

	// create one
	classItems := metadata.SetOneModelClassification{
		Data: metadata.Classification{
			Metadata: metadata.Metadata{
				Label: metadata.Label{metadata.LabelBusinessID: "test_biz"},
			},
			ClassificationID: classificationID,
		},
	}

	inputParams, err := json.Marshal(classItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("set one classificaiton:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/set/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.SetOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)
	for _, item := range clsResult.Data.Created {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range clsResult.Data.Updated {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range clsResult.Data.Exceptions {
		require.NotEqual(t, uint64(0), item.Code)
		require.NotEmpty(t, item.Message)
		require.NotNil(t, item.Data)
		require.NotEmpty(t, item.Data)
	}

	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("set one data result:%s", resultStr)
}

func createManyClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {

	// create many
	classItems := metadata.CreateManyModelClassifiaction{
		Data: []metadata.Classification{
			metadata.Classification{
				Metadata: metadata.Metadata{
					Label: metadata.Label{metadata.LabelBusinessID: "test_biz"},
				},
				ClassificationID: classificationID,
			},
		},
	}

	inputParams, err := json.Marshal(classItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create many classificaiton:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/createmany/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.CreatedManyOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)

	for _, item := range clsResult.Data.Created {
		require.NotEqual(t, uint64(0), item.ID)
	}

	for _, item := range clsResult.Data.Repeated {
		require.NotNil(t, item.Data)
		require.NotEmpty(t, item.Data)
	}

	for _, item := range clsResult.Data.Exceptions {
		require.NotNil(t, item.Data)
		require.NotEmpty(t, item.Data)
		require.NotEmpty(t, item.Message)
	}

	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("create many data result:%s", resultStr)

}

func createOneClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {

	// create one
	classItems := metadata.CreateOneModelClassification{
		Data: metadata.Classification{
			Metadata: metadata.Metadata{
				Label: metadata.Label{metadata.LabelBusinessID: "test_biz"},
			},
			ClassificationID: classificationID,
		},
	}

	inputParams, err := json.Marshal(classItems)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("create one classificaiton:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/create/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.CreatedOneOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), clsResult.Data.Created.ID)

	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("create one data result:%s", resultStr)
}

func queryClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ClassFieldClassificationID).Eq(classificationID))
	queryCond := metadata.QueryCondition{
		Fields: []string{},
		SortArr: []metadata.SearchSort{
			metadata.SearchSort{
				IsDsc: true,
				Field: metadata.ClassFieldClassificationID,
			},
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(queryCond)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("query classificaiton:%s", inputParams)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/read/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.ReadModelClassifitionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)
	require.Equal(t, clsResult.Data.Count, int64(len(clsResult.Data.Info)))
	for _, item := range clsResult.Data.Info {
		require.NotEqual(t, int64(0), item.ID)
		require.NotEmpty(t, item.OwnerID)
		require.NotEmpty(t, item.ClassificationID)
	}

	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("query data result:%s", resultStr)
}

func updateClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {

	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ClassFieldClassificationID).Eq(classificationID))
	cond.Element(mongo.Field(metadata.ClassFieldClassificationSupplierAccount).Eq("0"))
	queryCond := metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.ClassFieldClassificationName: "update_" + classificationID,
		},
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(queryCond)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("update classificaiton:%s", inputParams)

	dataResult, err := client.PUT("http://127.0.0.1:3308/api/v3/update/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.UpdatedOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)
	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("update data result:%s", resultStr)
}

func deleteClassification(t *testing.T, client *httpclient.HttpClient, classificationID string) {
	cond := mongo.NewCondition()
	cond.Element(mongo.Field(metadata.ClassFieldClassificationID).Eq(classificationID))
	queryCond := metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}

	inputParams, err := json.Marshal(queryCond)
	require.NoError(t, err)
	require.NotNil(t, inputParams)
	t.Logf("delete classificaiton:%s", inputParams)

	dataResult, err := client.DELETE("http://127.0.0.1:3308/api/v3/delete/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.DeletedOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)
	resultStr, err := json.Marshal(clsResult)
	require.NoError(t, err)
	t.Logf("delete data result:%s", resultStr)
}

func TestClassificationCRUD(t *testing.T) {

	// base
	startCoreService(t, "127.0.0.1", 3308)
	client := httpclient.NewHttpClient()
	classID := xid.New().String()
	t.Logf("create many:%s", classID)
	createManyClassification(t, client, classID)
	queryClassification(t, client, classID)
	updateClassification(t, client, classID)
	queryClassification(t, client, classID)
	deleteClassification(t, client, classID)
	queryClassification(t, client, classID)

	classID = xid.New().String()
	t.Logf("create one:%s", classID)
	createOneClassification(t, client, classID)
	queryClassification(t, client, classID)
	updateClassification(t, client, classID)
	queryClassification(t, client, classID)
	deleteClassification(t, client, classID)
	queryClassification(t, client, classID)

	classID = xid.New().String()
	t.Logf("set many:%s", classID)
	setManyClassification(t, client, classID)
	queryClassification(t, client, classID)
	updateClassification(t, client, classID)
	queryClassification(t, client, classID)
	deleteClassification(t, client, classID)
	queryClassification(t, client, classID)

	classID = xid.New().String()
	t.Logf("set one:%s", classID)
	setOneClassification(t, client, classID)
	queryClassification(t, client, classID)
	updateClassification(t, client, classID)
	queryClassification(t, client, classID)
	deleteClassification(t, client, classID)
	queryClassification(t, client, classID)
}
