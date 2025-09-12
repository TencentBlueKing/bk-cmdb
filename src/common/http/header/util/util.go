/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

	"configcenter/pkg/tenant/tools"
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
	httpheader.SetTenantID(newHeader, httpheader.GetTenantID(header))
	httpheader.SetAppCode(newHeader, httpheader.GetAppCode(header))
	httpheader.SetReqRealIP(newHeader, httpheader.GetReqRealIP(header))
	if httpheader.IsReqFromWeb(header) {
		httpheader.SetReqFromWeb(newHeader)
	}
	newHeader.Add(common.ReadReferenceKey, header.Get(common.ReadReferenceKey))

	return newHeader
}

// GenCommonHeader generate common cmdb http header, use default value if parameter is not set
func GenCommonHeader(user, tenantID, rid string) http.Header {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	if user == "" {
		user = common.CCSystemOperatorUserName
	}

	if tenantID == "" {
		tenantID = tools.GetDefaultTenant()
	}

	if rid == "" {
		rid = util.GenerateRID()
	}

	httpheader.SetUser(header, user)
	httpheader.SetTenantID(header, tenantID)
	httpheader.SetRid(header, rid)
	return header
}

// GenDefaultHeader generate default cmdb http header
func GenDefaultHeader() http.Header {
	return GenCommonHeader("", "", "")
}

// NewHeader take out the required header value and create a new header
func NewHeader(header http.Header) http.Header {
	newHeader := http.Header{}
	newHeader.Set("Content-Type", "application/json")

	httpheader.SetUser(newHeader, httpheader.GetUser(header))

	httpheader.SetTenantID(newHeader, httpheader.GetTenantID(header))

	httpheader.SetRid(newHeader, httpheader.GetRid(header))

	httpheader.SetLanguage(newHeader, httpheader.GetLanguage(header))

	httpheader.SetAppCode(newHeader, httpheader.GetAppCode(header))

	httpheader.SetTXId(newHeader, httpheader.GetTXId(header))
	httpheader.SetTXTimeout(newHeader, httpheader.GetTXTimeout(header))
	httpheader.SetTXTenant(newHeader, httpheader.GetTXTenant(header))

	if httpheader.IsReqFromWeb(header) {
		httpheader.SetReqFromWeb(newHeader)
	}

	return newHeader
}

// ConvertLegacyHeader convert legacy header to new http header, compatible for esb request
func ConvertLegacyHeader(header http.Header) http.Header {
	newHeader := NewHeader(header)

	if httpheader.GetUser(header) == "" {
		httpheader.SetUser(newHeader, header.Get(httpheader.BKHTTPHeaderUser))
	}

	// if multi tenant mode is not enabled and tenantID = "", set default tenant
	if httpheader.GetTenantID(header) == "" {
		tenantID := tools.GetDefaultTenant()
		if tenantID == common.BKSingleTenantID {
			httpheader.SetTenantID(newHeader, tenantID)
		}
	}

	if httpheader.GetRid(header) == "" {
		httpheader.SetRid(newHeader, header.Get(httpheader.BKHTTPCCRequestID))
	}

	if httpheader.GetLanguage(header) == "" {
		httpheader.SetLanguage(newHeader, header.Get(httpheader.BKHTTPLanguage))
	}

	if httpheader.GetAppCode(header) == "" {
		httpheader.SetAppCode(newHeader, header.Get(httpheader.BKHTTPRequestAppCode))
	}

	return newHeader
}

// NewHeaderFromContext new cmdb header by context
func NewHeaderFromContext(ctx context.Context) http.Header {
	rid := ctx.Value(common.ContextRequestIDField)
	ridValue, _ := rid.(string)

	user := ctx.Value(common.ContextRequestUserField)
	userValue, _ := user.(string)

	tenant := ctx.Value(common.ContextRequestTenantField)
	tenantValue, _ := tenant.(string)

	return GenCommonHeader(userValue, tenantValue, ridValue)
}

// BuildHeader build cmdb header by user & tenant
func BuildHeader(user string, tenantID string) http.Header {
	return GenCommonHeader(user, tenantID, "")
}
