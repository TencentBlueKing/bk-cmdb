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

package converter

import (
	"configcenter/src/common"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var v3ItemMap map[string]interface{}
var v2ItemMap map[string]interface{}

const (
	v3Data = "v3_test_data"
	v3Time = "2018-1-23 01:02:03"
	v2Time = "2018-01-23 01:02:03"
)

func init() {
	v3ItemMap = map[string]interface{}{
		common.BKAppIDField:    json.Number("0"),
		common.BKDefaultField:  json.Number("0"),
		common.BKAppNameField:  "app-name-test",
		common.LastTimeField:   v3Time,
		common.CreateTimeField: v3Time,
		common.BKTimeZoneField: v3Time,
	}

	v2ItemMap = map[string]interface{}{
		"ApplicationName": "app-name-test",
		"Default":         "0",
		"ApplicationID":   "0",
		"LastTime":        v2Time,
		"CreateTime":      v2Time,
		"TimeZone":        v3Time,
	}
}

func TestDecorateUserName(t *testing.T) {
	name := "test-name"
	vName := DecorateUserName(name)
	assert.Equal(t, name, vName)
}

func TestGetResDataV3(t *testing.T) {
	v3Resp := make(map[string]interface{}, 0)
	v3Resp["result"] = true
	v3Resp["data"] = v3Data

	by, _ := json.Marshal(v3Resp)
	v3, err := getResDataV3(string(by))
	assert.Nil(t, err)
	v3data, ok := v3.(string)
	assert.Equal(t, true, ok)
	assert.Equal(t, v3Data, v3data)
}

func TestConvertToV2Time(t *testing.T) {
	nTime := convertToV2Time(v3Time)
	assert.Equal(t, v2Time, nTime)
}

func TestConvertFieldsIntToStr(t *testing.T) {
	fields := []string{common.BKAppIDField, common.BKDefaultField, "nilKey"}

	v3ItemMap["nilKey"] = nil
	nItem, err := convertFieldsIntToStr(v3ItemMap, fields)
	assert.Nil(t, err)

	for _, i := range fields {
		_, ok := nItem[i].(string)
		assert.Equal(t, true, ok)
	}
}

func TestConvertOneApp(t *testing.T) {
	v2Item, err := convertOneApp(v3ItemMap)
	assert.Nil(t, err)
	assert.Equal(t, v2ItemMap["LastTime"], v2Item["LastTime"])
	assert.Equal(t, v2ItemMap["CreateTime"], v2Item["CreateTime"])
	assert.Equal(t, v2ItemMap["TimeZone"], v2Item["TimeZone"])
	assert.Equal(t, v2ItemMap["ApplicationID"], v2Item["ApplicationID"])
	assert.Equal(t, v2ItemMap["Default"], v2Item["Default"])
	assert.Equal(t, v2ItemMap["ApplicationName"], v2Item["ApplicationName"])
}
