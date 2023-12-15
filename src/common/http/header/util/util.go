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

// Package util defines http header related utility functions
package util

import (
	"context"
	"net/http"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/util"
)

// CCHeader generate cc header from http header
func CCHeader(header http.Header) http.Header {
	newHeader := make(http.Header)
	httpheader.SetRid(newHeader, httpheader.GetRid(header))
	httpheader.SetUser(newHeader, httpheader.GetUser(header))
	httpheader.SetUserToken(newHeader, httpheader.GetUserToken(header))
	httpheader.SetUserTicket(newHeader, httpheader.GetUserTicket(header))
	httpheader.SetLanguage(newHeader, httpheader.GetLanguage(header))
	httpheader.SetSupplierAccount(newHeader, httpheader.GetSupplierAccount(header))
	httpheader.SetAppCode(newHeader, httpheader.GetAppCode(header))
	httpheader.SetReqRealIP(newHeader, httpheader.GetReqRealIP(header))
	if httpheader.IsReqFromWeb(header) {
		httpheader.SetReqFromWeb(newHeader)
	}
	newHeader.Add(common.ReadReferenceKey, header.Get(common.ReadReferenceKey))

	return newHeader
}

// GenCommonHeader generate common cmdb http header, use default value if parameter is not set
func GenCommonHeader(user, supplierAccount, rid string) http.Header {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	if user == "" {
		user = common.CCSystemOperatorUserName
	}

	if supplierAccount == "" {
		supplierAccount = common.BKDefaultOwnerID
	}

	if rid == "" {
		rid = util.GenerateRID()
	}

	httpheader.SetUser(header, user)
	httpheader.SetSupplierAccount(header, supplierAccount)
	httpheader.SetRid(header, rid)
	return header
}

// GenDefaultHeader generate default cmdb http header
func GenDefaultHeader() http.Header {
	return GenCommonHeader("", "", "")
}

// ConvertLegacyHeader convert legacy header to new http header, compatible for esb request
func ConvertLegacyHeader(header http.Header) http.Header {
	if httpheader.GetUser(header) == "" {
		httpheader.SetUser(header, header.Get(httpheader.BKHTTPHeaderUser))
	}

	if httpheader.GetSupplierAccount(header) == "" {
		supplierAccount := header.Get(httpheader.BKHTTPOwner)
		if supplierAccount == "" {
			supplierAccount = header.Get(httpheader.BKHTTPOwnerID)
		}
		httpheader.SetSupplierAccount(header, supplierAccount)
	}

	if httpheader.GetRid(header) == "" {
		httpheader.SetRid(header, header.Get(httpheader.BKHTTPCCRequestID))
	}

	if httpheader.GetLanguage(header) == "" {
		httpheader.SetLanguage(header, header.Get(httpheader.BKHTTPLanguage))
	}

	if httpheader.GetAppCode(header) == "" {
		httpheader.SetAppCode(header, header.Get(httpheader.BKHTTPRequestAppCode))
	}
	return header
}

// NewHeaderFromContext new cmdb header by context
func NewHeaderFromContext(ctx context.Context) http.Header {
	rid := ctx.Value(common.ContextRequestIDField)
	ridValue, _ := rid.(string)

	user := ctx.Value(common.ContextRequestUserField)
	userValue, _ := user.(string)

	owner := ctx.Value(common.ContextRequestOwnerField)
	ownerValue, _ := owner.(string)

	return GenCommonHeader(userValue, ownerValue, ridValue)
}

// BuildHeader build cmdb header by user & supplier account
func BuildHeader(user string, supplierAccount string) http.Header {
	return GenCommonHeader(user, supplierAccount, "")
}
