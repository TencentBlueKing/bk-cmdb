/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// InterfaceToMapStr interface to mapstr
func InterfaceToMapStr(data interface{}) (map[string]interface{}, error) {

	resultData := make(map[string]interface{})
	switch value := data.(type) {
	case map[string]interface{}:
		return value, nil
	default:
		out, err := bson.Marshal(data)
		if err != nil {
			fmt.Printf("marshal error %v, data: %v", err, data)
			return nil, fmt.Errorf("marshal error %v", err)
		}
		if err = bson.Unmarshal(out, &resultData); err != nil {
			fmt.Printf("marshal error %v, data: %v", err, data)
			return nil, fmt.Errorf("unmarshal error %v", err)
		}
	}

	return resultData, nil
}

// CommonTenantTmp tenant template data struct for common type
type CommonTenantTmp TenantTmpData[map[string]interface{}]

// TenantTmpData tenant template data struct
type TenantTmpData[T any] struct {
	Type  string `json:"type"`
	IsPre bool   `json:"is_pre"`
	ID    int64  `json:"id"`
	Data  T      `json:"data"`
}

func main() {
	data := &CommonTenantTmp{
		Type:  "test",
		IsPre: true,
		ID:    1,
		Data: map[string]interface{}{
			"id":    2,
			"name":  "using",
			"study": "hhhh",
		},
	}

	result, err := InterfaceToMapStr(data)
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	fmt.Printf("test:%v", result["data"])
	fmt.Printf("%v", result)
}
