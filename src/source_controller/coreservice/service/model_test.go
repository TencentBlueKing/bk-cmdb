/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCreateClassification(t *testing.T) {

	startCoreService(t, "127.0.0.1", 3308)

	client := httpclient.NewHttpClient()

	inputParams := []byte(`
		{  
		  "metadata":{"business_object":"biz"},
		  "datas":[
			{
				"bk_classification_id" : "` + xid.New().String() + `",
				"bk_classification_name" : ""
			}]
		}`)

	dataResult, err := client.POST("http://127.0.0.1:3308/api/v3/createmany/model/classification", defaultHeader, inputParams)
	require.NoError(t, err)
	require.NotNil(t, dataResult)

	clsResult := metadata.CreatedManyOptionResult{}
	err = json.Unmarshal(dataResult, &clsResult)
	require.NoError(t, err)
	t.Logf("data result:%v", clsResult.Data.Created)

}
