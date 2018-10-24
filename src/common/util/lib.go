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
	"context"
	"net/http"
	"reflect"
	"sync/atomic"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
)

func InStrArr(arr []string, key string) bool {
	for _, a := range arr {
		if key == a {
			return true
		}
	}
	return false
}

func GetLanguage(header http.Header) string {
	return header.Get(common.BKHTTPLanguage)
}

// GetActionLanguage returns language form hender
func GetActionLanguage(req *restful.Request) string {
	language := req.HeaderParameter(common.BKHTTPLanguage)
	if "" == language {
		language = "zh-cn"
	}
	blog.V(5).Infof("request language: %s, header: %v", language, req.Request.Header)
	return language
}

// GetActionUser returns user form hender
func GetActionUser(req *restful.Request) string {
	user := req.HeaderParameter(common.BKHTTPHeaderUser)
	return user
}

// GetActionOnwerID returns owner_uin form hender
func GetActionOnwerID(req *restful.Request) string {
	ownerID := req.HeaderParameter(common.BKHTTPOwnerID)
	return ownerID
}

func GetUser(header http.Header) string {
	return header.Get(common.BKHTTPHeaderUser)
}

func GetOwnerID(header http.Header) string {
	return header.Get(common.BKHTTPOwnerID)
}

func GetOwnerIDAndUser(header http.Header) (string, string) {
	return header.Get(common.BKHTTPOwnerID), header.Get(common.BKHTTPHeaderUser)
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

// GetHTTPCCRequestID return configcenter request id from http header
func GetHTTPCCRequestID(header http.Header) string {
	rid := header.Get(common.BKHTTPCCRequestID)
	return rid
}

// GetHTTPCCTransaction return configcenter request id from http header
func GetHTTPCCTransaction(header http.Header) string {
	rid := header.Get(common.BKHTTPCCTransactionID)
	return rid
}

// GetDBContext returns a new context that contains JoinOption
func GetDBContext(parent context.Context, header http.Header) context.Context {
	return context.WithValue(parent, common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
}

// IsNil returns whether value is nil value, including map[string]interface{}{nil}, *Struct{nil}
func IsNil(value interface{}) bool {
	rflValue := reflect.ValueOf(value)
	if rflValue.IsValid() {
		return rflValue.IsNil()
	}
	return true
}

type AtomicBool int32

func NewBool(yes bool) *AtomicBool {
	var n = AtomicBool(0)
	if yes {
		n = AtomicBool(1)
	}
	return &n
}

func (b *AtomicBool) Set() {
	atomic.StoreInt32((*int32)(b), 1)
}

func (b *AtomicBool) UnSet() {
	atomic.StoreInt32((*int32)(b), 0)
}

func (b *AtomicBool) IsSet() bool {
	return atomic.LoadInt32((*int32)(b)) == 1
}

func (b *AtomicBool) SetTo(yes bool) {
	if yes {
		atomic.StoreInt32((*int32)(b), 1)
	} else {
		atomic.StoreInt32((*int32)(b), 0)
	}
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
