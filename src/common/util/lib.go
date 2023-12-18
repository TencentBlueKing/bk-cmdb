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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// InStrArr TODO
func InStrArr(arr []string, key string) bool {
	for _, a := range arr {
		if key == a {
			return true
		}
	}
	return false
}

// ExtractRequestIDFromContext TODO
func ExtractRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	rid := ctx.Value(common.ContextRequestIDField)
	ridValue, ok := rid.(string)
	if ok == true {
		return ridValue
	}
	return ""
}

// ExtractOwnerFromContext TODO
func ExtractOwnerFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	owner := ctx.Value(common.ContextRequestOwnerField)
	ownerValue, ok := owner.(string)
	if ok == true {
		return ownerValue
	}
	return ""
}

// NewContextFromGinContext new context from gin context
func NewContextFromGinContext(c *gin.Context) context.Context {
	header := c.Request.Header
	ctx := c.Request.Context()
	ctx = SetContextValueByHTTPHeader(ctx, header)
	return ctx
}

// NewContextFromHTTPHeader new context from http header
func NewContextFromHTTPHeader(header http.Header) context.Context {
	return SetContextValueByHTTPHeader(context.Background(), header)
}

// SetContextValueByHTTPHeader set context value by http header
func SetContextValueByHTTPHeader(ctx context.Context, header http.Header) context.Context {
	ctx = context.WithValue(ctx, common.ContextRequestIDField, httpheader.GetRid(header))
	ctx = context.WithValue(ctx, common.ContextRequestUserField, httpheader.GetUser(header))
	ctx = context.WithValue(ctx, common.ContextRequestOwnerField, httpheader.GetSupplierAccount(header))
	return ctx
}

// ExtractRequestUserFromContext TODO
func ExtractRequestUserFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	user := ctx.Value(common.ContextRequestUserField)
	userValue, ok := user.(string)
	if ok == true {
		return userValue
	}
	return ""
}

// AtomicBool TODO
type AtomicBool int32

// NewBool TODO
func NewBool(yes bool) *AtomicBool {
	var n = AtomicBool(0)
	if yes {
		n = AtomicBool(1)
	}
	return &n
}

// SetIfNotSet TODO
func (b *AtomicBool) SetIfNotSet() bool {
	return atomic.CompareAndSwapInt32((*int32)(b), 0, 1)
}

// Set TODO
func (b *AtomicBool) Set() {
	atomic.StoreInt32((*int32)(b), 1)
}

// UnSet TODO
func (b *AtomicBool) UnSet() {
	atomic.StoreInt32((*int32)(b), 0)
}

// IsSet TODO
func (b *AtomicBool) IsSet() bool {
	return atomic.LoadInt32((*int32)(b)) == 1
}

// SetTo TODO
func (b *AtomicBool) SetTo(yes bool) {
	if yes {
		atomic.StoreInt32((*int32)(b), 1)
	} else {
		atomic.StoreInt32((*int32)(b), 0)
	}
}

// IntSlice TODO
type IntSlice []int

// Len 用于排序
func (p IntSlice) Len() int { return len(p) }

// Less 用于排序
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }

// Swap 用于排序
func (p IntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Int64Slice TODO
type Int64Slice []int64

// Len 用于排序
func (p Int64Slice) Len() int { return len(p) }

// Less 用于排序
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }

// Swap 用于排序
func (p Int64Slice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// GenerateRID TODO
func GenerateRID() string {
	unused := "0000"
	id := xid.New()
	return fmt.Sprintf("cc%s%s", unused, id.String())
}

// Int64Join []int64 to string
func Int64Join(data []int64, separator string) string {
	var ret string
	for _, item := range data {
		ret += strconv.FormatInt(item, 10) + separator
	}
	return strings.Trim(ret, separator)
}

// BuildMongoField build mongodb sub item field key
func BuildMongoField(key ...string) string {
	return strings.Join(key, ".")
}

// BuildMongoSyncItemField build mongodb sub item synchronize field key
func BuildMongoSyncItemField(key string) string {
	return BuildMongoField(common.MetadataField, common.MetaDataSynchronizeField, key)
}

// GetDefaultCCError TODO
func GetDefaultCCError(header http.Header) errors.DefaultCCErrorIf {
	globalCCError := errors.GetGlobalCCError()
	if globalCCError == nil {
		return nil
	}
	language := httpheader.GetLanguage(header)
	return globalCCError.CreateDefaultCCErrorIf(language)
}

// SetHTTPReadPreference 在header头中设置mongodb read preference，这个是给调用子流程使用
func SetHTTPReadPreference(header http.Header, mode common.ReadPreferenceMode) http.Header {
	header.Set(common.ReadReferenceKey, mode.String())
	return header
}

// SetDBReadPreference 在context中设置mongodb read preference，给dal使用
func SetDBReadPreference(ctx context.Context, mode common.ReadPreferenceMode) context.Context {
	ctx = context.WithValue(ctx, common.ReadReferenceKey, mode.String())
	return ctx
}

// SetReadPreference 在context和header中设置mongodb read preference，给dal使用
func SetReadPreference(ctx context.Context, header http.Header, mode common.ReadPreferenceMode) (context.Context,
	http.Header) {
	ctx = SetDBReadPreference(ctx, mode)
	header = SetHTTPReadPreference(header, mode)
	return ctx, header
}

// GetDBReadPreference get mongodb read preference from context
func GetDBReadPreference(ctx context.Context) common.ReadPreferenceMode {
	val := ctx.Value(common.ReadReferenceKey)
	if val != nil {
		mode, ok := val.(string)
		if ok {
			return common.ReadPreferenceMode(mode)
		}
	}
	return common.NilMode
}

// GetHTTPReadPreference get mongodb read preference from http header
func GetHTTPReadPreference(header http.Header) common.ReadPreferenceMode {
	mode := header.Get(common.ReadReferenceKey)
	if mode == "" {
		return common.NilMode
	}
	return common.ReadPreferenceMode(mode)
}
