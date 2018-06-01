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
 
package actions_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestFavourite(t *testing.T) {
	id := testAddHostFavourite(t)
	t.Logf("id == %v ", id)
	defer testDeleteHostFavouriteByID(t, id)
	testGetHostFavourites(t)
	testEditHostFavouriteByID(t, id)
	testIncrHostFavouritesCount(t, id)
}
func testGetHostFavourites(t *testing.T) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"condition": map[string]interface{}{
			"IsDefault": 1,
			"Name":      "TestFavouriteName",
		},
		"limit": 10,
		"start": 0,
	}

	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/favorites/search", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)

}

func testAddHostFavourite(t *testing.T) (newid string) {
	server := CCAPITester(t)
	originData := map[string]interface{}{
		"Info":        `{"exact_search": true,"inner_ip": true,"outer_ip": true,"ip_list": [ "1.1.1.1","2.2.2.2"]}`,
		"inner_ip":    true,
		"outer_ip":    true,
		"QueryParams": `[{"object_id": "host","field": "operator_system","operator": "$nin","value": "123"}]`,
		"operator":    "$in",
		"Name":        "TestFavouriteName",
		"IsDefault":   1,
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/favorites", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
	id := body.Get("data.ID").String()
	require.NotZero(t, id, "response:[%s]", respbody)

	return id
}

func testEditHostFavouriteByID(t *testing.T, id string) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"Info":        `{"exact_search": true,"inner_ip": true,"outer_ip": true,"ip_list": [ "1.1.1.1","2.2.2.2"]}`,
		"inner_ip":    true,
		"outer_ip":    true,
		"QueryParams": `[{"object_id": "host","field": "operator_system","operator": "$nin","value": "123"}]`,
		"operator":    "$in",
		"Name":        "update",
		"IsDefault":   1,
	}
	reqbody, _ := json.Marshal(originData)
	r, err := http.NewRequest("PUT", server.URL+"/host/v1/hosts/favorites/"+id, bytes.NewBuffer(reqbody))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testDeleteHostFavouriteByID(t *testing.T, id string) {
	server := CCAPITester(t)

	originData := map[string]interface{}{}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("DELETE", server.URL+"/host/v1/hosts/favorites/"+id, bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testIncrHostFavouritesCount(t *testing.T, id string) {
	server := CCAPITester(t)

	originData := map[string]interface{}{}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("PUT", server.URL+"/host/v1/hosts/favorites/"+id+"/incr", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}
