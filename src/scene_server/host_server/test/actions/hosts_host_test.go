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

	"github.com/stretchr/testify/assert"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func testHostModuleRelation(t *testing.T, ids []int) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostID":        ids,
		"ModuleID":      []int{1},
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/modules", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testMoveIDleModule(t *testing.T, ids []int) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostID":        ids,
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/emptymodule", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testMoveFaultModule(t *testing.T, ids []int) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostID":        ids,
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/faultmodule", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testMoveHostToResourcePool(t *testing.T, ids []int) {
	server := CCAPITester(t)
	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostID":        ids,
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("PUT", server.URL+"/host/v1/hosts/resource", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testHostAdd(t *testing.T) {
	server := CCAPITester(t)
	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostInfo": map[int]map[string]interface{}{
			1: map[string]interface{}{
				"InnerIP":      "127.0.0.1",
				"OuterMAC":     "TestMac",
				"OSName":       "TestOS",
				"Mem":          64,
				"AgentVersion": "",
				"CreateTime":   "2018-01-30 15:57:00",
				"AssetID":      "",
				"SN":           "",
				"OuterIP":      "127.0.0.1",
				"HostName":     "TestHost",
				"MAC":          "TestMac",
				"Comment":      "",
				"ServiceTerm":  1,
				"Operator":     "",
				"BakOperator":  "",
				"OSVersion":    "",
				"CPU":          4,
				"CPUMhz":       2000,
				"CPUModule":    "",
				"Disk":         500,
				"aaa":          "1",
				"OSBit":        "2",
			},
		},
	}
	reqbody, err := json.Marshal(originData)
	require.NoError(t, err)
	r, err := http.NewRequest("POST", server.URL+"/host/v1/hosts/addhost", bytes.NewBuffer(reqbody))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func testAssignHostToApp(t *testing.T, ids []int) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"ApplicationID": 1,
		"HostID":        ids,
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/hosts/assgin", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func TestAddHostFromAgent(t *testing.T) {
	server := CCAPITester(t)

	originData := map[string]interface{}{}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/host/add/agent", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}

func TestAssignHostToAppModule(t *testing.T) {
	server := CCAPITester(t)

	originData := map[string]interface{}{
		"Ips":        []string{"1127.0.0.1"},
		"HostName":   []string{},
		"ModuleName": "string",
		"SetName":    "string",
		"AppName":    "string",
		"OsType":     "TestOS",
		"OwnerID":    "1",
	}
	reqbody, _ := json.Marshal(originData)
	r, _ := http.NewRequest("POST", server.URL+"/host/v1/host/add/module", bytes.NewBuffer(reqbody))
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	assert.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
}
