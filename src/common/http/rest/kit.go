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

package rest

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/sharding"
)

// Kit stores the metadata info of a request
type Kit struct {
	Rid             string
	Header          http.Header
	Ctx             context.Context
	CCError         errors.DefaultCCErrorIf
	User            string
	SupplierAccount string
}

// NewKit 产生一个新的kit， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewKit() *Kit {
	newHeader := headerutil.CCHeader(kit.Header)
	newKit := *kit
	newKit.Header = newHeader
	return &newKit
}

// NewHeader 产生一个新的header， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewHeader() http.Header {
	return headerutil.CCHeader(kit.Header)
}

// NewKitFromHeader generate a new kit from http header.
func NewKitFromHeader(header http.Header, errorIf errors.CCErrorIf) *Kit {
	return &Kit{
		Rid:             httpheader.GetRid(header),
		Header:          header,
		Ctx:             util.NewContextFromHTTPHeader(header),
		CCError:         errorIf.CreateDefaultCCErrorIf(httpheader.GetLanguage(header)),
		User:            httpheader.GetUser(header),
		SupplierAccount: httpheader.GetSupplierAccount(header),
	}
}

// NewKit generate a new kit
func NewKit() *Kit {
	return NewKitFromHeader(headerutil.GenDefaultHeader(), errors.GetGlobalCCError())
}

// ShardOpts returns sharding options
func (kit *Kit) ShardOpts() sharding.ShardOpts {
	return sharding.NewShardOpts().WithTenant(kit.SupplierAccount)
}

// SysShardOpts returns sharding options for system
func (kit *Kit) SysShardOpts() sharding.ShardOpts {
	return sharding.NewShardOpts().WithIgnoreTenant().WithTenant(kit.SupplierAccount)
}
