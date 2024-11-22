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

package header

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/require"
)

func TestGetActionLanguage(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/call", nil)

	language := GetLanguage(restful.NewRequest(req).Request.Header)

	SetLanguage(req.Header, "cn")
	language = GetLanguage(restful.NewRequest(req).Request.Header)
	require.Equal(t, "cn", language)

	SetLanguage(req.Header, "cnn")
	language = GetLanguage(restful.NewRequest(req).Request.Header)
	require.NotEqual(t, "cn", language)
}

func TestHeader(t *testing.T) {
	header := make(http.Header)
	SetLanguage(header, "zh")
	SetUser(header, "user")
	SetTenantID(header, "owner")
	SetRid(header, "rid")

	req := &http.Request{Header: header}
	r := restful.NewRequest(req)
	if GetLanguage(header) != "zh" {
		t.Fail()
	}
	if GetLanguage(r.Request.Header) != "zh" {
		t.Fail()
	}

	if GetUser(header) != "user" {
		t.Fail()
	}

	if GetTenantID(header) != "tenant_id" {
		t.Fail()
	}

	if GetRid(header) != "rid" {
		t.Fail()
	}
}
