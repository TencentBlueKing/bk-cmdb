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
 
package http

import (
	"encoding/json"
	"errors"
)

type APIRespone struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InternalError(code int, message string) error {

	_, err := createRespone(code, message, make(map[string]interface{}))

	return err
}

//func GetRespWithoutData(code int, message string) string {
//
//	ret, _ := createRespone(code, message, make(map[string]interface{}))
//
//	return ret
//}
//
//func GetRespone(code int, message string, data interface{}) string {
//
//	ret, _ := createRespone(code, message, data)
//
//	return ret
//}

func createRespone(code int, message string, data interface{}) (string, error) {
	bResult := false
	if 0 == code {
		bResult = true
	}
	resp := APIRespone{bResult, code, message, data}
	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(b), errors.New(string(b))
}
