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

package i18n

import (
	"context"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// I18NMiddleWare i18n middleware
func I18NMiddleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := pickTag(r, defaultI18NManager)
		p := message.NewPrinter(tag, message.Catalog(defaultI18NManager.Catalog()))
		printer := &TranslatePrinter{printer: p}
		ctx := contextWithTranslator(r.Context(), printer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithTranslator(ctx context.Context, trans *TranslatePrinter) context.Context {
	return context.WithValue(ctx, translatorKey, trans)
}

// pickTag pick language Tag from request
func pickTag(r *http.Request, m *I18NManager) language.Tag {
	if c, err := r.Cookie(HTTPCookieLanguage); err == nil && c.Value != "" {
		if t, e := language.Parse(c.Value); e == nil {
			return m.Match(t)
		}
	}

	if h := r.Header.Get(BKHTTPLanguage); h != "" {
		if t, e := language.Parse(h); e == nil {
			return m.Match(t)
		}
	}
	return m.Fallback()
}
