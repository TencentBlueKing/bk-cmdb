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

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func ConvLanguageToV3(language string) string {
	if "1" == language || "" == language {
		language = "zh-cn"
	} else if "2" == language {
		language = "en"
	}
	return language
}

// SliceStrToInt: 将字符串切片转换为整型切片
func SliceStrToInt(sliceStr []string) ([]int, error) {
	sliceInt := make([]int, 0)
	for _, idStr := range sliceStr {

		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return []int{}, err
		}
		sliceInt = append(sliceInt, id)
	}
	return sliceInt, nil
}

// ValidateFormData verify that the specified key is empty in formData
func ValidateFormData(formData url.Values, keys []string) (bool, string) {
	for _, key := range keys {
		if len(formData[key]) == 0 || formData[key][0] == "" {
			return false, fmt.Sprintf("param %s is empty!", key)
		}
	}
	return true, ""
}

// GetResMap extract http response body data and turn it into a Map
func GetResMap(resp *http.Response) (map[string]interface{}, error) {
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resMap := make(map[string]interface{})
	_ = json.Unmarshal(respData, &resMap)
	return resMap, nil
}
