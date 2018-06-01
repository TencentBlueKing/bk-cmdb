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

package validator

import (
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"testing"
)

func TestValidator(t *testing.T) {
	date := "2017-02-01 01:01:1"
	re := util.IsTime(date)
	fmt.Println(re)
	addr := "127.0.0.1:50001"
	jsonr := `{"SetName":"aaaaaaaaaa","cpu_core_count":8}`
	type JsonT struct {
		SetName string
	}
	var hostData JsonT
	json.Unmarshal([]byte(jsonr), &hostData)

	hostMap := make(map[string]interface{})
	hostMap["SetName"] = hostData.SetName
	valid := NewValidMap("0", "app", addr, nil)
	_, err := valid.ValidMap(hostMap, "", 0)
	fmt.Println(err)
}
