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

// Package header defines http header related logics
package header

// deprecated old http header keys, compatible for esb
// TODO remove these headers after esb is not used
const (
	// BKHTTPHeaderUser current request http request header fields name for login user
	BKHTTPHeaderUser = "BK_User"
	// BKHTTPLanguage the language key word
	BKHTTPLanguage = "HTTP_BLUEKING_LANGUAGE"
	// BKHTTPRequestAppCode is the blueking app code
	BKHTTPRequestAppCode = "Bk-App-Code"
	// BKHTTPCCRequestID cc request id cc_request_id
	BKHTTPCCRequestID = "Cc_Request_Id"
)

// new http header keys in the standard of api gateway
const (
	// BkApigwRidHeader is the blueking api gateway request id http header key
	BkApigwRidHeader = "X-Bkapi-Request-Id"

	// BkRidHeader is the request id http header key
	BkRidHeader = "X-Request-ID"

	// BkAuthHeader is the blueking api gateway authorization http header key
	BkAuthHeader = "X-Bkapi-Authorization"

	// BkJWTHeader is the blueking api gateway jwt http header key
	BkJWTHeader = "X-Bkapi-JWT"

	// AppCodeHeader is the blueking app code http header key, its value is from jwt info
	AppCodeHeader = "X-Bkcmdb-App-Code"

	// UserHeader is the username http header key, its value is from jwt info
	UserHeader = "X-Bkcmdb-User"

	// UserTokenHeader is the blueking user token http header key, its value is from common.HTTPCookieBKToken cookie
	UserTokenHeader = "X-Bkcmdb-User-Token"

	// UserTicketHeader is the blueking user ticket http header key, its value is from common.HTTPCookieBKTicket cookie
	UserTicketHeader = "X-Bkcmdb-User-Ticket"

	// LanguageHeader is the language http header key, its value is from common.HTTPCookieLanguage cookie
	LanguageHeader = "X-Bkcmdb-Language"

	// TenantHeader is tenant http header key
	TenantHeader = "X-Bk-Tenant-Id"

	// ReqFromWebHeader is the http header key that represents if request is from web server
	ReqFromWebHeader = "X-Bkcmdb-Request-From-Web"

	// ReqRealIPHeader is the request real ip http header key
	ReqRealIPHeader = "X-Real-Ip"

	// IsInnerReqHeader is the http header key that represents if request is an inner request
	IsInnerReqHeader = "X-Bkcmdb-Is-Inner-Request"
)
