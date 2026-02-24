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

package header

import (
	"net/http"

	"configcenter/src/common"
)

// GetRid get request id from http header
func GetRid(header http.Header) string {
	if rid := header.Get(BkRidHeader); rid != "" {
		return rid
	}
	// get request id from apigw rid header if not set
	return header.Get(BkApigwRidHeader)
}

// GetBkJWT get blueking api gateway jwt info from http header
func GetBkJWT(header http.Header) string {
	return header.Get(BkJWTHeader)
}

// GetAppCode get blueking app code from http header
func GetAppCode(header http.Header) string {
	return header.Get(AppCodeHeader)
}

// GetUser get username from http header
func GetUser(header http.Header) string {
	if user := header.Get(UserHeader); user != "" {
		return user
	}
	if user := header.Get("BK_User"); user != "" {
		return user
	}
	return ""
}

// GetUserToken get blueking user token from http header
func GetUserToken(header http.Header) string {
	return header.Get(UserTokenHeader)
}

// GetUserTicket get blueking user ticket from http header
func GetUserTicket(header http.Header) string {
	return header.Get(UserTicketHeader)
}

// GetLanguage get language from http header
func GetLanguage(header http.Header) string {
	return header.Get(LanguageHeader)
}

// GetTenantID get supplier account from http header
func GetTenantID(header http.Header) string {
	return header.Get(TenantHeader)
}

// IsReqFromWeb check if request is from web server
func IsReqFromWeb(header http.Header) bool {
	return header.Get(ReqFromWebHeader) == "true"
}

// GetReqRealIP get request real ip from http header
func GetReqRealIP(header http.Header) string {
	return header.Get(ReqRealIPHeader)
}

// GetTXId get transaction id from http header
func GetTXId(header http.Header) string {
	return header.Get(common.TransactionIdHeader)
}

// GetTXTimeout get transaction timeout from http header
func GetTXTimeout(header http.Header) string {
	return header.Get(common.TransactionTimeoutHeader)
}

// GetTXTenant get transaction tenantID from http header
func GetTXTenant(header http.Header) string {
	return header.Get(common.TransactionTenantIDHeader)
}

// SetRid set request id to http header
func SetRid(header http.Header, value string) {
	header.Set(BkRidHeader, value)
}

// SetBkAuth set blueking api gateway authorization info to http header
func SetBkAuth(header http.Header, value string) http.Header {
	h := make(http.Header)
	for key := range header {
		h.Set(key, header.Get(key))
	}
	h.Set(BkAuthHeader, value)
	return h
}

// SetBkJWT set blueking api gateway jwt info to http header
func SetBkJWT(header http.Header, value string) {
	header.Set(BkJWTHeader, value)
}

// SetAppCode set blueking app code to http header
func SetAppCode(header http.Header, value string) {
	header.Set(AppCodeHeader, value)
}

// SetUser set username to http header
func SetUser(header http.Header, value string) {
	header.Set(UserHeader, value)
	header.Set("BK_User", value)
}

// SetUserToken set blueking user token to http header
func SetUserToken(header http.Header, value string) {
	header.Set(UserTokenHeader, value)
}

// SetUserTicket set blueking user ticket to http header
func SetUserTicket(header http.Header, value string) {
	header.Set(UserTicketHeader, value)
}

// SetLanguage set language to http header
func SetLanguage(header http.Header, value string) {
	header.Set(LanguageHeader, value)
}

// SetTenantID set tenant id to http header
func SetTenantID(header http.Header, value string) {
	header.Set(TenantHeader, value)
}

// SetReqFromWeb set the request from web server flag to http header
func SetReqFromWeb(header http.Header) {
	header.Set(ReqFromWebHeader, "true")
}

// SetReqRealIP set request real ip to http header
func SetReqRealIP(header http.Header, value string) {
	header.Set(ReqRealIPHeader, value)
}

// SetTXId set transaction id to http header
func SetTXId(header http.Header, value string) {
	header.Set(common.TransactionIdHeader, value)
}

// SetTXTimeout set transaction timeout to http header
func SetTXTimeout(header http.Header, value string) {
	header.Set(common.TransactionTimeoutHeader, value)
}

// SetTXTenant set transaction tenant id to http header
func SetTXTenant(header http.Header, value string) {
	header.Set(common.TransactionTenantIDHeader, value)
}

// AddRid add request id to http header
func AddRid(header http.Header, value string) {
	if GetRid(header) != value {
		header.Add(BkRidHeader, value)
	}
}

// AddUser add user to http header
func AddUser(header http.Header, value string) {
	if GetUser(header) != value {
		header.Add(UserHeader, value)
	}
}

// AddTenantID add supplier account to http header
func AddTenantID(header http.Header, value string) {
	if GetTenantID(header) != value {
		header.Add(TenantHeader, value)
	}
}

// AddLanguage add language to http header
func AddLanguage(header http.Header, value string) {
	if GetLanguage(header) != value {
		header.Add(LanguageHeader, value)
	}
}

// IsInnerReq check if request is inner request
func IsInnerReq(header http.Header) bool {
	return header.Get(IsInnerReqHeader) == "true"
}

// SetIsInnerReqHeader set the request is inner flag to http header
func SetIsInnerReqHeader(header http.Header) {
	header.Set(IsInnerReqHeader, "true")
}
