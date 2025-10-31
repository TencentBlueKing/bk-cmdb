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

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// I18nMiddleware i18n middleware
func I18nMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang, err := i18n.GetDefaultManager().Validate(pickTag(r))
		if err != nil {
			log.Error(r.Context(), "invalid language", "lang", lang, log.E(err))
			err := &cerr.RespError{
				Code:        cerr.INVALID_REQUEST,
				DetailError: err,
			}
			rest.ApiRespError(err, w, r, cerr.INVALID_REQUEST)
			return
		}

		ctx := i18n.ContextWithLang(r.Context(), lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func pickTag(r *http.Request) i18n.LanguageType {
	if c, err := r.Cookie(i18n.HTTPCookieLanguage); err == nil && c.Value != "" {
		return i18n.LanguageType(c.Value)
	}

	if h := r.Header.Get(i18n.BKHTTPLanguage); h != "" {
		return i18n.LanguageType(h)
	}

	return i18n.DefaultLanguage
}
