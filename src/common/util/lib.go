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

package util

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"net/http"
	"reflect"
	"strings"

	restful "github.com/emicklei/go-restful"
)

func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func InStrArr(arr []string, key string) bool {
	for _, a := range arr {
		if key == a {
			return true
		}
	}
	return false
}

func ArrayUnique(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if !InArray(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

//StrArrayUnique get unique string array
func StrArrayUnique(a []string) (ret []string) {
	length := len(a)
	for i := 0; i < length; i++ {
		if !Contains(ret, a[i]) {
			ret = append(ret, a[i])
		}
	}
	return ret
}

//IntArrayUnique get unique int array
func IntArrayUnique(a []int) (ret []int) {
	length := len(a)
	for i := 0; i < length; i++ {
		if !ContainsInt(ret, a[i]) {
			ret = append(ret, a[i])
		}
	}
	return ret
}

func RemoveDuplicatesAndEmpty(slice []string) (ret []string) {
	for _, a := range slice {
		if strings.TrimSpace(a) != "" && !Contains(ret, a) {
			ret = append(ret, a)
		}
	}
	return
}

func StrArrDiff(slice1 []string, slice2 []string) []string {
	var diffStr []string
	for _, i := range slice1 {
		isIn := false
		for _, j := range slice2 {
			if i == j {
				isIn = true
				break
			}
		}
		if !isIn {
			diffStr = append(diffStr, i)
		}
	}
	return diffStr
}

// GetActionLanguage returns language form hender
func GetActionLanguage(req *restful.Request) string {
	language := req.HeaderParameter(common.BKHTTPLanguage)
	if "" == language {
		language = "zh-cn"
	}
	blog.Infof("request language: %s, header: %v", language, req.Request.Header)
	return language
}

// GetActionLanguage returns user form hender
func GetActionUser(req *restful.Request) string {
	user := req.HeaderParameter(common.BKHTTPHeaderUser)
	return user
}

// GetActionOnwerID returns owner_uin form hender
func GetActionOnwerID(req *restful.Request) string {
	ownerID := req.HeaderParameter(common.BKHTTPOwnerID)
	return ownerID
}

// GetActionOnwerIDAndUser returns owner_uin and user form hender
func GetActionOnwerIDAndUser(req *restful.Request) (string, string) {
	user := GetActionUser(req)
	ownerID := GetActionOnwerID(req)

	return ownerID, user
}

// GetActionLanguageByHTTPHeader return language from http header
func GetActionLanguageByHTTPHeader(header http.Header) string {
	language := header.Get(common.BKHTTPLanguage)
	if "" == language {
		return "zh-cn"
	}
	return language
}

// GetActionOnwerIDByHTTPHeader return owner from http header
func GetActionOnwerIDByHTTPHeader(header http.Header) string {
	ownerID := header.Get(common.BKHTTPOwnerID)
	return ownerID
}
