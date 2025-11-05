/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package service

import (
	"net/http"

	"github.com/TencentBlueKing/bk-cmdb/pkg/constant"
	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// I18nMiddleware i18n middleware
func I18nMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := pickTag(r)
		if err := i18n.Validate(lang); err != nil {
			log.Error(r.Context(), "invalid language", "lang", lang, log.E(err))
			err := &cerr.RespError{
				Code:        cerr.InvalidArgument,
				DetailError: err,
			}
			_ = rest.APIError(kit.DefaultKit(), err).Render(w)
			return
		}

		r.Header.Set(constant.HTTPLanguageHeader, string(lang))
		next.ServeHTTP(w, r)
	})
}

func pickTag(r *http.Request) i18n.LanguageType {
	if c, err := r.Cookie(constant.HTTPLanguageCookie); err == nil && c.Value != "" {
		return i18n.LanguageType(c.Value)
	}

	if h := r.Header.Get(constant.HTTPLanguageHeader); h != "" {
		return i18n.LanguageType(h)
	}

	// if language is not set, use default language
	return i18n.DefaultLang()
}

// Authentication 统一鉴权中间件
func Authentication(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		// TODO: 这里需要根据实际的鉴权逻辑来实现
		r.Header.Set(constant.AppCodeHeader, "test")
		r.Header.Set(constant.UserHeader, "test")
		r.Header.Set(constant.TenantHeader, "default")

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
